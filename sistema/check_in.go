package sistema

import (
	"container/heap"
	"fmt"
	"math/rand"
	"time"

	"github.com/untref-ayp2-guias/trabajo-pr-ctico-grupal-alpha.git/datos"
	"github.com/untref-ayp2-guias/trabajo-pr-ctico-grupal-alpha.git/estructuras"
)

var ModuloGlobalCheckIn *ModuloCheckIn

// ModuloCheckIn controla los mostradores y la cola
type ModuloCheckIn struct {
	Cola               []*estructuras.Pasajero
	Mostradores        []*estructuras.MostradorCheckIn
	MostradoresActivos int
}

// NewModuloCheckIn crea un módulo con N mostradores, inicia con 2 mostradores abiertos
func NewModuloCheckIn(n int) *ModuloCheckIn {
	mostradores := make([]*estructuras.MostradorCheckIn, n)
	activos := 0
	for i := 0; i < n; i++ {
		if activos < 2 {
			mostradores[i] = &estructuras.MostradorCheckIn{ID: i + 1, Disponible: true, Activo: true}
			activos++
		} else {
			mostradores[i] = &estructuras.MostradorCheckIn{ID: i + 1, Disponible: true, Activo: false}
		}

	}
	return &ModuloCheckIn{
		Cola:               make([]*estructuras.Pasajero, 0),
		Mostradores:        mostradores,
		MostradoresActivos: 1,
	}
}

// PlanificarLlegadasPasajeros agenda eventos de llegada al simulador
func PlanificarLlegadasPasajeros(vuelo *estructuras.Vuelo, actual time.Time, onSuccess func(time.Time, *estructuras.Pasajero)) {
	salida := vuelo.FechaHoraProgramada

	inicioCheckIn := salida.Add(-2 * time.Hour)
	finMaximoLlegada := salida.Add(-35 * time.Minute) // 10 minutos para simular pasajeros que llegan tarde y pierden el vuelo como yo lol

	minutosVentana := int(finMaximoLlegada.Sub(inicioCheckIn).Minutes())

	for _, pasajero := range vuelo.Pasajeros {
		if pasajero.YaLlego {
			continue
		}
		offsetMin := rand.Intn(minutosVentana)
		llegada := inicioCheckIn.Add(time.Duration(offsetMin) * time.Minute)
		pasajero.HoraLlegada = llegada
		pasajero.YaLlego = true

		//para agregar evento llegada pasajero
		onSuccess(llegada, pasajero)
	}

	for _, pasajero := range vuelo.ListaEspera {
		if pasajero.YaLlego {
			continue
		}
		offsetMin := rand.Intn(minutosVentana)
		llegada := inicioCheckIn.Add(time.Duration(offsetMin) * time.Minute)
		pasajero.HoraLlegada = llegada
		pasajero.YaLlego = true
		onSuccess(llegada, pasajero)
	}
}

// LlegadaPasajero encola al pasajero y trata de asignarlo
func (mc *ModuloCheckIn) LlegadaPasajero(p *estructuras.Pasajero, actual time.Time, onSuccess func(time.Time, *estructuras.MostradorCheckIn)) {
	vuelo, ok := datos.Vuelos[p.Vuelo]
	if ok && vuelo.CheckInCerrado {
		return
	}

	mc.Cola = append(mc.Cola, p)
	mc.IntentarAtender(actual, onSuccess)
}

// IntentarAtender asigna pasajeros a mostradores libres y activos
func (mc *ModuloCheckIn) IntentarAtender(actual time.Time, onSuccess func(tiempo time.Time, mostrador *estructuras.MostradorCheckIn)) {
	nuevaCola := []*estructuras.Pasajero{}

	for _, pasajero := range mc.Cola {
		vuelo := datos.Vuelos[pasajero.Vuelo]
		if vuelo.CheckInCerrado {
			continue
		}

		// Intentar asignarlo a un mostrador disponible
		asignado := false
		for _, m := range mc.Mostradores {
			if m.Disponible && m.Activo {
				duracion := m.Atender(pasajero, actual)
				fin := actual.Add(duracion)

				//para agregar evento de fin checkin pasajerio
				onSuccess(fin, m)

				asignado = true
				break
			}
		}

		if !asignado {
			// Si no se pudo atender todavía, se mantiene en la cola
			nuevaCola = append(nuevaCola, pasajero)
		}
	}

	mc.Cola = nuevaCola
}

