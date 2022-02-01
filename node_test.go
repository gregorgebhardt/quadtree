package quadtree

import (
	"reflect"
	"testing"
)

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
				make([]*Point, 0, 1), 0, nil}},
		{"Test2", args{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), 4},
			&Node{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
				make([]*Point, 0, 4), 0, nil}},
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
		want Point
	}{
		{"Test1", args{1., 1.}, Point{1., 1., nil}},
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
		insertPoints  []*Point
		area          *Area
		wantCollected []*Point
	}{
		{"Test1_NodeArea", NewNode(NewArea(Point{1., 1., nil}, Point{3., 3., nil}), 1),
			[]*Point{{1.5, 1.5, nil}},
			NewArea(Point{1., 1., nil}, Point{3., 3., nil}),
			[]*Point{{1.5, 1.5, nil}}},
		{"Test1_SubArea", NewNode(NewArea(Point{1., 1., nil}, Point{3., 3., nil}), 1),
			[]*Point{{1.5, 1.5, nil}},
			NewArea(Point{1., 1., nil}, Point{2., 2., nil}),
			[]*Point{{1.5, 1.5, nil}}},
		{"Test1_Overlap", NewNode(NewArea(Point{1., 1., nil}, Point{3., 3., nil}), 1),
			[]*Point{{1.5, 1.5, nil}},
			NewArea(Point{.0, .0, nil}, Point{2., 2., nil}),
			[]*Point{{1.5, 1.5, nil}}},
		{"Test1_SuperArea", NewNode(NewArea(Point{1., 1., nil}, Point{3., 3., nil}), 1),
			[]*Point{{1.5, 1.5, nil}},
			NewArea(Point{.0, .0, nil}, Point{4., 4., nil}),
			[]*Point{{1.5, 1.5, nil}}},
		{"Test1_EdgeIn1", NewNode(NewArea(Point{1., 1., nil}, Point{3., 3., nil}), 1),
			[]*Point{{1.5, 1.5, nil}},
			NewArea(Point{1.5, 1.0, nil}, Point{2., 2., nil}),
			[]*Point{{1.5, 1.5, nil}}},
		{"Test1_EdgeIn2", NewNode(NewArea(Point{1., 1., nil}, Point{3., 3., nil}), 1),
			[]*Point{{1.5, 1.5, nil}},
			NewArea(Point{1.0, 1.5, nil}, Point{2., 2., nil}),
			[]*Point{{1.5, 1.5, nil}}},
		{"Test1_CornerIn", NewNode(NewArea(Point{1., 1., nil}, Point{3., 3., nil}), 1),
			[]*Point{{1.5, 1.5, nil}},
			NewArea(Point{1.5, 1.5, nil}, Point{2., 2., nil}),
			[]*Point{{1.5, 1.5, nil}}},
		{"Test1_EdgeOut2", NewNode(NewArea(Point{1., 1., nil}, Point{3., 3., nil}), 1),
			[]*Point{{1.5, 1.5, nil}},
			NewArea(Point{1., 1., nil}, Point{1.5, 2., nil}),
			[]*Point{}},
		{"Test1_EdgeOut2", NewNode(NewArea(Point{1., 1., nil}, Point{3., 3., nil}), 1),
			[]*Point{{1.5, 1.5, nil}},
			NewArea(Point{1., 1., nil}, Point{2., 1.5, nil}),
			[]*Point{}},
		{"Test1_CornerOut", NewNode(NewArea(Point{1., 1., nil}, Point{3., 3., nil}), 1),
			[]*Point{{1.5, 1.5, nil}},
			NewArea(Point{1., 1., nil}, Point{1.5, 1.5, nil}),
			[]*Point{}},
		{"Test1_AreaOut", NewNode(NewArea(Point{1., 1., nil}, Point{3., 3., nil}), 1),
			[]*Point{{1.5, 1.5, nil}},
			NewArea(Point{2., 2., nil}, Point{3., 4.5, nil}),
			[]*Point{}},

		{"Test2_NodeArea", NewNode(NewArea(Point{1., 1., nil}, Point{3., 3., nil}), 1),
			[]*Point{
				{1.5, 1.5, nil}, // SW I
				{1.2, 1.7, nil}, // SW
				{1.1, 1.0, nil}, // SW
				{1.8, 1.7, nil}, // SW I
				{2.2, 1.5, nil}, // SE I
				{2.6, 1.1, nil}, // SE
				{1.2, 2.3, nil}, // NW
				{1.6, 2.1, nil}, // NW I
			},
			NewArea(Point{1., 1., nil}, Point{3., 3., nil}),
			[]*Point{
				{1.5, 1.5, nil}, // SW I
				{1.2, 1.7, nil}, // SW
				{1.1, 1.0, nil}, // SW
				{1.8, 1.7, nil}, // SW I
				{2.2, 1.5, nil}, // SE I
				{2.6, 1.1, nil}, // SE
				{1.2, 2.3, nil}, // NW
				{1.6, 2.1, nil}, // NW I
			}},
		{"Test2_SW", NewNode(NewArea(Point{1., 1., nil}, Point{3., 3., nil}), 1),
			[]*Point{
				{1.5, 1.5, nil}, // SW I
				{1.2, 1.7, nil}, // SW
				{1.1, 1.0, nil}, // SW
				{1.8, 1.7, nil}, // SW I
				{2.2, 1.5, nil}, // SE I
				{2.6, 1.1, nil}, // SE
				{1.2, 2.3, nil}, // NW
				{1.6, 2.1, nil}, // NW I
			},
			NewArea(Point{1., 1., nil}, Point{2., 2., nil}),
			[]*Point{
				{1.5, 1.5, nil}, // SW I
				{1.2, 1.7, nil}, // SW
				{1.1, 1.0, nil}, // SW
				{1.8, 1.7, nil}, // SW I
			}},
		{"Test2_SE", NewNode(NewArea(Point{2., 1., nil}, Point{3., 2., nil}), 1),
			[]*Point{
				{1.5, 1.5, nil}, // SW I
				{1.2, 1.7, nil}, // SW
				{1.1, 1.0, nil}, // SW
				{1.8, 1.7, nil}, // SW I
				{2.2, 1.5, nil}, // SE I
				{2.6, 1.1, nil}, // SE
				{1.2, 2.3, nil}, // NW
				{1.6, 2.1, nil}, // NW I
			},
			NewArea(Point{2., 1., nil}, Point{3., 2., nil}),
			[]*Point{
				{2.2, 1.5, nil}, // SE I
				{2.6, 1.1, nil}, // SE
			}},
		{"Test2_NW", NewNode(NewArea(Point{1., 1., nil}, Point{3., 3., nil}), 1),
			[]*Point{
				{1.5, 1.5, nil}, // SW I
				{1.2, 1.7, nil}, // SW
				{1.1, 1.0, nil}, // SW
				{1.8, 1.7, nil}, // SW I
				{2.2, 1.5, nil}, // SE I
				{2.6, 1.1, nil}, // SE
				{1.2, 2.3, nil}, // NW
				{1.6, 2.1, nil}, // NW I
			},
			NewArea(Point{1., 2., nil}, Point{2., 3., nil}),
			[]*Point{
				{1.2, 2.3, nil}, // NW
				{1.6, 2.1, nil}, // NW I
			}},
		{"Test2_NE", NewNode(NewArea(Point{1., 1., nil}, Point{3., 3., nil}), 1),
			[]*Point{
				{1.5, 1.5, nil}, // SW I
				{1.2, 1.7, nil}, // SW
				{1.1, 1.0, nil}, // SW
				{1.8, 1.7, nil}, // SW I
				{2.2, 1.5, nil}, // SE I
				{2.6, 1.1, nil}, // SE
				{1.2, 2.3, nil}, // NW
				{1.6, 2.1, nil}, // NW I
			},
			NewArea(Point{2., 2., nil}, Point{3., 3., nil}),
			[]*Point{}},
		{"Test2_Inner", NewNode(NewArea(Point{1., 1., nil}, Point{3., 3., nil}), 1),
			[]*Point{
				{1.5, 1.5, nil}, // SW I
				{1.2, 1.7, nil}, // SW
				{1.1, 1.0, nil}, // SW
				{1.8, 1.7, nil}, // SW I
				{2.2, 1.5, nil}, // SE I
				{2.6, 1.1, nil}, // SE
				{1.2, 2.3, nil}, // NW
				{1.6, 2.1, nil}, // NW I
			},
			NewArea(Point{1.5, 1.5, nil}, Point{2.5, 2.5, nil}),
			[]*Point{
				{1.5, 1.5, nil}, // SW I
				{1.8, 1.7, nil}, // SW I
				{2.2, 1.5, nil}, // SE I
				{1.6, 2.1, nil}, // NW I
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.node.points = append(tt.node.points, tt.insertPoints...)
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
		{"TestIn1", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), &Point{1.5, 1.5, nil}, true},
		{"TestIn2", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), &Point{2.5, 2.5, nil}, true},
		{"TestIn3", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), &Point{1.5, 2.5, nil}, true},
		{"TestIn4", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), &Point{2.5, 1.5, nil}, true},

		{"TestOut1", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), &Point{-1.5, -1.5, nil}, false},
		{"TestOut2", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), &Point{4.5, 2.5, nil}, false},
		{"TestOut3", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), &Point{1.5, 4.5, nil}, false},
		{"TestOut4", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), &Point{2.5, -1.5, nil}, false},

		{"TestNeg1", NewArea(NewPoint(-1., -1.), NewPoint(-3., -3.)), &Point{-1.5, -1.5, nil}, true},
		{"TestNeg2", NewArea(NewPoint(-1., -1.), NewPoint(-3., -3.)), &Point{-2.5, -2.5, nil}, true},
		{"TestNeg3", NewArea(NewPoint(-1., -1.), NewPoint(-3., -3.)), &Point{-3.5, -2.5, nil}, false},
		{"TestNeg4", NewArea(NewPoint(-1., -1.), NewPoint(-3., -3.)), &Point{2.5, 1.5, nil}, false},

		{"TestBorder1", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), &Point{1., 1., nil}, true},
		{"TestBorder2", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), &Point{1., 2., nil}, true},
		{"TestBorder3", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), &Point{3., 2., nil}, false},
		{"TestBorder4", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), &Point{1., 3., nil}, false},
		{"TestBorder5", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), &Point{3., 3., nil}, false},
		{"TestBorder6", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), &Point{3., 1., nil}, false},
		{"TestBorder7", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), &Point{2., 3., nil}, false},
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
		{"TestEmpty", NewNode(NewArea(Point{1., 1., nil}, Point{3., 3., nil}), 1),
			[]*Point{}, &Point{2.2, 1.4, nil}, false,
			&Node{NewArea(Point{1., 1., nil}, Point{3., 3., nil}),
				[]*Point{{2.2, 1.4, nil}}, 1, nil}},

		{"TestNotEmpty", NewNode(NewArea(Point{1., 1., nil}, Point{3., 3., nil}), 2),
			[]*Point{{1.1, 1.4, nil}}, &Point{2.2, 1.4, nil}, false,
			&Node{NewArea(Point{1., 1., nil}, Point{3., 3., nil}),
				[]*Point{{1.1, 1.4, nil}, {2.2, 1.4, nil}}, 2, nil}},

		{"TestEqual", NewNode(NewArea(Point{1., 1., nil}, Point{3., 3., nil}), 2),
			[]*Point{{2.2, 1.4, nil}}, &Point{2.2, 1.4, nil}, true,
			&Node{NewArea(Point{1., 1., nil}, Point{3., 3., nil}),
				[]*Point{{2.2, 1.4, nil}}, 1, nil}},

		{"TestFull", NewNode(NewArea(Point{1., 1., nil}, Point{3., 3., nil}), 1),
			[]*Point{{2.2, 1.4, nil}}, &Point{1.4, 1.1, nil}, false,
			&Node{NewArea(Point{1., 1., nil}, Point{3., 3., nil}),
				nil, 2, []*Node{
					NewNode(NewArea(Point{1., 2., nil}, Point{2., 3., nil}), 1),
					NewNode(NewArea(Point{2., 2., nil}, Point{3., 3., nil}), 1),
					{NewArea(Point{1., 1., nil}, Point{2., 2., nil}), []*Point{{1.4, 1.1, nil}}, 1, nil},
					{NewArea(Point{2., 1., nil}, Point{3., 2., nil}), []*Point{{2.2, 1.4, nil}}, 1, nil},
				}}},

		{"TestFullEqual", NewNode(NewArea(Point{1., 1., nil}, Point{3., 3., nil}), 1),
			[]*Point{{2.2, 1.4, nil}}, &Point{2.2, 1.4, nil}, true,
			&Node{NewArea(Point{1., 1., nil}, Point{3., 3., nil}),
				[]*Point{{2.2, 1.4, nil}}, 1, nil}},
	}
	for _, tt := range tests {
		tt.node.points = append(tt.node.points, tt.points...)
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
		points   []*Point
		num      int
		children []*Node
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"Test1", fields{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
			[]*Point{{2., 2., nil}}, 1, nil}, true},
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
		insertPoints []*Point
		search       *Point
		want         *Point
	}{
		{"Test1_ExistingPoint", NewNode(NewArea(Point{1., 1., nil}, Point{3., 3., nil}), 1),
			[]*Point{
				{1.5, 1.5, "p"}, // SW I
				{1.2, 1.7, "p"}, // SW
				{1.1, 1.0, "p"}, // SW
				{1.8, 1.7, "p"}, // SW I
				{2.2, 1.5, "p"}, // SE I
				{2.6, 1.1, "p"}, // SE
				{1.2, 2.3, "p"}, // NW
				{1.6, 2.1, "p"}, // NW I
			},
			&Point{1.5, 1.5, nil},
			&Point{1.5, 1.5, "p"}},
		{"Test1_MissingPoint", NewNode(NewArea(Point{1., 1., nil}, Point{3., 3., nil}), 1),
			[]*Point{
				{1.5, 1.5, "p"}, // SW I
				{1.2, 1.7, "p"}, // SW
				{1.1, 1.0, "p"}, // SW
				{1.8, 1.7, "p"}, // SW I
				{2.2, 1.5, "p"}, // SE I
				{2.6, 1.1, "p"}, // SE
				{1.2, 2.3, "p"}, // NW
				{1.6, 2.1, "p"}, // NW I
			},
			&Point{1.4, 1.5, nil},
			nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, p := range tt.insertPoints {
				_ = tt.node.Insert(p)
			}
			if got := tt.node.Get(tt.search); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_split(t *testing.T) {
	type fields struct {
		area     *Area
		points   []*Point
		num      int
		children []*Node
	}
	tests := []struct {
		name   string
		fields fields
		want   *Node
	}{
		{"Test1", fields{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
			[]*Point{
				{1.5, 1.5, nil},
				{1.5, 2.5, nil},
				{2.5, 1.5, nil},
				{2.5, 2.5, nil},
			}, 4, nil},
			&Node{NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), nil, 4,
				[]*Node{
					{NewArea(NewPoint(1., 2.), NewPoint(2., 3.)), []*Point{{1.5, 2.5, nil}}, 1, nil},
					{NewArea(NewPoint(2., 2.), NewPoint(3., 3.)), []*Point{{2.5, 2.5, nil}}, 1, nil},
					{NewArea(NewPoint(1., 1.), NewPoint(2., 2.)), []*Point{{1.5, 1.5, nil}}, 1, nil},
					{NewArea(NewPoint(2., 1.), NewPoint(3., 2.)), []*Point{{2.5, 1.5, nil}}, 1, nil},
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
		{"TestIn1", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), &Point{1.5, 1.5, nil}, SW},
		{"TestIn2", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), &Point{2.5, 2.5, nil}, NE},
		{"TestIn3", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), &Point{1.5, 2.5, nil}, NW},
		{"TestIn4", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), &Point{2.5, 1.5, nil}, SE},

		{"TestOut1", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), &Point{-1.5, -1.5, nil}, SW},
		{"TestOut2", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), &Point{4.5, 2.5, nil}, NE},
		{"TestOut3", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), &Point{1.5, 4.5, nil}, NW},
		{"TestOut4", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), &Point{2.5, -1.5, nil}, SE},

		{"TestNeg1", NewArea(NewPoint(-1., -1.), NewPoint(-3., -3.)), &Point{-1.5, -1.5, nil}, NE},
		{"TestNeg2", NewArea(NewPoint(-1., -1.), NewPoint(-3., -3.)), &Point{-2.5, -2.5, nil}, SW},
		{"TestNeg3", NewArea(NewPoint(-1., -1.), NewPoint(-3., -3.)), &Point{-1.5, -2.5, nil}, SE},
		{"TestNeg4", NewArea(NewPoint(-1., -1.), NewPoint(-3., -3.)), &Point{-2.5, -1.5, nil}, NW},

		{"TestBorder1", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), &Point{2., 2., nil}, NE},
		{"TestBorder2", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), &Point{2., 1.5, nil}, SE},
		{"TestBorder3", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), &Point{1.5, 2., nil}, NW},
		{"TestBorder4", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), &Point{1.999, 1.999, nil}, SW},
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
		{"TestEq1", &Point{1., 1., nil}, &Point{1., 1., nil}, true},
		{"TestEq2", &Point{1., 1., "nil"}, &Point{1., 1., nil}, true},
		{"TestNeq1", &Point{1., 2., nil}, &Point{1., 1., nil}, false},
		{"TestNeq2", &Point{3., 1., nil}, &Point{1., 1., nil}, false},
		{"TestNeq3", &Point{3., 1., nil}, &Point{1., 3., nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pointA.Equals(tt.pointB); got != tt.want {
				t.Errorf("Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}
