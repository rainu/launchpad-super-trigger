package gfx

type coord struct {
	X int
	Y int
}

type Quadrant int

const (
	FirstQuadrant  Quadrant = 0
	SecondQuadrant Quadrant = 1
	ThirdQuadrant  Quadrant = 2
	ForthQuadrant  Quadrant = 3
)

type Direction bool

const (
	AscDirection  Direction = true
	DescDirection Direction = false
)
