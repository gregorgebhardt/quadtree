package quadtree

import (
	"reflect"
	"testing"
)

func TestArea_NewArea(t *testing.T) {
	tests := []struct {
		name string
		l    *Point
		u    *Point
	}{
		{"Test1", NewPoint(1., 1.), NewPoint(3., 3.)},
		{"Test2", NewPoint(1., 3.), NewPoint(3., 1.)},
		{"Test3", NewPoint(3., 1.), NewPoint(1., 3.)},
		{"Test4", NewPoint(3., 3.), NewPoint(1., 1.)},
	}
	area := NewArea(NewPoint(1., 1.), NewPoint(3., 3.))
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
		{"Test1", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), args{2., 2.}, true},
		{"Test2", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), args{2., 4.}, false},
		{"Test3", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), args{4., 2.}, false},
		{"Test4", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), args{4., 4.}, false},
		{"Test5", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), args{-2., 2.}, false},
		{"Test6", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), args{2., -2.}, false},
		{"Test7", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), args{-2., -2.}, false},
		{"Test8", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), args{1., 1.}, true},
		{"Test9", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), args{1., 2.}, true},
		{"Test10", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), args{2., 1.}, true},
		{"Test11", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), args{3., 3.}, false},
		{"Test12", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), args{3., 1.}, false},
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
		{"TestSame", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
			NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), true},

		{"TestOnePointOut", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
			NewArea(NewPoint(0., 1.), NewPoint(3., 3.)), true},
		{"TestOnePointOut", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
			NewArea(NewPoint(1., 1.), NewPoint(4., 3.)), true},
		{"TestOnePointOut", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
			NewArea(NewPoint(1., 1.), NewPoint(3., 4.)), true},
		{"TestOnePointOut", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
			NewArea(NewPoint(1., 0.), NewPoint(3., 3.)), true},
		{"TestOnePointOut", NewArea(NewPoint(0., 1.), NewPoint(3., 3.)),
			NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), true},
		{"TestOnePointOut", NewArea(NewPoint(1., 1.), NewPoint(4., 3.)),
			NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), true},
		{"TestOnePointOut", NewArea(NewPoint(1., 1.), NewPoint(3., 4.)),
			NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), true},
		{"TestOnePointOut", NewArea(NewPoint(1., 0.), NewPoint(3., 3.)),
			NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), true},

		{"TestIn", NewArea(NewPoint(1.5, 1.5), NewPoint(2.3, 2.3)),
			NewArea(NewPoint(1., 1.), NewPoint(3., 3.)), true},
		{"TestIn", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
			NewArea(NewPoint(1.5, 1.5), NewPoint(2.3, 2.3)), true},

		{"TestOut", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
			NewArea(NewPoint(1., 4.), NewPoint(3., 3.)), true},
		{"TestOut", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
			NewArea(NewPoint(4., 1.), NewPoint(3., 3.)), true},
		{"TestOut", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
			NewArea(NewPoint(1., 1.), NewPoint(.5, 3.)), true},
		{"TestOut", NewArea(NewPoint(1., 1.), NewPoint(3., 3.)),
			NewArea(NewPoint(1., 1.), NewPoint(3., -3.)), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.area.intersects(tt.other); got != tt.want {
				t.Errorf("intersects() = %v, want %v", got, tt.want)
			}
		})
	}
}
