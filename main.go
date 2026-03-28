package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/untref-ayp2-guias/trabajo-pr-ctico-grupal-alpha.git/datos"
	"github.com/untref-ayp2-guias/trabajo-pr-ctico-grupal-alpha.git/estructuras"
	"github.com/untref-ayp2-guias/trabajo-pr-ctico-grupal-alpha.git/recursos"
	"github.com/untref-ayp2-guias/trabajo-pr-ctico-grupal-alpha.git/simulacion"
	"github.com/untref-ayp2-guias/trabajo-pr-ctico-grupal-alpha.git/sistema"
)

func seleccionarEscenario() string {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Seleccione un escenario:")
		fmt.Println("1. Tráfico aéreo bajo")
		fmt.Println("2. Tráfico aéreo medio")
		fmt.Println("3. Tráfico aéreo alto")
		fmt.Print("Opción: ")

		input, _ := reader.ReadString('\n')
		op := strings.TrimSpace(input)

		switch op {
		case "1":
			return "trafico_bajo"
		case "2":
			return "trafico_medio"
		case "3":
			return "trafico_alto"
		default:
			fmt.Println("Opción inválida. Intente nuevamente.")
		}
	}
}

func mostrarMenu(sim *simulacion.Simulador) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("--------------------------------")
		fmt.Printf("Tiempo actual: %s\n", sim.TiempoActual().Format("2006-01-02 15:04"))
		fmt.Println("--------------------------------")
		fmt.Println("Seleccione una opción:")
		fmt.Println("1. Avanzar al siguiente evento")
		fmt.Println("2. Adelantar 30 minutos")
		fmt.Println("3. Ejecutar hasta el fnal")
		fmt.Println("4. Buscar pasajero")
		fmt.Println("5. Habilitar mostrador")
		fmt.Println("6. Cerrar mostrador")
		fmt.Println("7. Consultar estado mostradores")
		fmt.Println("8. Descomprimir archivos de vuelo")
		fmt.Println("x. Salir")
		fmt.Print("Opción: ")
		input, _ := reader.ReadString('\n')
		op := strings.TrimSpace(input)

		switch op {
		case "1":
			sim.EjecutarSiguienteEvento()
		case "2":
			sim.EjecutarHasta(sim.TiempoActual().Add(30 * time.Minute))
		case "3":
			if sim.Eventos() == 0 {
				fmt.Print("No hay más eventos.\n")
			} else {
				sim.Ejecutar()
			}
		case "4":
			fmt.Print("Ingrese el documento del pasajero: ")
			input, _ := reader.ReadString('\n')
			doc := strings.TrimSpace(input)
			pasajeroInfo := estructuras.BuscarPasajeroPorDocumento(datos.Clientes, doc)
			fmt.Println(pasajeroInfo)
		case "5":
			fmt.Print("Número de mostrador a habilitar: ")
			numStr, _ := reader.ReadString('\n')
			numStr = strings.TrimSpace(numStr)
			num, err := strconv.Atoi(numStr)
			if err == nil {
				sistema.ModuloGlobalCheckIn.HabilitarMostrador(num)
			} else {
				fmt.Println("Número inválido.")
			}
		case "6":
			fmt.Print("Número de mostrador a cerrar: ")
			numStr, _ := reader.ReadString('\n')
			numStr = strings.TrimSpace(numStr)
			num, err := strconv.Atoi(numStr)
			if err == nil {
				sistema.ModuloGlobalCheckIn.CerrarMostrador(num)
			} else {
				fmt.Println("Número inválido.")
			}
		case "7":
			fmt.Printf("Cantidad de mostradores en el aeropuerto: %d\n", len(sistema.ModuloGlobalCheckIn.Mostradores))
			fmt.Printf("Cantidad de mostradores habilitados: %d\n", sistema.ModuloGlobalCheckIn.CantidadMostradoresActivos())
		case "8":
			carpeta := "archivos_vuelos"
			err := os.MkdirAll(carpeta, 0755)
			if err != nil {
				fmt.Printf("Error al crear la carpeta '%s': %v\n", carpeta, err)
				break
			}

			for _, vuelo := range datos.Vuelos {
				nombreArchivo := filepath.Join(carpeta, fmt.Sprintf("vuelo_%s.txt", vuelo.NumVuelo))
				archivoHuff := fmt.Sprintf("output/vuelo_%s.huff", vuelo.NumVuelo)

				if _, err := os.Stat(archivoHuff); os.IsNotExist(err) {
					fmt.Printf("El archivo %s no existe. Se omite.\n", archivoHuff)
					continue
				}
				err := recursos.DescomprimirArchivosHuffman(archivoHuff, nombreArchivo)
				if err != nil {
					fmt.Printf("Error al descomprimir %s: %v\n", archivoHuff, err)
				} else {
					fmt.Printf("Archivo %s generado correctamente.\n", nombreArchivo)
				}
			}
		case "x":
			fmt.Println("Saliendo de la simulación.")
			return
		default:
			fmt.Println("Opción inválida.")
		}
	}
}
func main() {
	escenario := seleccionarEscenario()

	err := datos.InicializarDatos(escenario)
	if err != nil {
		log.Fatalf("Error al inicializar datos: %v", err)
	}

	simulador := simulacion.NuevoSimulador()
	sistema.ModuloGlobalCheckIn = sistema.NewModuloCheckIn(10) //cantidad de mostradores que hay en el aeropuerto

	// genera el evento de inicio de checkin de todos los vuelos
	simulacion.GenerarEventosIniciales(simulador, datos.Vuelos)

	// genera las llegadas de los pasajeros al checkin
	for _, vuelo := range datos.Vuelos {
		sistema.PlanificarLlegadasPasajeros(vuelo, simulador.TiempoActual(), func(tiempo time.Time, pasajero *estructuras.Pasajero) {
			simulador.AgregarEvento(tiempo, simulacion.EventoLlegadaPasajero, pasajero)
		})
	}

	mostrarMenu(simulador)
}
