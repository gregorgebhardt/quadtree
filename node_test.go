package quadtree

import (
	"reflect"
	"testing"
)

type TestPoint struct {
	x, y float64
	s    string
}

func (p *TestPoint) X() float64 {
	return p.x
}

func (p *TestPoint) Y() float64 {
	return p.y
}

func (p *TestPoint) Equals(other PointPtr) bool {
	return p.x == other.X() && p.y == other.Y()
}

func TestNewNode(t *testing.T) {
	type args struct {
		a   *Area
		cap int
	}
	tests := []struct {
		name string
		args args
		want *Node
	}{
		{"Test1", args{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1},
			&Node{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
				make([]PointPtr, 0, 1), 0, nil}},
		{"Test2", args{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 4},
			&Node{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
				make([]PointPtr, 0, 4), 0, nil}},
		{"Test3", args{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 0},
			nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNode(tt.args.a, tt.args.cap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPoint(t *testing.T) {
	type args struct {
		x float64
		y float64
	}
	tests := []struct {
		name string
		args args
		want *Point
	}{
		{"Test1", args{1., 1.}, NewPoint(1., 1.)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPoint(tt.args.x, tt.args.y); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_collectAllIn(t *testing.T) {
	tests := []struct {
		name          string
		node          *Node
		insertPoints  []PointPtr
		area          *Area
		wantCollected []PointPtr
	}{
		{"Test1_NodeArea", NewNode(NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]PointPtr{NewPoint(1.5, 1.5)},
			NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
			[]PointPtr{NewPoint(1.5, 1.5)}},
		{"Test1_SubArea", NewNode(NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]PointPtr{NewPoint(1.5, 1.5)},
			NewArea(NewPoint(1., 1.), NewPoint(2., 2.)),
			[]PointPtr{NewPoint(1.5, 1.5)}},
		{"Test1_Overlap", NewNode(NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]PointPtr{NewPoint(1.5, 1.5)},
			NewArea(NewPoint(.0, .0), NewPoint(2., 2.)),
			[]PointPtr{NewPoint(1.5, 1.5)}},
		{"Test1_SuperArea", NewNode(NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]PointPtr{NewPoint(1.5, 1.5)},
			NewArea(NewPoint(.0, .0), NewPoint(4., 4.)),
			[]PointPtr{NewPoint(1.5, 1.5)}},
		{"Test1_EdgeIn1", NewNode(NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]PointPtr{NewPoint(1.5, 1.5)},
			NewArea(NewPoint(1.5, 1.0), NewPoint(2., 2.)),
			[]PointPtr{NewPoint(1.5, 1.5)}},
		{"Test1_EdgeIn2", NewNode(NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]PointPtr{NewPoint(1.5, 1.5)},
			NewArea(NewPoint(1.0, 1.5), NewPoint(2., 2.)),
			[]PointPtr{NewPoint(1.5, 1.5)}},
		{"Test1_CornerIn", NewNode(NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]PointPtr{NewPoint(1.5, 1.5)},
			NewArea(NewPoint(1.5, 1.5), NewPoint(2., 2.)),
			[]PointPtr{NewPoint(1.5, 1.5)}},
		{"Test1_EdgeOut2", NewNode(NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]PointPtr{NewPoint(1.5, 1.5)},
			NewArea(NewPoint(1., 1.), NewPoint(1.5, 2.)),
			[]PointPtr{}},
		{"Test1_EdgeOut2", NewNode(NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]PointPtr{NewPoint(1.5, 1.5)},
			NewArea(NewPoint(1., 1.), NewPoint(2., 1.5)),
			[]PointPtr{}},
		{"Test1_CornerOut", NewNode(NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]PointPtr{NewPoint(1.5, 1.5)},
			NewArea(NewPoint(1., 1.), NewPoint(1.5, 1.5)),
			[]PointPtr{}},
		{"Test1_AreaOut", NewNode(NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]PointPtr{NewPoint(1.5, 1.5)},
			NewArea(NewPoint(2., 2.), NewPoint(3., 4.5)),
			[]PointPtr{}},

		{"Test2_NodeArea", NewNode(NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]PointPtr{
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
			[]PointPtr{
				NewPoint(1.5, 1.5), // SW I
				NewPoint(1.2, 1.7), // SW
				NewPoint(1.1, 1.0), // SW
				NewPoint(1.8, 1.7), // SW I
				NewPoint(2.2, 1.5), // SE I
				NewPoint(2.6, 1.1), // SE
				NewPoint(1.2, 2.3), // NW
				NewPoint(1.6, 2.1), // NW I
			}},
		{"Test2_SW", NewNode(NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]PointPtr{
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
			[]PointPtr{
				NewPoint(1.5, 1.5), // SW I
				NewPoint(1.2, 1.7), // SW
				NewPoint(1.1, 1.0), // SW
				NewPoint(1.8, 1.7), // SW I
			}},
		{"Test2_SE", NewNode(NewArea(NewPoint(2., 1.), NewPoint(3., 2.)), 1),
			[]PointPtr{
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
			[]PointPtr{
				NewPoint(2.2, 1.5), // SE I
				NewPoint(2.6, 1.1), // SE
			}},
		{"Test2_NW", NewNode(NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]PointPtr{
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
			[]PointPtr{
				NewPoint(1.2, 2.3), // NW
				NewPoint(1.6, 2.1), // NW I
			}},
		{"Test2_NE", NewNode(NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]PointPtr{
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
			[]PointPtr{}},
		{"Test2_Inner", NewNode(NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]PointPtr{
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
			[]PointPtr{
				NewPoint(1.5, 1.5), // SW I
				NewPoint(1.8, 1.7), // SW I
				NewPoint(2.2, 1.5), // SE I
				NewPoint(1.6, 2.1), // NW I
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, p := range tt.insertPoints {
				tt.node.points = append(tt.node.points, p)
			}
			tt.node.num = len(tt.node.points)
			if gotCollected := tt.node.GetArea(tt.area); !reflect.DeepEqual(gotCollected, tt.wantCollected) {
				t.Errorf("GetArea() = %v, want %v", gotCollected, tt.wantCollected)
			}
		})
	}
}

func TestNode_contains(t *testing.T) {
	tests := []struct {
		name  string
		area  *Area
		point *Point
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
			n := NewNode(tt.area, 1)
			if got := n.contains(tt.point); got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_insert(t *testing.T) {
	tests := []struct {
		name        string
		node        *Node
		points      []*Point
		insertPoint *Point
		wantErr     bool
		wantNode    *Node
	}{
		{"TestEmpty", NewNode(NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]*Point{}, NewPoint(2.2, 1.4), false,
			&Node{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
				[]PointPtr{NewPoint(2.2, 1.4)}, 1, nil}},

		{"TestNotEmpty", NewNode(NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 2),
			[]*Point{NewPoint(1.1, 1.4)}, NewPoint(2.2, 1.4), false,
			&Node{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
				[]PointPtr{NewPoint(1.1, 1.4), NewPoint(2.2, 1.4)}, 2, nil}},

		{"TestEqual", NewNode(NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 2),
			[]*Point{NewPoint(2.2, 1.4)}, NewPoint(2.2, 1.4), true,
			&Node{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
				[]PointPtr{NewPoint(2.2, 1.4)}, 1, nil}},

		{"TestFull", NewNode(NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]*Point{NewPoint(2.2, 1.4)}, NewPoint(1.4, 1.1), false,
			&Node{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
				nil, 2, []*Node{
					NewNode(NewArea(NewPoint(1., 2.), NewPoint(2., 3.)), 1),
					NewNode(NewArea(NewPoint(2., 2.), NewPoint(3., 3.)), 1),
					{NewArea(NewPoint(1., 1.), NewPoint(2., 2.)), []PointPtr{NewPoint(1.4, 1.1)}, 1, nil},
					{NewArea(NewPoint(2., 1.), NewPoint(3., 2.)), []PointPtr{NewPoint(2.2, 1.4)}, 1, nil},
				}}},

		{"TestFullEqual", NewNode(NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]*Point{NewPoint(2.2, 1.4)}, NewPoint(2.2, 1.4), true,
			&Node{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
				[]PointPtr{NewPoint(2.2, 1.4)}, 1, nil}},
	}
	for _, tt := range tests {
		for _, p := range tt.points {
			tt.node.points = append(tt.node.points, p)
		}
		tt.node.num = len(tt.node.points)
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
		area     *Area
		points   []PointPtr
		num      int
		children []*Node
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"Test1", fields{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
			[]PointPtr{NewPoint(2., 2.)}, 1, nil}, true},
		{"Test1", fields{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), nil, 4,
			[]*Node{
				NewNode(NewArea(NewPoint(1., 2.), NewPoint(2., 3.)), 4),
				NewNode(NewArea(NewPoint(2., 2.), NewPoint(3., 3.)), 4),
				NewNode(NewArea(NewPoint(1., 1.), NewPoint(2., 2.)), 4),
				NewNode(NewArea(NewPoint(2., 1.), NewPoint(3., 2.)), 4),
			}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Node{
				area:     tt.fields.area,
				points:   tt.fields.points,
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
		node         *Node
		insertPoints []PointPtr
		search       *Point
		want         int
	}{
		{"Test1_ExistingPoint", NewNode(NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]PointPtr{
				&TestPoint{1.5, 1.5, "p"}, // SW I
				&TestPoint{1.2, 1.7, "p"}, // SW
				&TestPoint{1.1, 1.0, "p"}, // SW
				&TestPoint{1.8, 1.7, "p"}, // SW I
				&TestPoint{2.2, 1.5, "p"}, // SE I
				&TestPoint{2.6, 1.1, "p"}, // SE
				&TestPoint{1.2, 2.3, "p"}, // NW
				&TestPoint{1.6, 2.1, "p"}, // NW I
			},
			NewPoint(1.5, 1.5),
			0},
		{"Test1_MissingPoint", NewNode(NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 1),
			[]PointPtr{
				&TestPoint{1.5, 1.5, "p"}, // SW I
				&TestPoint{1.2, 1.7, "p"}, // SW
				&TestPoint{1.1, 1.0, "p"}, // SW
				&TestPoint{1.8, 1.7, "p"}, // SW I
				&TestPoint{2.2, 1.5, "p"}, // SE I
				&TestPoint{2.6, 1.1, "p"}, // SE
				&TestPoint{1.2, 2.3, "p"}, // NW
				&TestPoint{1.6, 2.1, "p"}, // NW I
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
		area     *Area
		points   []PointPtr
		num      int
		children []*Node
	}
	tests := []struct {
		name   string
		fields fields
		want   *Node
	}{
		{"Test1", fields{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
			[]PointPtr{
				NewPoint(1.5, 1.5),
				NewPoint(1.5, 2.5),
				NewPoint(2.5, 1.5),
				NewPoint(2.5, 2.5),
			}, 4, nil},
			&Node{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), nil, 4,
				[]*Node{
					{NewArea(NewPoint(1., 2.), NewPoint(2., 3.)), []PointPtr{NewPoint(1.5, 2.5)}, 1, nil},
					{NewArea(NewPoint(2., 2.), NewPoint(3., 3.)), []PointPtr{NewPoint(2.5, 2.5)}, 1, nil},
					{NewArea(NewPoint(1., 1.), NewPoint(2., 2.)), []PointPtr{NewPoint(1.5, 1.5)}, 1, nil},
					{NewArea(NewPoint(2., 1.), NewPoint(3., 2.)), []PointPtr{NewPoint(2.5, 1.5)}, 1, nil},
				}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Node{
				area:     tt.fields.area,
				points:   tt.fields.points,
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
		area  *Area
		point *Point
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
			n := NewNode(tt.area, 1)
			if got := n.whichQuadrant(tt.point); got != tt.want {
				t.Errorf("whichQuadrant() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPoint_equals(t *testing.T) {
	tests := []struct {
		name   string
		pointA *Point
		pointB *Point
		want   bool
	}{
		{"TestEq1", NewPoint(1., 1.), NewPoint(1., 1.), true},
		{"TestEq2", NewPoint(1., 1.), NewPoint(1., 1.), true},
		{"TestNeq1", NewPoint(1., 2.), NewPoint(1., 1.), false},
		{"TestNeq2", NewPoint(3., 1.), NewPoint(1., 1.), false},
		{"TestNeq3", NewPoint(3., 1.), NewPoint(1., 3.), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pointA.Equals(tt.pointB); got != tt.want {
				t.Errorf("Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}
