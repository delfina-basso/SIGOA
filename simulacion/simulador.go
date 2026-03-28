package simulacion

import (
	"container/heap"
	"fmt"
	"math/rand"
	"time"

	"github.com/untref-ayp2-guias/trabajo-pr-ctico-grupal-alpha.git/datos"
	"github.com/untref-ayp2-guias/trabajo-pr-ctico-grupal-alpha.git/estructuras"
	"github.com/untref-ayp2-guias/trabajo-pr-ctico-grupal-alpha.git/recursos"
	"github.com/untref-ayp2-guias/trabajo-pr-ctico-grupal-alpha.git/sistema"
)

type Simulador struct {
	tiempoActual time.Time
	eventos      ColaEventos
}

func NuevoSimulador() *Simulador {
	eventQueue := make(ColaEventos, 0)
	heap.Init(&eventQueue)

	//acomodar
	local, _ := time.LoadLocation("America/Argentina/Buenos_Aires")
	inicio := time.Date(2023, 10, 1, 5, 30, 0, 0, local)
	return &Simulador{
		tiempoActual: inicio,
		eventos:      eventQueue,
	}
}

func (s *Simulador) Eventos() int {
	return len(s.eventos)
}

func (s *Simulador) TiempoActual() time.Time {
	return s.tiempoActual
}

func (s *Simulador) AgregarEvento(tiempo time.Time, tipo TipoEvento, data any) {
	heap.Push(&s.eventos, &Evento{
		Tiempo: tiempo,
		Tipo:   tipo,
		Data:   data,
	})
}

func (s *Simulador) Ejecutar() {
	for s.eventos.Len() > 0 {
		evento := heap.Pop(&s.eventos).(*Evento)
		s.tiempoActual = evento.Tiempo
		s.ejecutarEvento(evento)
	}
}

func (s *Simulador) EjecutarHasta(hasta time.Time) {
	for s.eventos.Len() > 0 && !s.eventos[0].Tiempo.After(hasta) {
		evento := heap.Pop(&s.eventos).(*Evento)
		s.tiempoActual = evento.Tiempo
		s.ejecutarEvento(evento)
	}
	s.tiempoActual = hasta
	fmt.Printf("\n Tiempo de simulación adelantado hasta: %s\n\n", s.tiempoActual.Format("2006-01-02 15:04"))

}

