package estructuras

//estructura que representa un aeropuerto
type Aeropuerto struct {
	Provincia string
	Ciudad    string
	Nombre    string
	Cod_iata  string
}

func NewAeropuerto() *Aeropuerto {
	return &Aeropuerto{}
}
