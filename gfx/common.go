package gfx

/*
	Illustrated coordinate system of launchpad:

	+--------- arrow keys -----------+  +--- mode keys ---+
	{0, 8} {1, 8} {2, 8} {3, 8} {4, 8} {5, 8} {6, 8} {7, 8} |
	----------------------------------------------------------------
	{0, 0} {1, 0} {2, 0} {3, 0} {4, 0} {5, 0} {6, 0} {7, 0} | {8, 0} vol
	----------------------------------------------------------------
	{0, 1} {1, 1} {2, 1} {3, 1} {4, 1} {5, 1} {6, 1} {7, 1} | {8, 1} pan
	----------------------------------------------------------------
	{0, 2} {1, 2} {2, 2} {3, 2} {4, 2} {5, 2} {6, 2} {7, 2} | {8, 2} sndA
	----------------------------------------------------------------
	{0, 3} {1, 3} {2, 3} {3, 3} {4, 3} {5, 3} {6, 3} {7, 3} | {8, 3} sndB
	----------------------------------------------------------------
	{0, 4} {1, 4} {2, 4} {3, 4} {4, 4} {5, 4} {6, 4} {7, 4} | {8, 4} stop
	----------------------------------------------------------------
	{0, 5} {1, 5} {2, 5} {3, 5} {4, 5} {5, 5} {6, 5} {7, 5} | {8, 5} trk on
	----------------------------------------------------------------
	{0, 6} {1, 6} {2, 6} {3, 6} {4, 6} {5, 6} {6, 6} {7, 6} | {8, 6} solo
	----------------------------------------------------------------
	{0, 7} {1, 7} {2, 7} {3, 7} {4, 7} {5, 7} {6, 7} {7, 7} | {8, 7} arm
	----------------------------------------------------------------

*/

type coord struct {
	X int
	Y int
}

const (
	minY      = 0
	maxY      = 7
	minX      = 0
	maxX      = 7
	padHeight = 8
	padLength = 8
)

type Quadrant int

const (
	FirstQuadrant  Quadrant = 1
	SecondQuadrant Quadrant = 2
	ThirdQuadrant  Quadrant = 3
	ForthQuadrant  Quadrant = 4
)

type Direction bool

const (
	AscDirection  Direction = true
	DescDirection Direction = false
)
