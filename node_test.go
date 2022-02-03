package quadtree

import (
	"reflect"
	"testing"
)

type TestElement struct {
	x, y float64
	s    string
}

type TestNode Node[TestElement, float64]

func (e *TestElement) p() *Point[float64] {
	return NewPoint[float64](e.x, e.y)
}

func (e *TestElement) equals(other *TestElement) bool {
	return e.x == other.x && e.y == other.y
}

func TestNewNode(t *testing.T) {
	type args struct {
		a   *Area[float64]
		cap int
	}
	tests := []struct {
		name string
		args args
		want *Node[TestElement, float64]
	}{
		{"Test1", args{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1},
			&Node[TestElement, float64]{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
				make([]*TestElement, 0, 1), 0, nil}},
		{"Test2", args{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 4},
			&Node[TestElement, float64]{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
				make([]*TestElement, 0, 4), 0, nil}},
		{"Test3", args{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 0},
			nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNode[TestElement](tt.args.a, tt.args.cap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_collectAllIn(t *testing.T) {
	//NewNode := NewNode

	tests := []struct {
		name          string
		node          *Node[Point[float64], float64, *Point[float64]]
		insertPoints  []*Point[float64]
		area          *Area[float64]
		wantCollected []*Point[float64]
	}{
		{"Test1_NodeArea", NewNode[Point[float64]](NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]*Point[float64]{NewPoint(1.5, 1.5)},
			NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
			[]*Point[float64]{NewPoint(1.5, 1.5)}},
		{"Test1_SubArea", NewNode[Point[float64]](NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]*Point[float64]{NewPoint(1.5, 1.5)},
			NewArea(NewPoint(1., 1.), NewPoint(2., 2.)),
			[]*Point[float64]{NewPoint(1.5, 1.5)}},
		{"Test1_Overlap", NewNode[Point[float64]](NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]*Point[float64]{NewPoint(1.5, 1.5)},
			NewArea(NewPoint(.0, .0), NewPoint(2., 2.)),
			[]*Point[float64]{NewPoint(1.5, 1.5)}},
		{"Test1_SuperArea", NewNode[Point[float64]](NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]*Point[float64]{NewPoint(1.5, 1.5)},
			NewArea(NewPoint(.0, .0), NewPoint(4., 4.)),
			[]*Point[float64]{NewPoint(1.5, 1.5)}},
		{"Test1_EdgeIn1", NewNode[Point[float64]](NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]*Point[float64]{NewPoint(1.5, 1.5)},
			NewArea(NewPoint(1.5, 1.0), NewPoint(2., 2.)),
			[]*Point[float64]{NewPoint(1.5, 1.5)}},
		{"Test1_EdgeIn2", NewNode[Point[float64]](NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]*Point[float64]{NewPoint(1.5, 1.5)},
			NewArea(NewPoint(1.0, 1.5), NewPoint(2., 2.)),
			[]*Point[float64]{NewPoint(1.5, 1.5)}},
		{"Test1_CornerIn", NewNode[Point[float64]](NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]*Point[float64]{NewPoint(1.5, 1.5)},
			NewArea(NewPoint(1.5, 1.5), NewPoint(2., 2.)),
			[]*Point[float64]{NewPoint(1.5, 1.5)}},
		{"Test1_EdgeOut2", NewNode[Point[float64]](NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]*Point[float64]{NewPoint(1.5, 1.5)},
			NewArea(NewPoint(1., 1.), NewPoint(1.5, 2.)),
			[]*Point[float64]{}},
		{"Test1_EdgeOut2", NewNode[Point[float64]](NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]*Point[float64]{NewPoint(1.5, 1.5)},
			NewArea(NewPoint(1., 1.), NewPoint(2., 1.5)),
			[]*Point[float64]{}},
		{"Test1_CornerOut", NewNode[Point[float64]](NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]*Point[float64]{NewPoint(1.5, 1.5)},
			NewArea(NewPoint(1., 1.), NewPoint(1.5, 1.5)),
			[]*Point[float64]{}},
		{"Test1_AreaOut", NewNode[Point[float64]](NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]*Point[float64]{NewPoint(1.5, 1.5)},
			NewArea(NewPoint(2., 2.), NewPoint(3., 4.5)),
			[]*Point[float64]{}},

		{"Test2_NodeArea", NewNode[Point[float64]](NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]*Point[float64]{
				NewPoint(1.5, 1.5), // SW I
				NewPoint(1.2, 1.7), // SW
				NewPoint(1.1, 1.0), // SW
				NewPoint(1.8, 1.7), // SW I
				NewPoint(2.2, 1.5), // SE I
				NewPoint(2.6, 1.1), // SE
				NewPoint(1.2, 2.3), // NW
				NewPoint(1.6, 2.1), // NW I
			},
			NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
			[]*Point[float64]{
				NewPoint(1.5, 1.5), // SW I
				NewPoint(1.2, 1.7), // SW
				NewPoint(1.1, 1.0), // SW
				NewPoint(1.8, 1.7), // SW I
				NewPoint(2.2, 1.5), // SE I
				NewPoint(2.6, 1.1), // SE
				NewPoint(1.2, 2.3), // NW
				NewPoint(1.6, 2.1), // NW I
			}},
		{"Test2_SW", NewNode[Point[float64]](NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]*Point[float64]{
				NewPoint(1.5, 1.5), // SW I
				NewPoint(1.2, 1.7), // SW
				NewPoint(1.1, 1.0), // SW
				NewPoint(1.8, 1.7), // SW I
				NewPoint(2.2, 1.5), // SE I
				NewPoint(2.6, 1.1), // SE
				NewPoint(1.2, 2.3), // NW
				NewPoint(1.6, 2.1), // NW I
			},
			NewArea(NewPoint(1., 1.), NewPoint(2., 2.)),
			[]*Point[float64]{
				NewPoint(1.5, 1.5), // SW I
				NewPoint(1.2, 1.7), // SW
				NewPoint(1.1, 1.0), // SW
				NewPoint(1.8, 1.7), // SW I
			}},
		{"Test2_SE", NewNode[Point[float64]](NewArea(NewPoint(2., 1.), NewPoint(3., 2.)), 1),
			[]*Point[float64]{
				NewPoint(1.5, 1.5), // SW I
				NewPoint(1.2, 1.7), // SW
				NewPoint(1.1, 1.0), // SW
				NewPoint(1.8, 1.7), // SW I
				NewPoint(2.2, 1.5), // SE I
				NewPoint(2.6, 1.1), // SE
				NewPoint(1.2, 2.3), // NW
				NewPoint(1.6, 2.1), // NW I
			},
			NewArea(NewPoint(2., 1.), NewPoint(3., 2.)),
			[]*Point[float64]{
				NewPoint(2.2, 1.5), // SE I
				NewPoint(2.6, 1.1), // SE
			}},
		{"Test2_NW", NewNode[Point[float64]](NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]*Point[float64]{
				NewPoint(1.5, 1.5), // SW I
				NewPoint(1.2, 1.7), // SW
				NewPoint(1.1, 1.0), // SW
				NewPoint(1.8, 1.7), // SW I
				NewPoint(2.2, 1.5), // SE I
				NewPoint(2.6, 1.1), // SE
				NewPoint(1.2, 2.3), // NW
				NewPoint(1.6, 2.1), // NW I
			},
			NewArea(NewPoint(1., 2.), NewPoint(2., 3.)),
			[]*Point[float64]{
				NewPoint(1.2, 2.3), // NW
				NewPoint(1.6, 2.1), // NW I
			}},
		{"Test2_NE", NewNode[Point[float64]](NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]*Point[float64]{
				NewPoint(1.5, 1.5), // SW I
				NewPoint(1.2, 1.7), // SW
				NewPoint(1.1, 1.0), // SW
				NewPoint(1.8, 1.7), // SW I
				NewPoint(2.2, 1.5), // SE I
				NewPoint(2.6, 1.1), // SE
				NewPoint(1.2, 2.3), // NW
				NewPoint(1.6, 2.1), // NW I
			},
			NewArea(NewPoint(2., 2.), NewPoint(3., 3.)),
			[]*Point[float64]{}},
		{"Test2_Inner", NewNode[Point[float64]](NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]*Point[float64]{
				NewPoint(1.5, 1.5), // SW I
				NewPoint(1.2, 1.7), // SW
				NewPoint(1.1, 1.0), // SW
				NewPoint(1.8, 1.7), // SW I
				NewPoint(2.2, 1.5), // SE I
				NewPoint(2.6, 1.1), // SE
				NewPoint(1.2, 2.3), // NW
				NewPoint(1.6, 2.1), // NW I
			},
			NewArea(NewPoint(1.5, 1.5), NewPoint(2.5, 2.5)),
			[]*Point[float64]{
				NewPoint(1.5, 1.5), // SW I
				NewPoint(1.8, 1.7), // SW I
				NewPoint(2.2, 1.5), // SE I
				NewPoint(1.6, 2.1), // NW I
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, p := range tt.insertPoints {
				tt.node.elements = append(tt.node.elements, p)
			}
			tt.node.num = len(tt.node.elements)
			if gotCollected := tt.node.GetArea(tt.area); !reflect.DeepEqual(gotCollected, tt.wantCollected) {
				t.Errorf("GetArea() = %v, want %v", gotCollected, tt.wantCollected)
			}
		})
	}
}

func TestNode_contains(t *testing.T) {
	tests := []struct {
		name  string
		area  *Area[float64]
		point *Point[float64]
		want  bool
	}{
		{"TestIn1", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), NewPoint(1.5, 1.5), true},
		{"TestIn2", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), NewPoint(2.5, 2.5), true},
		{"TestIn3", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), NewPoint(1.5, 2.5), true},
		{"TestIn4", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), NewPoint(2.5, 1.5), true},

		{"TestOut1", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), NewPoint(-1.5, -1.5), false},
		{"TestOut2", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), NewPoint(4.5, 2.5), false},
		{"TestOut3", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), NewPoint(1.5, 4.5), false},
		{"TestOut4", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), NewPoint(2.5, -1.5), false},

		{"TestNeg1", NewArea(NewPoint(-1., -1.), NewPoint(-3., -3.)), NewPoint(-1.5, -1.5), true},
		{"TestNeg2", NewArea(NewPoint(-1., -1.), NewPoint(-3., -3.)), NewPoint(-2.5, -2.5), true},
		{"TestNeg3", NewArea(NewPoint(-1., -1.), NewPoint(-3., -3.)), NewPoint(-3.5, -2.5), false},
		{"TestNeg4", NewArea(NewPoint(-1., -1.), NewPoint(-3., -3.)), NewPoint(2.5, 1.5), false},

		{"TestBorder1", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), NewPoint(1., 1.), true},
		{"TestBorder2", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), NewPoint(1., 2.), true},
		{"TestBorder3", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), NewPoint(3., 2.), false},
		{"TestBorder4", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), NewPoint(1., 3.), false},
		{"TestBorder5", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), NewPoint(3., 3.), false},
		{"TestBorder6", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), NewPoint(3., 1.), false},
		{"TestBorder7", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), NewPoint(2., 3.), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewNode[Point[float64]](tt.area, 1)
			if got := n.contains(tt.point); got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_insert(t *testing.T) {
	tests := []struct {
		name        string
		node        *Node[Point[float64], float64, *Point[float64]]
		points      []*Point[float64]
		insertPoint *Point[float64]
		wantErr     bool
		wantNode    *Node[Point[float64], float64]
	}{
		{"TestEmpty", NewNode[Point[float64]](NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]*Point[float64]{}, NewPoint(2.2, 1.4), false,
			&Node[Point[float64], float64]{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
				[]*Point[float64]{NewPoint(2.2, 1.4)}, 1, nil}},

		{"TestNotEmpty", NewNode[Point[float64]](NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 2),
			[]*Point[float64]{NewPoint(1.1, 1.4)}, NewPoint(2.2, 1.4), false,
			&Node[Point[float64], float64]{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
				[]*Point[float64]{NewPoint(1.1, 1.4), NewPoint(2.2, 1.4)}, 2, nil}},

		{"TestEqual", NewNode[Point[float64]](NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 2),
			[]*Point[float64]{NewPoint(2.2, 1.4)}, NewPoint(2.2, 1.4), true,
			&Node[Point[float64], float64]{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
				[]*Point[float64]{NewPoint(2.2, 1.4)}, 1, nil}},

		{"TestFull", NewNode[Point[float64]](NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]*Point[float64]{NewPoint(2.2, 1.4)}, NewPoint(1.4, 1.1), false,
			&Node[Point[float64], float64]{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
				nil, 2, []*Node[Point[float64], float64]{
					NewNode[Point[float64]](NewArea(NewPoint(1., 2.), NewPoint(2., 3.)), 1),
					NewNode[Point[float64]](NewArea(NewPoint(2., 2.), NewPoint(3., 3.)), 1),
					{NewArea(NewPoint(1., 1.), NewPoint(2., 2.)), []*Point[float64]{NewPoint(1.4, 1.1)}, 1, nil},
					{NewArea(NewPoint(2., 1.), NewPoint(3., 2.)), []*Point[float64]{NewPoint(2.2, 1.4)}, 1, nil},
				}}},

		{"TestFullEqual", NewNode[Point[float64]](NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]*Point[float64]{NewPoint(2.2, 1.4)}, NewPoint(2.2, 1.4), true,
			&Node[Point[float64], float64]{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
				[]*Point[float64]{NewPoint(2.2, 1.4)}, 1, nil}},
	}
	for _, tt := range tests {
		for _, p := range tt.points {
			tt.node.elements = append(tt.node.elements, p)
		}
		tt.node.num = len(tt.node.elements)
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.node.Insert(tt.insertPoint); (err != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.node, tt.wantNode) {
				t.Errorf("node = %v, wantNode %v", tt.node, tt.wantNode)
			}
		})
	}
}

