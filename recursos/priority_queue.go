package recursos

type PriorityQueue[T any] struct {
	data    []T
	compare func(a, b T) bool // true si a tiene mayor prioridad que b (ej: max-heap: a > b)
}

func NewPriorityQueue[T any](compare func(a, b T) bool) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		data:    make([]T, 0),
		compare: compare,
	}
}

func (pq *PriorityQueue[T]) Len() int {
	return len(pq.data)
}

func (pq *PriorityQueue[T]) Empty() bool {
	return pq.Len() == 0
}

func (pq *PriorityQueue[T]) Peek() T {
	return pq.data[0]
}

func (pq *PriorityQueue[T]) Push(valor T) {
	pq.data = append(pq.data, valor)
	pq.upHeap(pq.Len() - 1)
}

func (pq *PriorityQueue[T]) Pop() T {
	n := pq.Len() - 1
	pq.swap(0, n)
	valor := pq.data[n]
	pq.data = pq.data[:n]
	pq.downHeap(0)
	return valor
}

func (pq *PriorityQueue[T]) swap(i, j int) {
	pq.data[i], pq.data[j] = pq.data[j], pq.data[i]
}

func (pq *PriorityQueue[T]) upHeap(i int) {
	for i > 0 {
		padre := (i - 1) / 2
		if pq.compare(pq.data[i], pq.data[padre]) {
			pq.swap(i, padre)
			i = padre
		} else {
			break
		}
	}
}

func (pq *PriorityQueue[T]) downHeap(i int) {
	n := pq.Len()
	for {
		izq := 2*i + 1
		der := 2*i + 2
		mayor := i

		if izq < n && pq.compare(pq.data[izq], pq.data[mayor]) {
			mayor = izq
		}
		if der < n && pq.compare(pq.data[der], pq.data[mayor]) {
			mayor = der
		}

		if mayor == i {
			break
		}
		pq.swap(i, mayor)
		i = mayor
	}
}
