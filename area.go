package quadtree

type Area struct {
	l, u, c Point
}

func NewArea(l, u *Point) *Area {
	a := &Area{*l, *u, Point{(l.X() + u.X()) / 2., (l.y + u.y) / 2.}}
	// swap points if necessary
	if a.l.x > a.u.x {
		a.l.x, a.u.x = a.u.x, a.l.x
	}
	if a.l.y > a.u.y {
		a.l.y, a.u.y = a.u.y, a.l.y
	}
	return a
}

func NewAreaAround(c Point, dX, dY float64) *Area {
	return &Area{Point{c.x - dX, c.y - dY}, Point{c.x + dX, c.y + dY}, c}
}

func (a *Area) containsPoint(p PointPtr) bool {
	return a.contains(p.X(), p.Y())
}

func (a *Area) contains(x, y float64) bool {
	// lower bound is inclusive, upper bound is exclusive
	return !(x < a.l.x || x >= a.u.x || y < a.l.y || y >= a.u.y)
}

func (a *Area) intersects(other *Area) bool {
	return !(a.l.x > other.u.x || a.l.y > other.u.y || a.u.x < other.l.x || a.u.y < other.l.y)
}

func (a *Area) split() [4]*Area {
	c := NewPoint((a.l.x+a.u.x)/2., (a.l.y+a.u.y)/2.)
	return [4]*Area{
		NewArea(NewPoint(a.l.x, a.u.y), c),
		NewArea(c, &a.u),
		NewArea(&a.l, c),
		NewArea(c, NewPoint(a.u.x, a.l.y)),
	}
}
