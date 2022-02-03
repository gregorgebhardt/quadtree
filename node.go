package quadtree

import "fmt"

type Element[V any, T Number] interface {
	*V
	P() *Point[T]
}

type Quadrant int

const (
	NW Quadrant = iota
	NE
	SW
	SE
)

type Node[V any, T Number, E Element[V, T]] struct {
	area     *Area[T]
	elements []E
	num      int

	children []*Node[V, T, E]
}

func NewNode[V any, T Number, E Element[V, T]](a *Area[T], cap int) *Node[V, T, E] {
	if cap <= 0 {
		return nil
	}
	return &Node[V, T, E]{area: a, elements: make([]E, 0, cap), children: nil}
}

func NewTree[V any, T Number, E Element[V, T]](xMin, xMax, yMin, yMax T, cap int) *Node[V, T, E] {
	a := NewArea(NewPoint(xMin, yMin), NewPoint(xMax, yMax))
	return &Node[V, T, E]{area: a, elements: make([]E, 0, cap), children: nil}
}

func (n *Node[V, T, E]) isLeaf() bool {
	return n.children == nil
}

func (n *Node[V, T, E]) contains(p *Point[T]) bool {
	if n == nil {
		return false
	}
	return n.area.containsPoint(p)
}

func (n *Node[V, T, E]) Get(p *Point[T]) E {
	if n.elements != nil {
		for _, e := range n.elements {
			if p.Equals(e.P()) {
				return e
			}
		}
		return nil
	} else {
		q := n.whichQuadrant(p)
		return n.children[q].Get(p)
	}
}

func (n *Node[V, T, E]) GetArea(a *Area[T]) []E {
	return n.GetAreaFiltered(a, func(_ E) bool { return true })
}

func (n *Node[V, T, E]) GetAreaFiltered(a *Area[T], f func(E) bool) (collected []E) {
	if n.isLeaf() {
		collected = make([]E, 0, len(n.elements))
		for _, e := range n.elements {
			if a.containsPoint(e.P()) && f(e) {
				collected = append(collected, e)
			}
		}
		return collected
	} else {
		c := make(chan []E)
		defer close(c)
		for _, child := range n.children {
			child := child
			go func() {
				if child.area.intersects(a) {
					c <- child.GetArea(a)
				} else {
					c <- make([]E, 0)
				}
			}()
		}
		collected = make([]E, 0, n.num)
		for i := 0; i < 4; i++ {
			collected = append(collected, <-c...)
		}
	}
	return collected
}

func (n *Node[V, T, E]) whichQuadrant(p *Point[T]) Quadrant {
	if p.y >= n.area.c.y {
		//	northern quadrants
		if p.x >= n.area.c.x {
			return NE
		}
		return NW
	} else {
		//	southern quadrants
		if p.x >= n.area.c.x {
			return SE
		}
		return SW
	}
}

func (n *Node[V, T, E]) split() {
	n.children = make([]*Node[V, T, E], 4)
	for i, a := range n.area.split() {
		n.children[i] = NewNode[V, T, E](a, cap(n.elements))
	}
	var q Quadrant
	for _, e := range n.elements {
		q = n.whichQuadrant(e.P())
		_ = n.children[q].Insert(e)
	}
	n.elements = nil
}

type ElementError[V any, T Number, E Element[V, T]] struct {
	msg string
	e   E
}

func (e *ElementError[V, T, E]) Error() string {
	return fmt.Sprintf("%s:\n%v", e.msg, e.e)
}

func PointExistsError[V any, T Number, E Element[V, T]](e E) *ElementError[V, T, E] {
	return &ElementError[V, T, E]{"Point does already exist in Quadtree.", e}
}

func (n *Node[V, T, E]) Insert(e E) error {
	if n.isLeaf() && len(n.elements) < cap(n.elements) {
		for _, b := range n.elements {
			if b.P().Equals(e.P()) {
				return PointExistsError[V, T, E](b)
			}
		}
		n.elements = append(n.elements, e)
		n.num++
		return nil
	} else {
		if n.isLeaf() {
			for _, b := range n.elements {
				if b.P().Equals(e.P()) {
					return PointExistsError[V, T, E](b)
				}
			}
			n.split()
		}
		q := n.whichQuadrant(e.P())
		err := n.children[q].Insert(e)
		if err == nil {
			n.num++
		}
		return err
	}
}
