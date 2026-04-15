package battle

import (
	"reflect"
	"testing"
)

func Test_field_find(t *testing.T) {
	f := field{
		{2, 2, 2, 0, 1, 0, 0, 1, 1, 2},
		{2, 2, 2, 0, 0, 0, 0, 0, 2, 2},
		{2, 2, 2, 1, 0, 0, 0, 2, 4, 4},
		{2, 2, 2, 0, 0, 1, 1, 3, 4, 3},
		{2, 2, 2, 0, 0, 0, 0, 2, 4, 3},
		{4, 2, 2, 2, 2, 2, 2, 2, 4, 3},
		{3, 3, 2, 4, 4, 1, 3, 2, 4, 3},
		{4, 2, 4, 3, 4, 2, 2, 4, 2, 2},
		{4, 3, 4, 4, 4, 2, 2, 3, 3, 3},
		{4, 4, 4, 2, 2, 2, 2, 2, 2, 4},
	}
	tests := []struct {
		name string
		args []uint8
		want []point
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: nil,
			want: nil,
		},
		{
			name: "",
			args: []uint8{fieldFree},
			want: []point{3000, 5000, 6000, 3100, 4100, 5100, 6100, 7100, 4200, 5200, 6200, 3300, 4300, 3400, 4400, 5400, 6400},
		},
		{
			name: "",
			args: []uint8{fieldShip, fieldShot},
			want: []point{4010, 7010, 8010, 3210, 5310, 6310, 7330, 9330, 9430, 9530, 630, 1630, 5610, 6630, 9630, 3730, 1830, 7830, 8830, 9830},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := f.find(0, tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("find() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
