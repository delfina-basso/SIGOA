package estructuras

import (
	"fmt"
	"time"
)

// estructura que representa un vuelo
type Vuelo struct {
	NumVuelo              string
	FechaHoraProgramada   time.Time
	Destino               string // cod_iata
	AeronaveAsignada      string // matricula
	Estado                string
	CheckInCerrado        bool
	Pasajeros             []*Pasajero
	ListaEspera           []*Pasajero
	ConfiguracionAsientos []*ConfiguracionAsientos
	ZonaHabilitada        int
	Carga                 []*Carga
	Equipaje              []float64 //equipaje[0] es peso total y equipaje[1] es volumen total de equipjaes despachados
}

func NewVuelo() *Vuelo {
	return &Vuelo{
		Pasajeros:             make([]*Pasajero, 0),
		ListaEspera:           make([]*Pasajero, 0),
		ConfiguracionAsientos: make([]*ConfiguracionAsientos, 0),
		Carga:                 make([]*Carga, 0),
		Equipaje:              make([]float64, 2),
	}
}

func (v *Vuelo) String() string {
	return fmt.Sprintf("Vuelo %s a %s %s", v.NumVuelo, v.Destino, v.FechaHoraProgramada.Format("15:04"))
}
