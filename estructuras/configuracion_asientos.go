package estructuras

import (
	"fmt"
	"strconv"
)

// estructura que representa la configuracion de asientos de una aeronave
type ConfiguracionAsientos struct {
	CodAeronave    string
	Zona           int
	AsientoInicial int
	AsientoFinal   int
}

func NewConfiguracionAsientos() *ConfiguracionAsientos {
	return &ConfiguracionAsientos{}
}

func AsignarZonasEmbarque(vuelos map[string]*Vuelo, configAsientos map[string][]*ConfiguracionAsientos) {
	for _, vuelo := range vuelos {
		configs, ok := configAsientos[vuelo.AeronaveAsignada]
		if !ok {
			fmt.Printf("No hay configuración de zonas para la aeronave %s\n", vuelo.AeronaveAsignada)
			continue
		}

		for _, p := range vuelo.Pasajeros {
			asientoNum, err := strconv.Atoi(p.AsientoAsignado)
			if err != nil {
				continue
			}
			for _, cfg := range configs {
				if asientoNum >= cfg.AsientoInicial && asientoNum <= cfg.AsientoFinal {
					p.ZonaEmbarque = fmt.Sprintf("%d", cfg.Zona)
					break
				}
			}
		}
	}
}
