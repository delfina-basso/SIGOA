package estructuras

import (
	"fmt"
	"time"
)

// estructura que representa un pasajero
type Pasajero struct {
	Documento       string
	Nombre          string
	Apellido        string
	Categoria       string
	Reserva         string
	EstadoReserva   string
	Vuelo           string
	AsientoAsignado string
	ZonaEmbarque    string
	PerdioVuelo     bool
	//checkin
	HoraLlegada time.Time
	YaLlego     bool
	HizoCheckIn bool
	//embarque
	HoraLlegadaEmbarque time.Time
	Embarcado           bool
}

func NewPasajero() *Pasajero {
	return &Pasajero{}
}

func BuscarPasajeroPorDocumento(Pasajeros map[string]*Pasajero, documento string) string {
	if pasajero, ok := Pasajeros[documento]; ok {
		categoria := pasajero.Categoria
		if pasajero.Categoria == "" {
			categoria = "No frecuente"
		}
		return fmt.Sprintf("Documento: %s. Nombre: %s. Apellido: %s. Categoria: %s. Reserva: %s. Vuelo: %s.",
			pasajero.Documento, pasajero.Nombre, pasajero.Apellido, categoria, pasajero.Reserva, pasajero.Vuelo)
	} else {
		return fmt.Sprintf("No se encontró un pasajero con el documento: %s", documento)
	}
}

func (p *Pasajero) String() string {
	return fmt.Sprintf("%s %s (%s)", p.Nombre, p.Apellido, p.Documento)
}

// Estructura para dar prioridad
type PasajeroPrioritario struct {
	Pasajero *Pasajero
	Indice   int
}

// Mapa de prioridad de categorías
var prioridadCategoria = map[string]int{
	"Platino": 1,
	"Oro":     2,
	"Plata":   3,
	"":        4,
}

type ColaPrioridad []*PasajeroPrioritario

func (cp ColaPrioridad) Len() int {
	return len(cp)
}

func (cp ColaPrioridad) Less(i, j int) bool {
	p1 := cp[i].Pasajero
	p2 := cp[j].Pasajero

	pri1 := prioridadCategoria[p1.Categoria]
	pri2 := prioridadCategoria[p2.Categoria]

	if pri1 != pri2 {
		return pri1 < pri2
	}
	return p1.HoraLlegada.Before(p2.HoraLlegada)
}

func (cp ColaPrioridad) Swap(i, j int) {
	cp[i], cp[j] = cp[j], cp[i]
	cp[i].Indice = i
	cp[j].Indice = j
}

func (cp *ColaPrioridad) Push(x any) {
	n := len(*cp)
	item := x.(*PasajeroPrioritario)
	item.Indice = n
	*cp = append(*cp, item)
}

func (cp *ColaPrioridad) Pop() any {
	old := *cp
	n := len(old)
	item := old[n-1]
	*cp = old[0 : n-1]
	return item
}
