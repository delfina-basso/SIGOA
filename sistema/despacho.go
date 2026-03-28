package sistema

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/untref-ayp2-guias/trabajo-pr-ctico-grupal-alpha.git/estructuras"
	"github.com/untref-ayp2-guias/trabajo-pr-ctico-grupal-alpha.git/recursos"
)

var ColaDespacho = recursos.NewPriorityQueue[*estructuras.Vuelo](CompararVuelos)

// CompararVuelos define la prioridad para la cola de despacho
func CompararVuelos(a, b *estructuras.Vuelo) bool {
	if a.Estado == "Listo para despegar" && b.Estado != "Listo para despegar" {
		return true
	}
	if b.Estado == "Listo para despegar" && a.Estado != "Listo para despegar" {
		return false
	}
	return a.FechaHoraProgramada.Before(b.FechaHoraProgramada)
}

type BordeEdificio struct {
	X      int
	Altura int
	Tipo   int // 1 = inicio, -1 = fin
}

type puntoHorizonte struct {
	X int
	y int
}

// CalcularLineaHorizonte calcula los puntos de la línea del horizonte y los devuelve.
func CalcularLineaHorizonte(edificios []*estructuras.Edificio) []puntoHorizonte {
	var bordes []BordeEdificio

	//crea dos bordes por edificio, uno donde empieza y otro donde termina
	for _, edificio := range edificios {
		bordes = append(bordes, BordeEdificio{X: edificio.Xi, Altura: edificio.Altura, Tipo: 1})  // inicio
		bordes = append(bordes, BordeEdificio{X: edificio.Xf, Altura: edificio.Altura, Tipo: -1}) // fin
	}

	// Ordenar bordes
	sort.Slice(bordes, func(i, j int) bool {
		//por X
		if bordes[i].X != bordes[j].X {
			return bordes[i].X < bordes[j].X
		}
		//por inicio antes que fin, si un edificio empieza en el mismo x que otro termina, primero va el edificio que empieza
		if bordes[i].Tipo != bordes[j].Tipo {
			return bordes[i].Tipo > bordes[j].Tipo
		}
		//por altura descendente, si dos edificios empiezan o terminan ambos en el mismo x, primero va el mas alto
		return bordes[i].Altura > bordes[j].Altura
	})

	//crea la priority queue
	alturas := recursos.NewPriorityQueue(func(a, b int) bool {
		return a > b // max-heap
	})

	//se agrega el 0 para que sea la primer altura maxima de la priority queue
	alturas.Push(0)

	//esto es la forma de eliminar las alturas ya utilizadas porque la pq al ser un heap no puede eliminar elementos cualquiera
	contador := map[int]int{0: 1} //inicializa la altura 0 como que apareció una vez (altura piso)

	var resultado []puntoHorizonte
	alturaMaximaAnterior := 0

	for _, borde := range bordes {
		//si el borde es un inicio agrega el edificio a la pq
		if borde.Tipo == 1 {
			alturas.Push(borde.Altura)
			contador[borde.Altura]++
		} else { //si el borde es un final quita el edificio de la pq
			contador[borde.Altura]--
		}

		// Elimina los bordes de edificios cuyas alturas ya fueron calculadas (inicio +1, fin -1 = 0)
		// hace Pop de la pq hasta que esté vacía (no haya más bordes de deificios) o hasta que encuentre uno que quede por calcular
		for !alturas.Empty() && contador[alturas.Peek()] == 0 {
			alturas.Pop()
		}

		//la altura maxima luego de agregar o quitar el edificio
		alturaMaximaActual := alturas.Peek()

		//si la altura maxima cambia, el punto es parte del resultado
		if alturaMaximaActual != alturaMaximaAnterior {
			punto := puntoHorizonte{X: borde.X, y: alturaMaximaActual}
			resultado = append(resultado, punto)
			alturaMaximaAnterior = alturaMaximaActual
		}
	}

	return resultado
}

func GenerarArchivoLineaHorizonte(puntos []puntoHorizonte) error {
	carpeta := "output"
	err := os.MkdirAll(carpeta, 0755) //no permite que se escriba en ella
	if err != nil {
		return fmt.Errorf("error al crear la carpeta '%s': %w", carpeta, err)
	}

	// crear carpeta si no existe
	nombreArchivo := filepath.Join(carpeta, "linea_horizonte.txt")

	file, err := os.Create(nombreArchivo)
	if err != nil {
		return fmt.Errorf("error al crear el archivo de línea del horizonte: %w", err)
	}
	defer file.Close()

	for _, p := range puntos {
		fmt.Fprintf(file, "%d,%d\n", p.X, p.y)
	}

	return nil
}