// ProcesarFinCheckIn finaliza atención y busca siguiente
func (mc *ModuloCheckIn) ProcesarFinCheckInPasajero(m *estructuras.MostradorCheckIn, actual time.Time, onSuccess func(time.Time, *estructuras.MostradorCheckIn)) {
	p := m.PasajeroActual
	if p != nil {
		vuelo, ok := datos.Vuelos[p.Vuelo]
		if ok && p.EstadoReserva == "Confirmada" {
			asignarAsientoYZona(p, vuelo)

			//si tiene equipaje printea el mensaje todo junto
			equipaje, ok := datos.Equipajes[p.Documento]
			if ok && equipaje.Bultos > 0 {
				vuelo.Equipaje[0] += equipaje.PesoTotal
				vuelo.Equipaje[1] += equipaje.VolumenTotal
				equipaje.NumeroVuelo = vuelo.NumVuelo
				fmt.Printf("[%s] Mostrador %d: %sTarjeta de embarque y ticket de equipaje de %s %s (%s) emitidos -> Vuelo: %s. Asiento: %s. Zona: %s. Bultos: %d.Peso Total: %.2f kg.%s\n",
					actual.Format("2006-01-02 15:04"), m.ID, "\033[32m", p.Nombre, p.Apellido, p.Documento, p.Vuelo, p.AsientoAsignado, p.ZonaEmbarque, equipaje.Bultos, equipaje.PesoTotal, "\033[0m")

			} else { //no tiene equipaje
				fmt.Printf("[%s] Mostrador %d: %sTarjeta de embarque de %s %s emitida -> Vuelo: %s. Asiento: %s. Zona: %s%s\n",
					actual.Format("2006-01-02 15:04"), m.ID, "\033[32m", p.Nombre, p.Apellido, p.Vuelo, p.AsientoAsignado, p.ZonaEmbarque, "\033[0m")
			}
		}

		m.FinalizarAtencionMostrador(actual)
		mc.IntentarAtender(actual, func(tiempo time.Time, mostrador *estructuras.MostradorCheckIn) {
			onSuccess(tiempo, mostrador)
		})
	}
}

func asignarAsientoYZona(p *estructuras.Pasajero, vuelo *estructuras.Vuelo) bool {
	asientosOcupados := make(map[string]bool)
	for _, pasajero := range vuelo.Pasajeros {
		asientosOcupados[pasajero.AsientoAsignado] = true
	}

	for _, config := range vuelo.ConfiguracionAsientos {
		for i := config.AsientoInicial; i <= config.AsientoFinal; i++ {
			asientoStr := fmt.Sprintf("%d", i)
			if !asientosOcupados[asientoStr] {
				p.AsientoAsignado = asientoStr
				p.ZonaEmbarque = fmt.Sprintf("%d", config.Zona)
				return true
			}
		}
	}

	fmt.Printf("Pasajero %s %s (%s) no pudo ser asignado a un asiento porque el vuelo %s está completo.\n",
		p.Nombre, p.Apellido, p.Documento, p.Vuelo)
	p.PerdioVuelo = true

	return false
}

// HabilitarMostrador activa un mostrador específico
func (mc *ModuloCheckIn) HabilitarMostrador(id int) {
	if id > 0 && id <= len(mc.Mostradores) {
		if !mc.Mostradores[id-1].Activo {
			mc.Mostradores[id-1].Activo = true
			mc.MostradoresActivos++
			fmt.Printf("Mostrador %d habilitado.\n", id)
		} else {
			fmt.Printf("Mostrador %d ya se encuentra habilitado.\n", id)
		}
	} else {
		fmt.Printf("No existe el mostrador %d. El aeropuerto dispone de %d mostradores.\n", id, len(mc.Mostradores))
	}
}

// CerrarMostrador desactiva un mostrador (lo marca como inactivo)
func (mc *ModuloCheckIn) CerrarMostrador(id int) {
	if id > 0 && id <= len(mc.Mostradores) {
		if mc.Mostradores[id-1].Activo {
			if mc.CantidadMostradoresActivos() == 1 {
				fmt.Printf("Debe haber por lo menos 1 mostrador activo\n")
				return
			}
			mc.Mostradores[id-1].Activo = false
			mc.MostradoresActivos--
			fmt.Printf("Mostrador %d cerrado.\n", id)
		} else {
			fmt.Printf("Mostrador %d ya se encuentra cerrado.\n", id)
		}
	} else {
		fmt.Printf("No existe el mostrador %d.\n", id)
	}
}

// devolver cantidad de mostradores activos
func (mc *ModuloCheckIn) CantidadMostradoresActivos() int {
	return mc.MostradoresActivos
}

