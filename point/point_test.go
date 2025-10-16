package point

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name string
		args args
		want P
	}{
		{
			name: "Create point (3,4)",
			args: args{x: 3, y: 4},
			want: P{X: 3, Y: 4},
		},
		{
			name: "Create point (-1,0)",
			args: args{x: -1, y: 0},
			want: P{X: -1, Y: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.x, tt.args.y); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestP_Equal(t *testing.T) {
	type args struct {
		other P
	}
	tests := []struct {
		name string
		p    P
		args args
		want bool
	}{
		{
			name: "Equal points",
			p:    P{X: 1, Y: 2},
			args: args{other: P{X: 1, Y: 2}},
			want: true,
		},
		{
			name: "Unequal points",
			p:    P{X: 1, Y: 2},
			args: args{other: P{X: 2, Y: 3}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Equal(tt.args.other); got != tt.want {
				t.Errorf("P.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestP_Neighboring(t *testing.T) {
	type args struct {
		other P
	}
	tests := []struct {
		name string
		p    P
		args args
		want bool
	}{
		{
			name: "Neighboring points (horizontal)",
			p:    P{X: 1, Y: 2},
			args: args{other: P{X: 2, Y: 2}},
			want: true,
		},
		{
			name: "Neighboring points (vertical)",
			p:    P{X: 1, Y: 2},
			args: args{other: P{X: 1, Y: 3}},
			want: true,
		},
		{
			name: "Neighboring points (diagonal)",
			p:    P{X: 1, Y: 2},
			args: args{other: P{X: 2, Y: 3}},
			want: true,
		},
		{
			name: "Non-neighboring points",
			p:    P{X: 1, Y: 2},
			args: args{other: P{X: 3, Y: 3}},
			want: false,
		},
		{
			name: "Same point",
			p:    P{X: 1, Y: 2},
			args: args{other: P{X: 1, Y: 2}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Neighboring(tt.args.other); got != tt.want {
				t.Errorf("P.Neighboring() = %v, want %v", got, tt.want)
			}
		})
	}
}
