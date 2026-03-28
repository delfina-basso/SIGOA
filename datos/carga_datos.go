package datos

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/untref-ayp2-guias/trabajo-pr-ctico-grupal-alpha.git/estructuras"
)

var local *time.Location

func HoraLocal() *time.Location {
	local, _ = time.LoadLocation("America/Argentina/Buenos_Aires")
	return local
}

// CargarDatosVuelos carga los datos de vuelos
func CargarDatosVuelos(filepath string) (map[string]*estructuras.Vuelo, error) {
	archivo, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error al abrir el archivo de vuelos: %w", err)
	}
	defer archivo.Close() //se usa porque antes hicimos open, para evitar problemas futuros

	vuelos := make(map[string]*estructuras.Vuelo)

	scanner := bufio.NewScanner(archivo)
	salteaEncabezado := false
	for scanner.Scan() {
		if !salteaEncabezado { //sirve para no leer la primeera linea de encabezados
			salteaEncabezado = true
			continue
		}

		linea := scanner.Text()
		partes := strings.Split(linea, ";")

		fechaHoraString := partes[1]
		fechaHora, _ := time.ParseInLocation("2006-01-02 15:04:05", fechaHoraString, HoraLocal()) // AAAA-MM-DD HH:MM:SS

		//creacion de la estructura Vuelo
		vuelo := estructuras.NewVuelo()
		vuelo.NumVuelo = partes[0]
		vuelo.FechaHoraProgramada = fechaHora
		vuelo.Destino = partes[2]
		vuelo.AeronaveAsignada = partes[3]
		vuelo.Estado = "Programado"

		vuelos[vuelo.NumVuelo] = vuelo
	}

	configs, _ := CargarDatosConfiguracionAsientos("data/configuracion_asientos.txt")

	for _, vuelo := range vuelos {
		if cfg, ok := configs[vuelo.AeronaveAsignada]; ok {
			vuelo.ConfiguracionAsientos = cfg
		}
	}

	fmt.Printf("%d vuelos cargados.\n", len(vuelos))
	return vuelos, nil
}

// CargarDatosAeronaves carga los datos de las aeronaves
func CargarDatosAeronaves(filepath string) (map[string]*estructuras.Aeronave, error) {

	archivo, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error al abrir el archivo de aeronaves: %w", err)
	}
	defer archivo.Close()

	aeronaves := make(map[string]*estructuras.Aeronave)

	scanner := bufio.NewScanner(archivo)
	salteaEncabezado := false
	for scanner.Scan() {
		if !salteaEncabezado {
			salteaEncabezado = true
			continue
		}

		linea := scanner.Text()
		partes := strings.Split(linea, ";")

		asientos, _ := strconv.Atoi(partes[1])
		capacidadCarga, _ := strconv.ParseFloat(partes[2], 64)
		volumenCarga, _ := strconv.ParseFloat(partes[3], 64)

		//creacion de la estructura Aeronave
		aeronave := estructuras.NewAeronave()
		aeronave.Matricula = partes[0]
		aeronave.Asientos = asientos
		aeronave.CapacidadCarga = capacidadCarga
		aeronave.VolumenCarga = volumenCarga

		aeronaves[aeronave.Matricula] = aeronave
	}

	fmt.Printf("%d aeronaves cargadas.\n", len(aeronaves))
	return aeronaves, nil
}

// CargarDatosClientes carga los datos de los clientes (pasajeros)
func CargarDatosClientes(filepath string) (map[string]*estructuras.Pasajero, error) {
	archivo, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error al abrir el archivo de clientes: %w", err)
	}
	defer archivo.Close()

	clientes := make(map[string]*estructuras.Pasajero)

	scanner := bufio.NewScanner(archivo)
	salteaEncabezado := false
	for scanner.Scan() {
		if !salteaEncabezado {
			salteaEncabezado = true
			continue
		}

		linea := scanner.Text()
		partes := strings.Split(linea, ";")

		// Creación de la estructura Pasajero
		pasajero := estructuras.NewPasajero()
		pasajero.Nombre = partes[0]
		pasajero.Apellido = partes[1]
		pasajero.Documento = partes[2]
		pasajero.Categoria = partes[3]

		clientes[pasajero.Documento] = pasajero
	}

	fmt.Printf("%d clientes cargados.\n", len(clientes))
	return clientes, nil
}

func CargarDatosEquipaje(path string) (map[string]*estructuras.Equipaje, error) {
	archivo, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer archivo.Close()

	equipajes := make(map[string]*estructuras.Equipaje)
	scanner := bufio.NewScanner(archivo)

	salteaEncabezado := false
	for scanner.Scan() {
		if !salteaEncabezado {
			salteaEncabezado = true
			continue
		}

		linea := scanner.Text()
		partes := strings.Split(linea, ";")

		equipaje := estructuras.NewEquipaje()
		equipaje.DocumentoPasajero = partes[0]
		equipaje.Bultos, _ = strconv.Atoi(partes[1])
		equipaje.PesoTotal, _ = strconv.ParseFloat(partes[2], 64)
		equipaje.VolumenTotal, _ = strconv.ParseFloat(partes[3], 64)

		equipajes[equipaje.DocumentoPasajero] = equipaje
	}
	fmt.Printf("%d equipajes cargados.\n", len(equipajes))
	return equipajes, nil
}

