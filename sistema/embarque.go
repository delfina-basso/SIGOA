package sistema

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/untref-ayp2-guias/trabajo-pr-ctico-grupal-alpha.git/datos"
	"github.com/untref-ayp2-guias/trabajo-pr-ctico-grupal-alpha.git/estructuras"
)

func GenerarLlegadasEmbarque(vuelo *estructuras.Vuelo, inicioEmbarque time.Time, onSuccess func(tiempo time.Time, pasajero *estructuras.Pasajero)) {

	var zona1, zona2, zona3, conCategoria []*estructuras.Pasajero

	for _, p := range vuelo.Pasajeros {
		if p.PerdioVuelo || !p.HizoCheckIn {
			continue
		}

		if p.Categoria != "" {
			conCategoria = append(conCategoria, p)
		} else {
			switch p.ZonaEmbarque {
			case "1":
				zona1 = append(zona1, p)
			case "2":
				zona2 = append(zona2, p)
			case "3":
				zona3 = append(zona3, p)
			}
		}
	}

	finEmbarque := inicioEmbarque.Add(30 * time.Minute)

	for _, p := range conCategoria {
		ventana := int(finEmbarque.Sub(inicioEmbarque).Minutes())
		llegada := inicioEmbarque.Add(time.Duration(rand.Intn(ventana)) * time.Minute)
		p.HoraLlegadaEmbarque = llegada
		onSuccess(llegada, p)
	}

	for _, p := range zona1 {
		inicioZona := inicioEmbarque
		llegada := inicioZona.Add(time.Duration(rand.Intn(15)) * time.Minute)
		p.HoraLlegadaEmbarque = llegada
		onSuccess(llegada, p)
	}

	for _, p := range zona2 {
		inicioZona := inicioEmbarque.Add(10 * time.Minute)
		llegada := inicioZona.Add(time.Duration(rand.Intn(15)) * time.Minute)
		p.HoraLlegadaEmbarque = llegada
		onSuccess(llegada, p)
	}

	for _, p := range zona3 {
		inicioZona := inicioEmbarque.Add(20 * time.Minute)
		llegada := inicioZona.Add(time.Duration(rand.Intn(15)) * time.Minute)
		p.HoraLlegadaEmbarque = llegada
		onSuccess(llegada, p)
	}
}

// GenerarArchivoVuelo genera un archivo temporal con todos los datos repsectivos al vuelo.
func GenerarArchivoVuelo(vuelo *estructuras.Vuelo, tiempo time.Time,
	listaEmbarcados []*estructuras.Pasajero, listaEquipaje []*estructuras.Equipaje,
	listaNoPresentados []*estructuras.Pasajero, listaEspera []*estructuras.Pasajero,
	listaCargas []*estructuras.Carga) error {

	carpeta := "output"
	err := os.MkdirAll(carpeta, 0755) //no permite que se escriba en ella
	if err != nil {
		return fmt.Errorf("error al crear la carpeta '%s': %w", carpeta, err)
	}

	// Crear archivo de salida
	nombreArchivo := filepath.Join(carpeta, fmt.Sprintf("vuelo_%s.txt", vuelo.NumVuelo))
	file, err := os.Create(nombreArchivo)
	if err != nil {
		return fmt.Errorf("error al crear archivo de vuelo: %w", err)
	}
	defer file.Close()

	// Datos del vuelo
	fmt.Fprintf(file, "INFORME DEL VUELO %s\n", vuelo.NumVuelo)
	fmt.Fprintf(file, "\nFecha y hora programada: %s\n", vuelo.FechaHoraProgramada.Format("2006-01-02 15:04"))
	fmt.Fprintf(file, "Destino: %s, %s.\n", vuelo.Destino, datos.Aeropuertos[vuelo.Destino].Provincia)
	fmt.Fprintf(file, "Aeronave: %s \n", vuelo.AeronaveAsignada)

	// Pasajeros embarcados
	fmt.Fprintln(file, "\nPasajeros embarcados:")
	if len(listaEmbarcados) == 0 {
		fmt.Fprintf(file, "Ninguno.\n")
	} else {
		for _, p := range listaEmbarcados {
			fmt.Fprintf(file, " - %s %s (%s) Categoría: %s\n", p.Nombre, p.Apellido, p.Documento, p.Categoria)
		}
	}

	// Equipaje
	fmt.Fprintln(file, "\nEquipaje de los pasajeros embarcados:")
	if len(listaEquipaje) == 0 {
		fmt.Fprintf(file, "Ninguno.\n")
	} else {
		for _, e := range listaEquipaje {
			pasajero := datos.Clientes[e.DocumentoPasajero]
			fmt.Fprintf(file, " - %s %s (%s) | Bultos: %d | Peso: %.2f kg\n",
				pasajero.Nombre, pasajero.Apellido, e.DocumentoPasajero, e.Bultos, e.PesoTotal)
		}
	}

	// No presentados
	fmt.Fprintln(file, "\nPasajeros que no se presentaron:")
	if len(listaNoPresentados) == 0 {
		fmt.Fprintf(file, "Ninguno.\n")
	} else {
		for _, p := range listaNoPresentados {
			fmt.Fprintf(file, " - %s %s (%s)\n", p.Nombre, p.Apellido, p.Documento)
		}
	}

	// Lista de espera
	fmt.Fprintln(file, "\nPasajeros restantes en lista de espera:")
	if len(listaEspera) == 0 {
		fmt.Fprintf(file, "Ninguno.\n")
	} else {
		for _, p := range listaEspera {
			fmt.Fprintf(file, " - %s %s (%s)\n", p.Nombre, p.Apellido, p.Documento)
		}
	}

	// Cargas
	fmt.Fprintln(file, "\nLista de cargas:")
	if len(vuelo.Carga) == 0 {
		fmt.Fprintf(file, "Ninguno.\n")
	} else {
		for _, c := range vuelo.Carga {
			fmt.Fprintf(file, " - Destino: %s | Peso: %d kg | Volumen: %.2f \n", c.Destino, c.Peso, c.Volumen)
		}
	}

	fmt.Printf("Vuelo %s: Se generó el archivo de datos del vuelo correctamente.\n", vuelo.NumVuelo)
	return nil
}
