package estructuras

//estructura que representa una aeronave
type Aeronave struct {
	Matricula      string
	Asientos       int
	CapacidadCarga float64
	VolumenCarga   float64
}

func NewAeronave() *Aeronave {
	return &Aeronave{}
}
