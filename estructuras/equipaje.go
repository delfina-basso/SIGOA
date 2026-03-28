package estructuras

type Equipaje struct {
	DocumentoPasajero string
	NumeroVuelo       string
	Bultos            int
	PesoTotal         float64
	VolumenTotal      float64
}

func NewEquipaje() *Equipaje {
	return &Equipaje{}
}
