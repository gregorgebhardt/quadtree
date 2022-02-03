package quadtree

import "fmt"

//type Element[T Number, V any] struct {
//	p     Point[T]
//	value V
//}

type Element[V any, T Number] interface {
	*V
	p() *Point[T]
	equals(other *V) bool
}

//func (e *Element[T, V]) equals(other *Element[T, V]) bool {
//	return e.p.equals(other.p)
//}

type Quadrant int

const (
	NW Quadrant = iota
	NE
	SW
	SE
)

type Node[T Number, V any, E Element[V, T]] struct {
	area     *Area[T]
	elements []E
	num      int

	children []*Node[T, V, E]
}

func NewNode[T Number, V any, E Element[V, T]](a *Area[T], cap int) *Node[T, V, E] {
	if cap <= 0 {
		return nil
	}
	return &Node[T, V, E]{area: a, elements: make([]E, 0, cap), children: nil}
}

func NewTree[T Number, V any, E Element[V, T]](xMin, xMax, yMin, yMax T, cap int) *Node[T, V, E] {
	a := NewArea(NewPoint(xMin, yMin), NewPoint(xMax, yMax))
	return &Node[T, V, E]{area: a, elements: make([]E, 0, cap), children: nil}
}

func (n *Node[T, V, E]) isLeaf() bool {
	return n.children == nil
}

func (n *Node[T, V, E]) contains(p *Point[T]) bool {
	if n == nil {
		return false
	}
	return n.area.containsPoint(p)
}

func (n *Node[T, V, E]) Get(p *Point[T]) E {
	if n.elements != nil {
		for _, e := range n.elements {
			if p.equals(e.p()) {
				return e
			}
		}
		return nil
	} else {
		q := n.whichQuadrant(p)
		return n.children[q].Get(p)
	}
}

func (n *Node[T, V, E]) GetArea(a *Area[T]) []E {
	return n.GetAreaFiltered(a, func(_ E) bool { return true })
}

func (n *Node[T, V, E]) GetAreaFiltered(a *Area[T], f func(E) bool) (collected []E) {
	if n.isLeaf() {
		collected = make([]E, 0, len(n.elements))
		for _, e := range n.elements {
			if a.containsPoint(e.p()) && f(e) {
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

func (n *Node[T, V, E]) whichQuadrant(p *Point[T]) Quadrant {
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

func (n *Node[T, V, E]) split() {
	n.children = make([]*Node[T, V, E], 4)
	for i, a := range n.area.split() {
		n.children[i] = NewNode[T, V, E](a, cap(n.elements))
	}
	var q Quadrant
	for _, e := range n.elements {
		q = n.whichQuadrant(e.p())
		_ = n.children[q].Insert(e)
	}
	n.elements = nil
}

type ElementError[T Number, V any, E Element[V, T]] struct {
	msg string
	e   E
}

func (e *ElementError[T, V, E]) Error() string {
	return fmt.Sprintf("%s:\n%v", e.msg, e.e)
}

func PointExistsError[T Number, V any, E Element[V, T]](e E) *ElementError[T, V, E] {
	return &ElementError[T, V, E]{"Point does already exist in Quadtree.", e}
}

func (n *Node[T, V, E]) Insert(e E) error {
	if n.isLeaf() && len(n.elements) < cap(n.elements) {
		for _, b := range n.elements {
			if b.equals(e) {
				return PointExistsError[T, V, E](b)
			}
		}
		n.elements = append(n.elements, e)
		n.num++
		return nil
	} else {
		if n.isLeaf() {
			for _, b := range n.elements {
				if b.equals(e) {
					return PointExistsError[T, V, E](b)
				}
			}
			n.split()
		}
		q := n.whichQuadrant(e.p())
		err := n.children[q].Insert(e)
		if err == nil {
			n.num++
		}
		return err
	}
}
