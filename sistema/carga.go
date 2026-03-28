package sistema

import (
	"fmt"
	"time"

	"github.com/untref-ayp2-guias/trabajo-pr-ctico-grupal-alpha.git/datos"
	"github.com/untref-ayp2-guias/trabajo-pr-ctico-grupal-alpha.git/estructuras"
	"github.com/untref-ayp2-guias/trabajo-pr-ctico-grupal-alpha.git/recursos"
)

// Carga representa el proceso de carga de un vuelo
type Carga struct {
	Estado          string    // "pendiente", "en curso", "finalizada"
	HoraCheckInFin  time.Time // hora en que finaliza el check-in
	HoraVuelo       time.Time // hora de salida del vuelo
	HoraInicioCarga time.Time // calculada: misma que HoraCheckInFin
	HoraLimiteCarga time.Time // calculada: 1 hora antes del vuelo
}

// NuevaCarga crea una instancia del proceso de carga
func NuevaCarga(horaCheckInFin, horaVuelo time.Time) *Carga {
	return &Carga{
		Estado:          "pendiente",
		HoraCheckInFin:  horaCheckInFin,
		HoraVuelo:       horaVuelo,
		HoraInicioCarga: horaCheckInFin,
		HoraLimiteCarga: horaVuelo.Add(-1 * time.Hour),
	}
}

// ActualizarEstado actualiza el estado del proceso de carga según la hora actual
func (c *Carga) ActualizarEstado(ahora time.Time) {
	switch {
	case ahora.Before(c.HoraInicioCarga):
		c.Estado = "pendiente"
	case ahora.After(c.HoraInicioCarga) && ahora.Before(c.HoraLimiteCarga):
		c.Estado = "en curso"
	default:
		c.Estado = "finalizada"
	}
}

// MostrarEstado imprime el estado actual
func (c *Carga) MostrarEstado() {
	fmt.Printf("Estado de carga: %s\n", c.Estado)
}

// BuscarCargaPorDestino devuelve todas las cargas para un destino específico
func BuscarCargaPorDestino(cargas []*estructuras.Carga, destino string) []estructuras.Carga {
	var resultado []estructuras.Carga
	for _, carga := range cargas {
		if carga.Destino == destino {
			resultado = append(resultado, *carga)
			fmt.Printf("Destino: %s. Peso: %d. Volumen: %.2f\n",
				carga.Destino, carga.Peso, carga.Volumen)
		}
	}
	return resultado
}

// asignarCarga asigna las cargas al vuelo teniendo en cuenta los equipajes
func AsignarCarga(vuelo *estructuras.Vuelo, cargas []*estructuras.Carga) []*estructuras.Carga {
	aeronave := datos.Aeronaves[vuelo.AeronaveAsignada]
	var asignadas []*estructuras.Carga
	pesoActual := vuelo.Equipaje[0]
	volumenActual := vuelo.Equipaje[1]

	// filtramos cargas por destino
	var cargasDestino []*estructuras.Carga
	for _, carga := range cargas {
		if carga.Destino == vuelo.Destino {
			cargasDestino = append(cargasDestino, carga)
		}
	}

	// ordenamos por peso descendente
	cargasDestino = recursos.MergeSort(cargasDestino, func(a, b *estructuras.Carga) bool {
		return a.Peso > b.Peso
	})

	// asignamos cargas sin exceder capacidad
	for _, carga := range cargasDestino {
		nuevoPeso := pesoActual + float64(carga.Peso)
		nuevoVolumen := volumenActual + carga.Volumen

		if nuevoPeso <= aeronave.CapacidadCarga && nuevoVolumen <= aeronave.VolumenCarga {
			asignadas = append(asignadas, carga)
			pesoActual = nuevoPeso
			volumenActual = nuevoVolumen
		}
	}
	// eliminamos cargas asignadas del slice global
	for _, asignada := range asignadas {
		for i := 0; i < len(datos.Cargas); i++ {
			if datos.Cargas[i] == asignada {
				datos.Cargas = append(datos.Cargas[:i], datos.Cargas[i+1:]...)
				i--
				break
			}
		}
	}
	vuelo.Carga = asignadas
	return asignadas
}
