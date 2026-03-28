package estructuras

//estructura que representa un edificio
type Edificio struct {
	Xi     int
	Altura int
	Xf     int
}

func NewEdificio() *Edificio {
	return &Edificio{}
}
