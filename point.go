package quadtree

type PointPtr interface {
	X() float64
	Y() float64
	Equals(ptr PointPtr) bool
}

type Point struct {
	x, y float64
}

func (p *Point) X() float64 {
	return p.x
}

func (p *Point) Y() float64 {
	return p.y
}

func NewPoint(x, y float64) *Point {
	return &Point{x: x, y: y}
}

func (p *Point) Equals(other PointPtr) bool {
	return p.x == other.X() && p.y == other.Y()
}