func (mc *ModuloCheckIn) ProcesarCierreCheckIn(vuelo *estructuras.Vuelo, actual time.Time) {
	vuelo.CheckInCerrado = true
	fmt.Printf("[%s] %sCerrando check-in del vuelo %s. Los mostradores no atenderán pasajeros del vuelo que estén en la cola.%s\n", actual.Format("2006-01-02 15:04"), "\033[34m", vuelo.NumVuelo, "\033[0m")
	//eliminar de la cola de espera a los que no llegaron a un mostrador antes de que cierre el checkin
	nuevaCola := []*estructuras.Pasajero{}
	for _, p := range mc.Cola {
		if p.Vuelo == vuelo.NumVuelo {
			p.PerdioVuelo = true
			fmt.Printf("[%s] %s%s %s perdió el vuelo %s por no llegar a tiempo a los mostradores.%s\n", actual.Format("2006-01-02 15:04"), "\033[31m", p.Nombre, p.Apellido, p.Vuelo, "\033[0m")
		} else {
			nuevaCola = append(nuevaCola, p)
		}
	}
	mc.Cola = nuevaCola

	mc.AsignarAsientosDesdeListaEspera(vuelo, datos.ConfiguracionAsientos, actual)
}

// AsignarAsientosDesdeListaEspera asigna asientos por prioridad al cerrar check-in
func (mc *ModuloCheckIn) AsignarAsientosDesdeListaEspera(
	vuelo *estructuras.Vuelo,
	configs map[string][]*estructuras.ConfiguracionAsientos, actual time.Time) {

	fmt.Printf("[%s] %sVuelo %s: Asignando asientos a pasajeros en lista de espera:%s\n",
		actual.Format("2006-01-02 15:04"), "\033[34m", vuelo.NumVuelo, "\033[0m")

	if len(vuelo.ListaEspera) == 0 {
		fmt.Printf("[%s] %sNo hay pasajeros en lista de espera.%s\n",
			actual.Format("2006-01-02 15:04"), "\033[31m", "\033[0m")
		return
	}

	configsVuelo := configs[vuelo.AeronaveAsignada]

	//Crear mapa de asientos ocupados
	asientosOcupados := make(map[string]bool)
	for _, p := range vuelo.Pasajeros {
		if p.AsientoAsignado != "" {
			asientosOcupados[p.AsientoAsignado] = true
		}
	}

	//Crear cola de prioridad
	colap := &estructuras.ColaPrioridad{}
	heap.Init(colap)
	for _, pasajero := range vuelo.ListaEspera {
		heap.Push(colap, &estructuras.PasajeroPrioritario{Pasajero: pasajero})

	}

	//Asignar asientos
	for _, config := range configsVuelo {
		for asiento := config.AsientoInicial; asiento <= config.AsientoFinal && colap.Len() > 0; asiento++ {
			numAsiento := fmt.Sprintf("%d", asiento)
			if !asientosOcupados[numAsiento] {
				item := heap.Pop(colap).(*estructuras.PasajeroPrioritario)
				pasajero := item.Pasajero
				pasajero.AsientoAsignado = numAsiento
				pasajero.ZonaEmbarque = fmt.Sprintf("%d", config.Zona)
				pasajero.HizoCheckIn = true
				vuelo.Pasajeros = append(vuelo.Pasajeros, pasajero)
				asientosOcupados[numAsiento] = true

				//misma logica que procesarfincheckinpasjero
				equipaje, ok := datos.Equipajes[pasajero.Documento]
				if ok && equipaje.Bultos > 0 {
					vuelo.Equipaje[0] += equipaje.PesoTotal
					vuelo.Equipaje[1] += equipaje.VolumenTotal
					equipaje.NumeroVuelo = vuelo.NumVuelo
					fmt.Printf("[%s] %sTarjeta de embarque y ticket de equipaje de %s %s (%s) emitidos -> Vuelo: %s. Asiento: %s. Zona: %s. Bultos: %d.Peso Total: %.2f kg.%s\n",
						actual.Format("2006-01-02 15:04"), "\033[32m", pasajero.Nombre, pasajero.Apellido, pasajero.Documento, pasajero.Vuelo, pasajero.AsientoAsignado, pasajero.ZonaEmbarque, equipaje.Bultos, equipaje.PesoTotal, "\033[0m")
				} else {
					fmt.Printf("[%s] %sTarjeta de embarque de %s %s emitida -> Vuelo: %s. Asiento: %s. Zona: %s%s\n",
						actual.Format("2006-01-02 15:04"),
						"\033[32m", pasajero.Nombre, pasajero.Apellido, pasajero.Vuelo, pasajero.AsientoAsignado, pasajero.ZonaEmbarque, "\033[0m")
				}

			}
		}

		//Dejar en lista de espera los que no consiguieron asiento
		restantes := []*estructuras.Pasajero{}
		for colap.Len() > 0 {
			item := heap.Pop(colap).(*estructuras.PasajeroPrioritario)
			restantes = append(restantes, item.Pasajero)
		}
		vuelo.ListaEspera = restantes
	}
}
