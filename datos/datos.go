package datos

import (
	"fmt"

	"github.com/untref-ayp2-guias/trabajo-pr-ctico-grupal-alpha.git/estructuras"
)

var (
	Vuelos                map[string]*estructuras.Vuelo
	Aeronaves             map[string]*estructuras.Aeronave
	Clientes              map[string]*estructuras.Pasajero
	Reservas              map[string]*estructuras.Reserva
	Aeropuertos           map[string]*estructuras.Aeropuerto
	Equipajes             map[string]*estructuras.Equipaje
	Cargas                []*estructuras.Carga
	Edificios             []*estructuras.Edificio
	ConfiguracionAsientos map[string][]*estructuras.ConfiguracionAsientos
)

func InicializarDatos(escenario string) error {
	Vuelos = make(map[string]*estructuras.Vuelo)
	Aeronaves = make(map[string]*estructuras.Aeronave)
	Clientes = make(map[string]*estructuras.Pasajero)
	Reservas = make(map[string]*estructuras.Reserva)
	Aeropuertos = make(map[string]*estructuras.Aeropuerto)
	Equipajes = make(map[string]*estructuras.Equipaje)
	Cargas = []*estructuras.Carga{}
	Edificios = []*estructuras.Edificio{}
	ConfiguracionAsientos = make(map[string][]*estructuras.ConfiguracionAsientos)

	var err error

	vuelosPath := fmt.Sprintf("data/%s/vuelos.txt", escenario)
	Vuelos, err = CargarDatosVuelos(vuelosPath)
	if err != nil {
		return fmt.Errorf("error cargando vuelos: %w", err)
	}

	Aeronaves, err = CargarDatosAeronaves("data/aeronaves.txt")
	if err != nil {
		return fmt.Errorf("error cargando aeronaves: %w", err)
	}

	clientesPath := fmt.Sprintf("data/%s/clientes.txt", escenario)
	Clientes, err = CargarDatosClientes(clientesPath)
	if err != nil {
		return fmt.Errorf("error cargando clientes: %w", err)
	}

	Aeropuertos, err = CargarDatosAeropuertos("data/aeropuertos.txt")
	if err != nil {
		return fmt.Errorf("error cargando aeropuertos: %w", err)
	}

	equipajePath := fmt.Sprintf("data/%s/equipaje.txt", escenario)
	Equipajes, err = CargarDatosEquipaje(equipajePath)
	if err != nil {
		return fmt.Errorf("error cargando equipajes: %w", err)
	}

	cargasPath := fmt.Sprintf("data/%s/cargas.txt", escenario)
	Cargas, err = CargarDatosCargas(cargasPath)
	if err != nil {
		return fmt.Errorf("error cargando cargas: %w", err)
	}

	Edificios, err = CargarDatosEdificios("data/edificios.txt")
	if err != nil {
		return fmt.Errorf("error cargando edificios: %w", err)
	}

	ConfiguracionAsientos, err = CargarDatosConfiguracionAsientos("data/configuracion_asientos.txt")
	if err != nil {
		return fmt.Errorf("error cargando configuración de asientos: %w", err)
	}

	reservasPath := fmt.Sprintf("data/%s/reservas.txt", escenario)

	Reservas, err = CargarDatosReservas(reservasPath)
	if err != nil {
		return fmt.Errorf("error cargando reservas: %w", err)
	}

	fmt.Println("Todos los datos del escenario fueron cargados correctamente.")
	return nil
}
