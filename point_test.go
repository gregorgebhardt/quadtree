package quadtree

import (
	"reflect"
	"testing"
)

func TestNewPoint(t *testing.T) {
	type args struct {
		x float64
		y float64
	}
	tests := []struct {
		name string
		args args
		want *Point[float64]
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

func TestPoint_equals(t *testing.T) {
	tests := []struct {
		name   string
		pointA *Point[float64]
		pointB *Point[float64]
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
			if got := tt.pointA.equals(tt.pointB); got != tt.want {
				t.Errorf("Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}