func (s *Simulador) ejecutarEvento(e *Evento) {
	s.tiempoActual = e.Tiempo
	color := colorEvento(e.Tipo)
	reset := "\033[0m"

	switch e.Tipo {
	case EventoInicioCheckIn:
		vuelo := e.Data.(*estructuras.Vuelo)
		fmt.Printf("[%s] %sVuelo %s: Inicio del Check-In.%s\n",
			e.Tiempo.Format("2006-01-02 15:04"), color, vuelo.NumVuelo, reset)
		vuelo.Estado = "En Check-In"
		s.AgregarEvento(vuelo.FechaHoraProgramada.Add(-45*time.Minute), EventoCierreCheckIn, vuelo)

	case EventoLlegadaPasajero:
		pasajero := e.Data.(*estructuras.Pasajero)
		vuelo := datos.Vuelos[pasajero.Vuelo]
		cierreCheckIn := vuelo.FechaHoraProgramada.Add(-46 * time.Minute)

		if pasajero.HoraLlegada.After(cierreCheckIn) {
			pasajero.PerdioVuelo = true
			fmt.Printf("[%s] %sPasajero %s %s llegó después del cierre del check-in y perdió el vuelo %s.%s\n",
				e.Tiempo.Format("2006-01-02 15:04"), "\033[31m", pasajero.Nombre, pasajero.Apellido, pasajero.Vuelo, reset)
			return
		}

		fmt.Printf("[%s] Pasajero %s %s llegó al aeropuerto para el vuelo %s y se dirige a la cola de los mostradores.\n",
			e.Tiempo.Format("2006-01-02 15:04"),
			pasajero.Nombre, pasajero.Apellido, pasajero.Vuelo)

		sistema.ModuloGlobalCheckIn.LlegadaPasajero(pasajero, s.TiempoActual(), func(tiempo time.Time, mostrador *estructuras.MostradorCheckIn) {
			s.AgregarEvento(tiempo, EventoFinCheckIn, mostrador)
		})

	case EventoFinCheckIn:
		sistema.ModuloGlobalCheckIn.ProcesarFinCheckInPasajero(e.Data.(*estructuras.MostradorCheckIn), s.TiempoActual(), func(tiempo time.Time, mostrador *estructuras.MostradorCheckIn) {
			s.AgregarEvento(tiempo, EventoFinCheckIn, mostrador) //no es recursivo, es para que se procese el siguiente pasajero
		}) //no se llama a cierre checkin porque ese evento se agrega en generar eventos iniciales

	case EventoCierreCheckIn:
		vuelo := e.Data.(*estructuras.Vuelo)
		sistema.ModuloGlobalCheckIn.ProcesarCierreCheckIn(vuelo, s.TiempoActual())
		fmt.Printf("[%s] %sVuelo %s: Cierre del Check-In.%s\n",
			e.Tiempo.Format("2006-01-02 15:04"), color, vuelo.NumVuelo, reset)

		// fmt.Printf("%sLISTA DE PASAJEROS CON CHECK-IN CONFIRMADO:%s\n", color, reset)
		// for _, p := range vuelo.Pasajeros {
		// 	if p.HizoCheckIn && p.AsientoAsignado != "" {
		// 		fmt.Printf(" - %s %s. (%s)\n",
		// 			p.Nombre, p.Apellido, p.Documento)
		// 	}
		// }
		vuelo.Estado = "Check-In Cerrado"
		s.AgregarEvento(s.TiempoActual(), EventoInicioCarga, vuelo)

	case EventoInicioCarga:
		vuelo := e.Data.(*estructuras.Vuelo)
		vuelo.Estado = "En proceso de carga"
		sistema.AsignarCarga(vuelo, datos.Cargas)

		fmt.Printf("[%s] %sVuelo %s: Inicio de carga. %s\n",
			e.Tiempo.Format("2006-01-02 15:04"), color, vuelo.NumVuelo, reset)
		for _, c := range vuelo.Carga {
			fmt.Printf("Carga asignada: peso %dkg, volumen %.2fm³.\n", int(c.Peso), c.Volumen)
		}
		s.AgregarEvento(vuelo.FechaHoraProgramada.Add(-35*time.Minute), EventoFinCarga, vuelo)

	case EventoFinCarga:
		vuelo := e.Data.(*estructuras.Vuelo)
		vuelo.Estado = "Fin del proceso de carga"
		fmt.Printf("[%s] %sVuelo %s: Finalizó el proceso de carga. %d cargas asignadas. %s\n",
			e.Tiempo.Format("2006-01-02 15:04"), color, vuelo.NumVuelo, len(vuelo.Carga), reset)
		s.AgregarEvento(s.TiempoActual(), EventoInicioEmbarque, vuelo)

	case EventoInicioEmbarque:
		vuelo := e.Data.(*estructuras.Vuelo)
		vuelo.Estado = "En proceso de embarque"
		vuelo.ZonaHabilitada = 1
		fmt.Printf("[%s] %sVuelo %s: Inicio del proceso de embarque. Habilitación de zona 1.%s\n",
			e.Tiempo.Format("2006-01-02 15:04"), color, vuelo.NumVuelo, reset)
		inicioEmbarque := vuelo.FechaHoraProgramada.Add(-35 * time.Minute)

		sistema.GenerarLlegadasEmbarque(vuelo, inicioEmbarque, func(tiempo time.Time, pasajero *estructuras.Pasajero) {
			s.AgregarEvento(tiempo, EventoLlegadaEmbarque, pasajero)
		})

		s.AgregarEvento((inicioEmbarque.Add(10 * time.Minute)), EventoFinZonaEmbarque, vuelo)
		s.AgregarEvento((inicioEmbarque.Add(20 * time.Minute)), EventoFinZonaEmbarque, vuelo)
		s.AgregarEvento(vuelo.FechaHoraProgramada.Add(-5*time.Minute), EventoFinEmbarque, vuelo)

	case EventoLlegadaEmbarque:
		pasajero := e.Data.(*estructuras.Pasajero)
		vuelo := datos.Vuelos[pasajero.Vuelo]

		inicioEmbarque := vuelo.FechaHoraProgramada.Add(-35 * time.Minute)

		cierresZonas := map[string]time.Time{
			"1": inicioEmbarque.Add(10 * time.Minute),
			"2": inicioEmbarque.Add(20 * time.Minute),
			"3": inicioEmbarque.Add(30 * time.Minute),
		}

		llegada := e.Tiempo
		llegoTarde := false

		if pasajero.Categoria != "" {
			if llegada.After(cierresZonas["3"]) {
				llegoTarde = true
			}
		} else {
			cierreZona := cierresZonas[pasajero.ZonaEmbarque]
			if llegada.After(cierreZona) {
				llegoTarde = true
			}
		}

		if llegoTarde {
			pasajero.PerdioVuelo = true
			fmt.Printf("[%s] %sPasajero %s %s llegó después de finalizado el embarque y perdió el vuelo %s.%s\n",
				llegada.Format("2006-01-02 15:04"), "\033[31m", pasajero.Nombre, pasajero.Apellido, vuelo.NumVuelo, reset)
			return
		}

		fmt.Printf("[%s] %sPasajero %s %s embarcó el vuelo %s.%s\n",
			llegada.Format("2006-01-02 15:04"), "\033[32m", pasajero.Nombre, pasajero.Apellido, vuelo.NumVuelo, "\033[0m")

		pasajero.Embarcado = true

	case EventoFinZonaEmbarque:
		vuelo := e.Data.(*estructuras.Vuelo)

		zona := vuelo.ZonaHabilitada
		fmt.Printf("[%s] %sVuelo %s: Finalizó el embarque de la zona %v, habilitación de la zona %v.%s\n",
			e.Tiempo.Format("2006-01-02 15:04"), color, vuelo.NumVuelo, zona, (zona + 1), reset)
		vuelo.ZonaHabilitada++

	case EventoFinEmbarque:
		vuelo := e.Data.(*estructuras.Vuelo)

		color := colorEvento(EventoFinEmbarque)
		reset := "\033[0m"

		fmt.Printf("[%s] %sVuelo %s: Finalizó el embarque. Aeronave %s en preparación para el despegue.%s\n",
			e.Tiempo.Format("2006-01-02 15:04"), color, vuelo.NumVuelo, vuelo.AeronaveAsignada, reset)

		listaEmbarcados := make([]*estructuras.Pasajero, 0)
		// fmt.Printf("%sPasajeros embarcados:%s\n", color, reset)
		for _, p := range vuelo.Pasajeros {
			if p.Embarcado {
				// fmt.Printf(" - %s %s (%s)\n", p.Nombre, p.Apellido, p.Documento)
				listaEmbarcados = append(listaEmbarcados, p)
			}
		}
		listaEquipaje := make([]*estructuras.Equipaje, 0)
		// fmt.Printf("%sEquipaje de los pasajeros embarcados:%s\n", color, reset)
		for _, p := range vuelo.Pasajeros {
			if p.Embarcado {
				if equipaje, ok := datos.Equipajes[p.Documento]; ok {
					// fmt.Printf(" - %s %s. (%s) Bultos: %d. Peso: %.2f kg\n",
					// p.Nombre, p.Apellido, equipaje.DocumentoPasajero, equipaje.Bultos, equipaje.PesoTotal)
					listaEquipaje = append(listaEquipaje, equipaje)
				}
			}
		}
		listaNoPresentados := make([]*estructuras.Pasajero, 0)
		// fmt.Printf("%sPasajeros que no se presentaron:%s\n", color, reset)
		for _, p := range vuelo.Pasajeros {
			if p.PerdioVuelo || !p.HizoCheckIn || !p.Embarcado {
				// fmt.Printf(" - %s %s (%s)\n", p.Nombre, p.Apellido, p.Documento)
				listaNoPresentados = append(listaNoPresentados, p)
			}
		}
		// fmt.Printf("%sPasajeros restantes en lista de espera:%s\n", color, reset)
		// if len(vuelo.ListaEspera) == 0 {
		// 	fmt.Printf("Ninguno.\n")
		// } else {
		// 	for _, p := range vuelo.ListaEspera {
		// 		fmt.Printf(" - %s %s (%s)\n", p.Nombre, p.Apellido, p.Documento)
		// 	}
		// }

		sistema.GenerarArchivoVuelo(vuelo, e.Tiempo, listaEmbarcados, listaEquipaje, listaNoPresentados, vuelo.ListaEspera, vuelo.Carga)
		recursos.ComprimirConHuffman(vuelo)

		vuelo.Estado = "Listo para despegar"

		sistema.ColaDespacho.Push(vuelo)
		fmt.Printf("[%s] %sVuelo %s: Listo para despegar, esperando instrucciones de control.%s\n",
			s.TiempoActual().Format("2006-01-02 15:04"), "\033[35m", vuelo.NumVuelo, "\033[0m")

		// Programar despacho entre hora programada y 3 minutos después
		delay := time.Duration(rand.Intn(4)) * time.Minute
		tiempoDespacho := vuelo.FechaHoraProgramada.Add(delay)

		s.AgregarEvento(tiempoDespacho, EventoDespachoVuelo, vuelo)

	case EventoDespachoVuelo:
		vuelo := e.Data.(*estructuras.Vuelo)

		fmt.Printf("[%s] %sControl: Calculando línea del horizonte para el vuelo %s.%s\n", s.TiempoActual().Format("2006-01-02 15:04"),
			color, vuelo.NumVuelo, reset)
		err := sistema.GenerarArchivoLineaHorizonte(sistema.CalcularLineaHorizonte(datos.Edificios))
		if err != nil {
			fmt.Println("Error al generar línea del horizonte:", err)
		}

		fmt.Printf("[%s] %sControl: Vuelo %s autorizado a despegar.%s\n", s.TiempoActual().Format("2006-01-02 15:04"),
			color, vuelo.NumVuelo, reset)

		sistema.ColaDespacho.Pop()

		aeropuerto := datos.Aeropuertos[vuelo.Destino]
		fmt.Printf("[%s] %sVuelo %s: Despachado. Destino: %s, %s. Hora programada: %s. Hora de salida real: %s.%s\n",
			s.TiempoActual().Format("2006-01-02 15:04"), color, vuelo.NumVuelo, vuelo.Destino, aeropuerto.Provincia,
			vuelo.FechaHoraProgramada.Format("15:04"), s.TiempoActual().Format("15:04"), reset)

	default:
		fmt.Printf("[%s] %sEvento: %s -> %v%s\n",
			e.Tiempo.Format("2006-01-02 15:04"), color, e.Tipo, e.Data,
			reset)
	}
}