func CargarDatosConfiguracionAsientos(filepath string) (map[string][]*estructuras.ConfiguracionAsientos, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error al abrir el archivo de configuración de asientos: %w", err)
	}
	defer file.Close()

	configs := make(map[string][]*estructuras.ConfiguracionAsientos)
	scanner := bufio.NewScanner(file)
	salteaEncabezado := false
	for scanner.Scan() {
		if !salteaEncabezado {
			salteaEncabezado = true
			continue
		}
		linea := scanner.Text()
		partes := strings.Split(linea, ";")

		zona, _ := strconv.Atoi(partes[1])
		inicial, _ := strconv.Atoi(partes[2])
		fin, _ := strconv.Atoi(partes[3])

		configuracion := &estructuras.ConfiguracionAsientos{
			CodAeronave:    partes[0],
			Zona:           zona,
			AsientoInicial: inicial,
			AsientoFinal:   fin,
		}
		configs[partes[0]] = append(configs[partes[0]], configuracion)
	}

	return configs, nil
}

// CargarDatosReservas carga los datos de las reservas y las asigna a los vuelos y pasajeros existentes.
// Los mapas de vuelos y clientes deben ser cargados antes.
func CargarDatosReservas(filepath string) (map[string]*estructuras.Reserva, error) {
	archivo, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error al abrir el archivo de reservas: %w", err)
	}
	defer archivo.Close()

	reservas := make(map[string]*estructuras.Reserva)
	scanner := bufio.NewScanner(archivo)
	salteaEncabezado := false
	for scanner.Scan() {
		if !salteaEncabezado {
			salteaEncabezado = true
			continue
		}

		linea := scanner.Text()
		partes := strings.Split(linea, ";")

		documentoPasajero := partes[0]
		codReserva := partes[1]
		numeroVuelo := partes[2]
		estadoReserva := partes[3]

		pasajero, existe := Clientes[documentoPasajero]
		if !existe {
			continue
		}

		vuelo, existe := Vuelos[numeroVuelo]
		if !existe {
			continue
		}

		pasajero.Vuelo = numeroVuelo
		pasajero.Reserva = codReserva
		pasajero.EstadoReserva = estadoReserva

		reserva, existe := reservas[codReserva]
		if !existe {
			reserva = &estructuras.Reserva{CodReserva: codReserva, Pasajeros: make([]*estructuras.Pasajero, 0, 5)}
			reservas[codReserva] = reserva
		}
		reserva.Pasajeros = append(reserva.Pasajeros, pasajero)

		switch estadoReserva {
		case "Confirmada":
			pasajero.Vuelo = numeroVuelo
			vuelo.Pasajeros = append(vuelo.Pasajeros, pasajero)

		case "Lista de espera":
			pasajero.Vuelo = numeroVuelo
			vuelo.ListaEspera = append(vuelo.ListaEspera, pasajero)

		case "Cancelada":
			pasajero.Vuelo = numeroVuelo
		}
	}

	fmt.Printf("%d reservas procesadas. Pasajeros asignados a sus respectivos vuelos.\n", len(reservas))
	return reservas, nil
}

// CargarAeropuertos carga los datos de los aeropuertos.
func CargarDatosAeropuertos(filepath string) (map[string]*estructuras.Aeropuerto, error) {
	archivo, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error al abrir el archivo de aeropuertos: %w", err)
	}
	defer archivo.Close()

	aeropuertos := make(map[string]*estructuras.Aeropuerto)

	scanner := bufio.NewScanner(archivo)
	salteaEncabezado := false
	for scanner.Scan() {
		if !salteaEncabezado {
			salteaEncabezado = true
			continue
		}

		linea := scanner.Text()
		partes := strings.Split(linea, ";")

		//creacion de la estructura Aeropuerto
		aeropuerto := estructuras.NewAeropuerto()
		aeropuerto.Provincia = partes[0]
		aeropuerto.Ciudad = partes[1]
		aeropuerto.Nombre = partes[2]
		aeropuerto.Cod_iata = partes[3]

		aeropuertos[aeropuerto.Cod_iata] = aeropuerto
	}

	fmt.Printf("%d aeropuertos cargados.\n", len(aeropuertos))
	return aeropuertos, nil
}

// CargarCargas carga los datos de las cargas.
func CargarDatosCargas(filepath string) ([]*estructuras.Carga, error) {
	archivo, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error al abrir el archivo de cargas: %w", err)
	}
	defer archivo.Close()

	var cargas []*estructuras.Carga

	scanner := bufio.NewScanner(archivo)
	salteaEncabezado := false
	for scanner.Scan() {
		if !salteaEncabezado {
			salteaEncabezado = true
			continue
		}

		linea := scanner.Text()
		partes := strings.Split(linea, ";")

		peso, _ := strconv.Atoi(partes[1])
		volumen, _ := strconv.ParseFloat(partes[2], 64)

		//creacion de la estructura Carga
		carga := estructuras.NewCarga()
		carga.Destino = partes[0]
		carga.Peso = peso
		carga.Volumen = volumen

		cargas = append(cargas, carga)
	}

	fmt.Printf("%d cargas listas.\n", len(cargas))
	return cargas, nil
}

// CargarEdificios carga los datos de los edificios.
func CargarDatosEdificios(filepath string) ([]*estructuras.Edificio, error) {
	archivo, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error al abrir el archivo de edificios: %w", err)
	}
	defer archivo.Close()

	var edificios []*estructuras.Edificio

	scanner := bufio.NewScanner(archivo)
	salteaEncabezado := false
	for scanner.Scan() {
		if !salteaEncabezado {
			salteaEncabezado = true
			continue
		}

		linea := scanner.Text()
		partes := strings.Split(linea, ";")

		xi, _ := strconv.Atoi(partes[0])
		altura, _ := strconv.Atoi(partes[1])
		xf, _ := strconv.Atoi(partes[2])

		//creacion de la estructura edificio
		edificio := estructuras.NewEdificio()
		edificio.Xi = xi
		edificio.Altura = altura
		edificio.Xf = xf

		edificios = append(edificios, edificio)
	}

	fmt.Printf("%d edificios cargados.\n", len(edificios))
	return edificios, nil
}
