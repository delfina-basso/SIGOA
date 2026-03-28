package estructuras

//estructura que representa una carga
type Carga struct {
	Destino string
	Peso    int
	Volumen float64
}

func NewCarga() *Carga {
	return &Carga{}
}
