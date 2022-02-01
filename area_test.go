package quadtree

import (
	"reflect"
	"testing"
)

func TestArea_NewArea(t *testing.T) {
	tests := []struct {
		name string
		l    Point
		u    Point
	}{
		{"Test1", Point{x: 1., y: 1.}, Point{x: 3., y: 3.}},
		{"Test2", Point{x: 1., y: 3.}, Point{x: 3., y: 1.}},
		{"Test3", Point{x: 3., y: 1.}, Point{x: 1., y: 3.}},
		{"Test4", Point{x: 3., y: 3.}, Point{x: 1., y: 1.}},
	}
	area := NewArea(Point{x: 1., y: 1.}, Point{x: 3., y: 3.})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t_area := NewArea(tt.l, tt.u)
			if !reflect.DeepEqual(area, t_area) {
				t.Errorf("area = %v, NewArea: %v", area, t_area)
			}
		})
	}
}

func TestArea_contains(t *testing.T) {
	type args struct {
		x float64
		y float64
	}
	tests := []struct {
		name string
		area *Area
		args args
		want bool
	}{
		{"Test1", NewArea(Point{x: 1., y: 1.}, Point{x: 3., y: 3.}), args{2., 2.}, true},
		{"Test2", NewArea(Point{x: 1., y: 1.}, Point{x: 3., y: 3.}), args{2., 4.}, false},
		{"Test3", NewArea(Point{x: 1., y: 1.}, Point{x: 3., y: 3.}), args{4., 2.}, false},
		{"Test4", NewArea(Point{x: 1., y: 1.}, Point{x: 3., y: 3.}), args{4., 4.}, false},
		{"Test5", NewArea(Point{x: 1., y: 1.}, Point{x: 3., y: 3.}), args{-2., 2.}, false},
		{"Test6", NewArea(Point{x: 1., y: 1.}, Point{x: 3., y: 3.}), args{2., -2.}, false},
		{"Test7", NewArea(Point{x: 1., y: 1.}, Point{x: 3., y: 3.}), args{-2., -2.}, false},
		{"Test8", NewArea(Point{x: 1., y: 1.}, Point{x: 3., y: 3.}), args{1., 1.}, true},
		{"Test9", NewArea(Point{x: 1., y: 1.}, Point{x: 3., y: 3.}), args{1., 2.}, true},
		{"Test10", NewArea(Point{x: 1., y: 1.}, Point{x: 3., y: 3.}), args{3., 3.}, true},
		{"Test11", NewArea(Point{x: 1., y: 1.}, Point{x: 3., y: 3.}), args{3., 1.}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.area.contains(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArea_intersects(t *testing.T) {
	tests := []struct {
		name  string
		area  *Area
		other *Area
		want  bool
	}{
		{"TestSame", NewArea(Point{x: 1., y: 1.}, Point{x: 3., y: 3.}),
			NewArea(Point{x: 1., y: 1.}, Point{x: 3., y: 3.}), true},

		{"TestOnePointOut", NewArea(Point{x: 1., y: 1.}, Point{x: 3., y: 3.}),
			NewArea(Point{x: 0., y: 1.}, Point{x: 3., y: 3.}), true},
		{"TestOnePointOut", NewArea(Point{x: 1., y: 1.}, Point{x: 3., y: 3.}),
			NewArea(Point{x: 1., y: 1.}, Point{x: 4., y: 3.}), true},
		{"TestOnePointOut", NewArea(Point{x: 1., y: 1.}, Point{x: 3., y: 3.}),
			NewArea(Point{x: 1., y: 1.}, Point{x: 3., y: 4.}), true},
		{"TestOnePointOut", NewArea(Point{x: 1., y: 1.}, Point{x: 3., y: 3.}),
			NewArea(Point{x: 1., y: 0.}, Point{x: 3., y: 3.}), true},
		{"TestOnePointOut", NewArea(Point{x: 0., y: 1.}, Point{x: 3., y: 3.}),
			NewArea(Point{x: 1., y: 1.}, Point{x: 3., y: 3.}), true},
		{"TestOnePointOut", NewArea(Point{x: 1., y: 1.}, Point{x: 4., y: 3.}),
			NewArea(Point{x: 1., y: 1.}, Point{x: 3., y: 3.}), true},
		{"TestOnePointOut", NewArea(Point{x: 1., y: 1.}, Point{x: 3., y: 4.}),
			NewArea(Point{x: 1., y: 1.}, Point{x: 3., y: 3.}), true},
		{"TestOnePointOut", NewArea(Point{x: 1., y: 0.}, Point{x: 3., y: 3.}),
			NewArea(Point{x: 1., y: 1.}, Point{x: 3., y: 3.}), true},

		{"TestIn", NewArea(Point{x: 1.5, y: 1.5}, Point{x: 2.3, y: 2.3}),
			NewArea(Point{x: 1., y: 1.}, Point{x: 3., y: 3.}), true},
		{"TestIn", NewArea(Point{x: 1., y: 1.}, Point{x: 3., y: 3.}),
			NewArea(Point{x: 1.5, y: 1.5}, Point{x: 2.3, y: 2.3}), true},

		{"TestOut", NewArea(Point{x: 1., y: 1.}, Point{x: 3., y: 3.}),
			NewArea(Point{x: 1., y: 4.}, Point{x: 3., y: 3.}), true},
		{"TestOut", NewArea(Point{x: 1., y: 1.}, Point{x: 3., y: 3.}),
			NewArea(Point{x: 4., y: 1.}, Point{x: 3., y: 3.}), true},
		{"TestOut", NewArea(Point{x: 1., y: 1.}, Point{x: 3., y: 3.}),
			NewArea(Point{x: 1., y: 1.}, Point{x: .5, y: 3.}), true},
		{"TestOut", NewArea(Point{x: 1., y: 1.}, Point{x: 3., y: 3.}),
			NewArea(Point{x: 1., y: 1.}, Point{x: 3., y: -3.}), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.area.intersects(tt.other); got != tt.want {
				t.Errorf("intersects() = %v, want %v", got, tt.want)
			}
		})
	}
}
