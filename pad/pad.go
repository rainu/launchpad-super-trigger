package pad

type Lighter interface {
	Light(x, y, g, r int) error
}