func TestNode_isLeaf(t *testing.T) {
	type fields struct {
		area     *Area[float64]
		points   []*Point[float64]
		num      int
		children []*Node[Point[float64], float64]
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"Test1", fields{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
			[]*Point[float64]{NewPoint(2., 2.)}, 1, nil}, true},
		{"Test1", fields{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), nil, 4,
			[]*Node[Point[float64], float64]{
				NewNode[Point[float64]](NewArea(NewPoint(1., 2.), NewPoint(2., 3.)), 4),
				NewNode[Point[float64]](NewArea(NewPoint(2., 2.), NewPoint(3., 3.)), 4),
				NewNode[Point[float64]](NewArea(NewPoint(1., 1.), NewPoint(2., 2.)), 4),
				NewNode[Point[float64]](NewArea(NewPoint(2., 1.), NewPoint(3., 2.)), 4),
			}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Node[Point[float64], float64]{
				area:     tt.fields.area,
				elements: tt.fields.points,
				num:      tt.fields.num,
				children: tt.fields.children,
			}
			if got := n.isLeaf(); got != tt.want {
				t.Errorf("isLeaf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_search(t *testing.T) {
	tests := []struct {
		name         string
		node         *Node[TestElement, float64]
		insertPoints []*TestElement
		search       *Point[float64]
		want         int
	}{
		{"Test1_ExistingPoint", NewNode[TestElement](NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]*TestElement{
				&TestElement{1.5, 1.5, "p"}, // SW I
				&TestElement{1.2, 1.7, "p"}, // SW
				&TestElement{1.1, 1.0, "p"}, // SW
				&TestElement{1.8, 1.7, "p"}, // SW I
				&TestElement{2.2, 1.5, "p"}, // SE I
				&TestElement{2.6, 1.1, "p"}, // SE
				&TestElement{1.2, 2.3, "p"}, // NW
				&TestElement{1.6, 2.1, "p"}, // NW I
			},
			NewPoint(1.5, 1.5),
			0},
		{"Test1_MissingPoint", NewNode[TestElement](NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]*TestElement{
				&TestElement{1.5, 1.5, "p"}, // SW I
				&TestElement{1.2, 1.7, "p"}, // SW
				&TestElement{1.1, 1.0, "p"}, // SW
				&TestElement{1.8, 1.7, "p"}, // SW I
				&TestElement{2.2, 1.5, "p"}, // SE I
				&TestElement{2.6, 1.1, "p"}, // SE
				&TestElement{1.2, 2.3, "p"}, // NW
				&TestElement{1.6, 2.1, "p"}, // NW I
			},
			NewPoint(1.4, 1.5),
			-1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, p := range tt.insertPoints {
				_ = tt.node.Insert(p)
			}

			if got := tt.node.Get(tt.search); got != nil && got != tt.insertPoints[tt.want] {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			} else if got == nil && tt.want != -1 {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_split(t *testing.T) {
	type fields struct {
		area     *Area[float64]
		points   []*Point[float64]
		num      int
		children []*Node[Point[float64], float64]
	}
	tests := []struct {
		name   string
		fields fields
		want   *Node[Point[float64], float64]
	}{
		{"Test1", fields{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
			[]*Point[float64]{
				NewPoint(1.5, 1.5),
				NewPoint(1.5, 2.5),
				NewPoint(2.5, 1.5),
				NewPoint(2.5, 2.5),
			}, 4, nil},
			&Node[Point[float64], float64]{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), nil, 4,
				[]*Node[Point[float64], float64]{
					{NewArea(NewPoint(1., 2.), NewPoint(2., 3.)), []*Point[float64]{NewPoint(1.5, 2.5)}, 1, nil},
					{NewArea(NewPoint(2., 2.), NewPoint(3., 3.)), []*Point[float64]{NewPoint(2.5, 2.5)}, 1, nil},
					{NewArea(NewPoint(1., 1.), NewPoint(2., 2.)), []*Point[float64]{NewPoint(1.5, 1.5)}, 1, nil},
					{NewArea(NewPoint(2., 1.), NewPoint(3., 2.)), []*Point[float64]{NewPoint(2.5, 1.5)}, 1, nil},
				}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Node[Point[float64], float64]{
				area:     tt.fields.area,
				elements: tt.fields.points,
				num:      tt.fields.num,
				children: tt.fields.children,
			}
			n.split()
			if !reflect.DeepEqual(n, tt.want) {
				t.Errorf("Get() = %v, want %v", n, tt.want)
			}
		})
	}
}

func TestNode_whichQuadrant(t *testing.T) {

	tests := []struct {
		name  string
		area  *Area[float64]
		point *Point[float64]
		want  Quadrant
	}{
		{"TestIn1", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), NewPoint(1.5, 1.5), SW},
		{"TestIn2", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), NewPoint(2.5, 2.5), NE},
		{"TestIn3", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), NewPoint(1.5, 2.5), NW},
		{"TestIn4", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), NewPoint(2.5, 1.5), SE},

		{"TestOut1", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), NewPoint(-1.5, -1.5), SW},
		{"TestOut2", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), NewPoint(4.5, 2.5), NE},
		{"TestOut3", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), NewPoint(1.5, 4.5), NW},
		{"TestOut4", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), NewPoint(2.5, -1.5), SE},

		{"TestNeg1", NewArea(NewPoint(-1., -1.), NewPoint(-3., -3.)), NewPoint(-1.5, -1.5), NE},
		{"TestNeg2", NewArea(NewPoint(-1., -1.), NewPoint(-3., -3.)), NewPoint(-2.5, -2.5), SW},
		{"TestNeg3", NewArea(NewPoint(-1., -1.), NewPoint(-3., -3.)), NewPoint(-1.5, -2.5), SE},
		{"TestNeg4", NewArea(NewPoint(-1., -1.), NewPoint(-3., -3.)), NewPoint(-2.5, -1.5), NW},

		{"TestBorder1", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), NewPoint(2., 2.), NE},
		{"TestBorder2", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), NewPoint(2., 1.5), SE},
		{"TestBorder3", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), NewPoint(1.5, 2.), NW},
		{"TestBorder4", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), NewPoint(1.999, 1.999), SW},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewNode[Point[float64]](tt.area, 1)
			if got := n.whichQuadrant(tt.point); got != tt.want {
				t.Errorf("whichQuadrant() = %v, want %v", got, tt.want)
			}
		})
	}
}
