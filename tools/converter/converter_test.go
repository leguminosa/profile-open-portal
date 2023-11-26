package converter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToInt(t *testing.T) {
	tests := []struct {
		name string
		v    interface{}
		want int
	}{
		{
			name: "int",
			v:    1,
			want: 1,
		},
		{
			name: "int8",
			v:    int8(1),
			want: 1,
		},
		{
			name: "int16",
			v:    int16(1),
			want: 1,
		},
		{
			name: "int32",
			v:    int32(1),
			want: 1,
		},
		{
			name: "int64",
			v:    int64(1),
			want: 1,
		},
		{
			name: "float32",
			v:    float32(1),
			want: 1,
		},
		{
			name: "float64",
			v:    float64(1),
			want: 1,
		},
		{
			name: "string",
			v:    "1",
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToInt(tt.v)
			assert.Equal(t, tt.want, got)
		})
	}
}
