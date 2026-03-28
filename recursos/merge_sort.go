package recursos

/* MergeSort ordena un slice de cualquier tipo usando la funcion menor(i, j) bool.
 El slice se pasa como []T y la función menor compara elementos en posiciones i y j.
 Ejemplos de uso: vuelosOrdenados := recursos.MergeSort(vuelos, func(a, b Vuelo) bool {
    				return a.HoraSalida.Before(b.HoraSalida)
				  })
				pasajerosOrdenados := recursos.MergeSort(pasajeros, func(a, b Pasajero) bool{
    				return a.Nombre < b.Nombre
				})
*/
func MergeSort[T any](arr []T, menor func(a, b T) bool) []T {
	n := len(arr)
	if n <= 1 {
		return arr
	}
	mid := n / 2
	izq := MergeSort(arr[:mid], menor)
	der := MergeSort(arr[mid:], menor)
	return merge(izq, der, menor)
}

// junta todos los slices en un array ordenado
func merge[T any](izq, der []T, menor func(a, b T) bool) []T {
	result := make([]T, 0, len(izq)+len(der))
	i, j := 0, 0
	for i < len(izq) && j < len(der) {
		if menor(izq[i], der[j]) {
			result = append(result, izq[i])
			i++
		} else {
			result = append(result, der[j])
			j++
		}
	}
	result = append(result, izq[i:]...)
	result = append(result, der[j:]...)
	return result
}