func (s *Simulador) EjecutarSiguienteEvento() {
	if s.eventos.Len() == 0 {
		fmt.Println("No hay más eventos programados.")
		return
	}
	evento := heap.Pop(&s.eventos).(*Evento)
	s.tiempoActual = evento.Tiempo
	s.ejecutarEvento(evento)
}

//esta en desuso, quizas la necesite luego
//pone bonito los datos impresos por pantalla
// func formatearDatoEvento(data any) string {
// 	switch v := data.(type) {
// 	case *estructuras.Vuelo:
// 		return fmt.Sprintf("Vuelo %s (%s -> %s)", v.NumVuelo, v.AeronaveAsignada, v.Destino)
// 	case *estructuras.Pasajero:
// 		return fmt.Sprintf("%s %s (DNI: %s)", v.Nombre, v.Apellido, v.Documento)
// 	default:
// 		return fmt.Sprintf("%v", data)
// 	}
// }

func colorEvento(tipo TipoEvento) string {
	switch tipo {
	case "LlegadaPasajero":
		return "\033[34m" // Azul
	case "InicioCheckIn":
		return "\033[34m" // Azul
	case "CierreCheckIn":
		return "\033[34m" // Azul
	case "InicioCarga":
		return "\033[33m" // Amarillo
	case "FinCarga":
		return "\033[33m" // Amarillo
	case "InicioEmbarque":
		return "\033[36m" // Celeste
	case "FinZonaEmbarque":
		return "\033[36m" // Celeste
	case "FinEmbarque":
		return "\033[36m" // Celeste
	case "DespachoVuelo":
		return "\033[35m" // Violeta
	default:
		return "\033[0m"

		//"\033[31m" // Rojo
	}
}
