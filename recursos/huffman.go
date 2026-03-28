package recursos

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"

	"github.com/untref-ayp2-guias/trabajo-pr-ctico-grupal-alpha.git/estructuras"
)

// Nodo representa un nodo en el árbol de Huffman.
type Nodo struct {
	Valor      rune
	Frecuencia int
	Izq, Der   *Nodo
}

//cp = ColaPrioridad

// Cola de prioridad (implementa heap.Interface)
type ColaPrioridad []*Nodo

func (cp ColaPrioridad) Len() int {
	return len(cp) // longitud de la cola de prioridad
}

// compara dos nodos por frecuencia
// Less() comparador, devuelve true si el nodo i tiene menor frecuencia que el nodo j
func (cp ColaPrioridad) Less(i, j int) bool {
	// el de menor frecuencia tiene mayor prioridad
	return cp[i].Frecuencia < cp[j].Frecuencia
}

// Swap() intercambia los nodos
func (cp ColaPrioridad) Swap(i, j int) {
	cp[i], cp[j] = cp[j], cp[i]
}
func (cp *ColaPrioridad) Push(x any) {
	*cp = append(*cp, x.(*Nodo)) // agrega un nuevo nodo a la cola de prioridad
}

// elimina y devuelve el nodo con menor frecuencia
func (cp *ColaPrioridad) Pop() any {
	n := len(*cp)
	x := (*cp)[n-1]
	*cp = (*cp)[:n-1]
	return x
}

// Construye el árbol de Huffman a partir de una frecuencia de caracteres
func ConstruirArbol(frecuencias map[rune]int) *Nodo {
	cp := &ColaPrioridad{}
	for r, f := range frecuencias {
		heap.Push(cp, &Nodo{Valor: r, Frecuencia: f}) // Crear un nodo para cada carácter y su frecuencia
	}
	heap.Init(cp) // Inicializar la cola de prioridad

	// mientras haya más de un nodo en la cola
	// saca a los dos nodos con menor frecuencia
	for cp.Len() > 1 {
		izq := heap.Pop(cp).(*Nodo)
		der := heap.Pop(cp).(*Nodo)
		nuevo := &Nodo{ // crear un nuevo nodo que combina los dos nodos
			Frecuencia: izq.Frecuencia + der.Frecuencia,
			Izq:        izq, // asignar el nodo izquierdo
			Der:        der, // asignar el nodo derecho
		}
		heap.Push(cp, nuevo) // agregar el nuevo nodo a la cola de prioridad
	}
	return heap.Pop(cp).(*Nodo) // devolver el nodo raíz del árbol de Huffman
}

// Genera el código Huffman para cada carácter
func GenerarCodigos(n *Nodo, prefijo string, codigos map[rune]string) {
	if n == nil {
		return
	}
	if n.Izq == nil && n.Der == nil { // si es una hoja, asigna el código al carácter
		// n.Valor es el carácter y prefijo es el código binario
		codigos[n.Valor] = prefijo
	}
	GenerarCodigos(n.Izq, prefijo+"0", codigos) // llama recursivamente al hijo izquierdo con prefijo "0"
	GenerarCodigos(n.Der, prefijo+"1", codigos) // llama recursivamente al hijo derecho con prefijo "1"
}

// Codifica un texto usando el árbol de Huffman
// agrega el código correspondiente al carácter
// si el carácter no está en el mapa de códigos, se ignora
func Codificar(texto string, codigos map[rune]string) string {
	codificado := ""
	for _, r := range texto { // recorre cada carácter del texto
		codificado += codigos[r]
	}
	return codificado
}

// Decodifica una cadena binaria usando el árbol de Huffman
func Decodificar(codificado string, raiz *Nodo) string {
	resultado := ""
	nodo := raiz
	for _, bit := range codificado { // recorre cada bit del código
		if bit == '0' {
			nodo = nodo.Izq
		} else {
			nodo = nodo.Der
		}
		if nodo.Izq == nil && nodo.Der == nil { // si es una hoja, agrega el valor al resultado
			resultado += string(nodo.Valor)
			nodo = raiz // reinicia al nodo raíz
		}
	}
	return resultado
}

// CodificarHuffman genera el archivo .huff a partir de un archivo de texto
func CodificarHuffman(inputPath, outputPath string) error {
	inputData, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("error al leer el archivo: %w", err)
	}
	texto := string(inputData)

	frecuencias := make(map[rune]int)
	for _, r := range texto {
		frecuencias[r]++
	}

	raiz := ConstruirArbol(frecuencias)

	codigos := make(map[rune]string)
	GenerarCodigos(raiz, "", codigos)

	codificado := Codificar(texto, codigos)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error al crear el archivo comprimido: %w", err)
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)

	for r, c := range codigos {
		fmt.Fprintf(writer, "%d|%s\n", r, c)
	}
	fmt.Fprintln(writer, "==")
	fmt.Fprintln(writer, codificado)

	return writer.Flush()
}

func ComprimirConHuffman(vuelo *estructuras.Vuelo) error {
	carpeta := "output"
	err := os.MkdirAll(carpeta, 0755)
	if err != nil {
		return fmt.Errorf("error al crear la carpeta '%s': %w", carpeta, err)
	}

	input := fmt.Sprintf("output/vuelo_%s.txt", vuelo.NumVuelo)
	output := fmt.Sprintf("output/vuelo_%s.huff", vuelo.NumVuelo)
	err = CodificarHuffman(input, output)
	if err != nil {
		fmt.Println("Error al generar archivo .huff:", err)
	}
	err = os.Remove(input)
	if err != nil {
		fmt.Println("No se pudo borrar el archivo .txt original:", err)
	}
	return nil
}

func DescomprimirArchivosHuffman(inputPath, outputPath string) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("error al abrir el archivo .huff: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	codigos := make(map[string]rune)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "==" {
			break
		}
		var r rune
		var codigo string
		fmt.Sscanf(line, "%d|%s", &r, &codigo)
		codigos[codigo] = r
	}

	// Leer el texto codificado
	var codificado string
	for scanner.Scan() {
		codificado += scanner.Text()
	}

	// Invertir mapa: código binario -> carácter
	type nodo struct {
		Izq, Der *nodo
		Valor    *rune
	}
	raiz := &nodo{}

	for bin, r := range codigos {
		actual := raiz
		for _, bit := range bin {
			if bit == '0' {
				if actual.Izq == nil {
					actual.Izq = &nodo{}
				}
				actual = actual.Izq
			} else {
				if actual.Der == nil {
					actual.Der = &nodo{}
				}
				actual = actual.Der
			}
		}
		actual.Valor = new(rune)
		*actual.Valor = r
	}

	// Decodificar
	var resultado string
	actual := raiz
	for _, bit := range codificado {
		if bit == '0' {
			actual = actual.Izq
		} else {
			actual = actual.Der
		}
		if actual.Valor != nil {
			resultado += string(*actual.Valor)
			actual = raiz
		}
	}

	// Escribir el resultado en archivo .txt
	err = os.WriteFile(outputPath, []byte(resultado), 0644)
	if err != nil {
		return fmt.Errorf("error al generar archivo archivo .txt: %w", err)
	}

	return nil
}
