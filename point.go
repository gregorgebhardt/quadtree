package quadtree

import "constraints"

type Number interface {
	constraints.Integer | constraints.Float
}

type Point[T Number] struct {
	x, y T
}

func NewPoint[T Number](x, y T) *Point[T] {
	return &Point[T]{x: x, y: y}
}

func (p *Point[T]) equals(other *Point[T]) bool {
	return p.x == other.x && p.y == other.y
}
