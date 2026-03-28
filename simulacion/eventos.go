package simulacion

import (
	"time"

	"github.com/untref-ayp2-guias/trabajo-pr-ctico-grupal-alpha.git/estructuras"
)

type TipoEvento string

const (
	EventoLlegadaPasajero TipoEvento = "LlegadaPasajero"
	EventoInicioCheckIn   TipoEvento = "InicioCheckIn" //habilitacion del checkin del vuelo
	EventoCierreCheckIn   TipoEvento = "CierreCheckIn" // cierre del checkin del vuelo
	EventoFinCheckIn      TipoEvento = "FinCheckIn"    //fin del checkin de los pasajeros

	EventoInicioCarga TipoEvento = "InicioCarga"
	EventoFinCarga    TipoEvento = "FinCarga"

	EventoLlegadaEmbarque TipoEvento = "LlegadaEmbarque" //llegadas de pasajeros al embarque
	EventoInicioEmbarque  TipoEvento = "InicioEmbarque"  //inicio de la zona 1
	EventoFinZonaEmbarque TipoEvento = "FinZonaEmbarque" //fin de zonas 1 y 2
	EventoFinEmbarque     TipoEvento = "FinEmbarque"     // fin de zona 3 y embarque

	EventoDespachoVuelo TipoEvento = "DespachoVuelo"
)

type Evento struct {
	Tiempo time.Time
	Tipo   TipoEvento
	Data   any
}

type ColaEventos []*Evento

func (ce ColaEventos) Len() int {
	return len(ce)
}

func (ce ColaEventos) Less(i, j int) bool {
	return ce[i].Tiempo.Before(ce[j].Tiempo)
}

func (ce ColaEventos) Swap(i, j int) {
	ce[i], ce[j] = ce[j], ce[i]
}

func (ce *ColaEventos) Push(x any) {
	*ce = append(*ce, x.(*Evento))
}

func (ce *ColaEventos) Pop() any {
	old := *ce
	n := len(old)
	elem := old[n-1]
	*ce = old[:n-1]
	return elem
}

func GenerarEventosIniciales(sim *Simulador, vuelos map[string]*estructuras.Vuelo) {
	for _, vuelo := range vuelos {
		horaVuelo := vuelo.FechaHoraProgramada
		sim.AgregarEvento(horaVuelo.Add(-2*time.Hour), EventoInicioCheckIn, vuelo)
	}
}
