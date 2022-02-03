package quadtree

type Area[T Number] struct {
	l, u, c Point[T]
}

func NewArea[T Number](l, u *Point[T]) *Area[T] {
	a := &Area[T]{*l, *u, Point[T]{(l.x + u.x) / 2., (l.y + u.y) / 2.}}
	// swap elements if necessary
	if a.l.x > a.u.x {
		a.l.x, a.u.x = a.u.x, a.l.x
	}
	if a.l.y > a.u.y {
		a.l.y, a.u.y = a.u.y, a.l.y
	}
	return a
}

func NewAreaAround[T Number](c *Point[T], dX, dY T) *Area[T] {
	return &Area[T]{Point[T]{c.x - dX, c.y - dY}, Point[T]{c.x + dX, c.y + dY}, *c}
}

func (a *Area[T]) containsPoint(p *Point[T]) bool {
	return a.contains(p.x, p.y)
}

func (a *Area[T]) contains(x, y T) bool {
	// lower bound is inclusive, upper bound is exclusive
	return !(x < a.l.x || x >= a.u.x || y < a.l.y || y >= a.u.y)
}

func (a *Area[T]) intersects(other *Area[T]) bool {
	return !(a.l.x > other.u.x || a.l.y > other.u.y || a.u.x < other.l.x || a.u.y < other.l.y)
}

func (a *Area[T]) split() [4]*Area[T] {
	c := NewPoint[T]((a.l.x+a.u.x)/2., (a.l.y+a.u.y)/2.)
	return [4]*Area[T]{
		NewArea[T](NewPoint[T](a.l.x, a.u.y), c),
		NewArea[T](c, &a.u),
		NewArea[T](&a.l, c),
		NewArea[T](c, NewPoint[T](a.u.x, a.l.y)),
	}
}
