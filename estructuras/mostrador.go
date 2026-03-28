package estructuras

import (
	"fmt"
	"math/rand"
	"time"
)

// MostradorCheckIn representa un mostrador de atención
type MostradorCheckIn struct {
	ID             int
	Disponible     bool
	Activo         bool
	PasajeroActual *Pasajero
	TiempoRestante time.Duration
}

// Atender inicia la atención de un pasajero y devuelve la duración de la atención
func (m *MostradorCheckIn) Atender(p *Pasajero, actual time.Time) time.Duration {
	m.PasajeroActual = p
	m.Disponible = false
	duracion := time.Duration(rand.Intn(5)+1) * time.Minute
	m.TiempoRestante = duracion

	fmt.Printf("[%s] Mostrador %d: Atendiendo a %s %s (%s)\n",
		actual.Format("2006-01-02 15:04"), m.ID, p.Nombre, p.Apellido, p.Documento)

	return duracion
}

// FinalizarAtencionMostrador finaliza la atención
func (m *MostradorCheckIn) FinalizarAtencionMostrador(actual time.Time) {
	p := m.PasajeroActual
	if p.EstadoReserva == "Lista de espera" {
		fmt.Printf("[%s] Mostrador %d: %s %s en lista de espera para el vuelo %s.\n",
			actual.Format("2006-01-02 15:04"), m.ID, p.Nombre, p.Apellido, p.Vuelo)
	}
	p.HizoCheckIn = true
	m.PasajeroActual = nil
	m.Disponible = true
}
