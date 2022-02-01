package quadtree

import "fmt"

type Quadrant int

const (
	NW Quadrant = iota
	NE
	SW
	SE
)

type Node struct {
	area   *Area
	points []PointPtr
	num    int

	children []*Node
}

func NewNode(a *Area, cap int) *Node {
	if cap <= 0 {
		return nil
	}
	return &Node{area: a, points: make([]PointPtr, 0, cap), children: nil}
}

func (n *Node) isLeaf() bool {
	return n.children == nil
}

func (n *Node) contains(p PointPtr) bool {
	if n == nil {
		return false
	}
	return n.area.containsPoint(p)
}

func (n *Node) Get(point PointPtr) PointPtr {
	if n.points != nil {
		for _, p := range n.points {
			if point.Equals(p) {
				return p
			}
		}
		return nil
	} else {
		q := n.whichQuadrant(point)
		return n.children[q].Get(point)
	}
}

func (n *Node) whichQuadrant(p PointPtr) Quadrant {
	if p.Y() >= n.area.c.y {
		//	northern quadrants
		if p.X() >= n.area.c.x {
			return NE
		}
		return NW
	} else {
		//	southern quadrants
		if p.X() >= n.area.c.x {
			return SE
		}
		return SW
	}
}

func (n *Node) split() {
	n.children = make([]*Node, 4)
	for i, a := range n.area.split() {
		n.children[i] = NewNode(a, cap(n.points))
	}
	var q Quadrant
	for _, p := range n.points {
		q = n.whichQuadrant(p)
		_ = n.children[q].Insert(p)
	}
	n.points = nil
}

type PointError struct {
	msg string
	p   PointPtr
}

func (e *PointError) Error() string {
	return fmt.Sprintf("%s:\n%v", e.msg, e.p)
}

func PointExistsError(p PointPtr) *PointError {
	return &PointError{"Point does already exist in Quadtree.", p}
}

func (n *Node) Insert(p PointPtr) error {
	if n.isLeaf() && len(n.points) < cap(n.points) {
		for _, b := range n.points {
			if b.Equals(p) {
				return PointExistsError(b)
			}
		}
		n.points = append(n.points, p)
		n.num++
		return nil
	} else {
		if n.isLeaf() {
			for _, b := range n.points {
				if b.Equals(p) {
					return PointExistsError(b)
				}
			}
			n.split()
		}
		q := n.whichQuadrant(p)
		err := n.children[q].Insert(p)
		if err == nil {
			n.num++
		}
		return err
	}
}

func (n *Node) GetArea(a *Area) (collected []PointPtr) {
	if n.isLeaf() {
		collected = make([]PointPtr, 0, len(n.points))
		for _, p := range n.points {
			if a.containsPoint(p) {
				collected = append(collected, p)
			}
		}
		return collected
	} else {
		c := make(chan []PointPtr)
		defer close(c)
		for _, child := range n.children {
			child := child
			go func() {
				if child.area.intersects(a) {
					c <- child.GetArea(a)
				} else {
					c <- make([]PointPtr, 0)
				}
			}()
		}
		collected = make([]PointPtr, 0, n.num)
		for i := 0; i < 4; i++ {
			collected = append(collected, <-c...)
		}
	}
	return collected
}
