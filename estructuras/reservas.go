package estructuras

//estructura que representa una reserva
type Reserva struct {
	CodReserva        string
	Pasajeros         []*Pasajero
	DocumentoPasajero string
	NumVuelo          string
	EstadoReserva     string
}

func NewReserva() *Reserva {
	return &Reserva{}
}
