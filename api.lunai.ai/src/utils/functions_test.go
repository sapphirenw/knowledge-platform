package utils

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestChunkStringEqualUntilN(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		n              int
		expectedOutput []string
	}{
		{
			name:           "empty string",
			input:          "",
			n:              10,
			expectedOutput: []string{},
		},
		{
			name:           "n larger than string length",
			input:          "hello",
			n:              10,
			expectedOutput: []string{"hello"},
		},
		{
			name:           "n equals string length",
			input:          "hello",
			n:              5,
			expectedOutput: []string{"hello"},
		},
		{
			name:           "n is one",
			input:          "hello",
			n:              1,
			expectedOutput: []string{"h", "e", "l", "l", "o"},
		},
		{
			name:           "n divides string length",
			input:          "hellohello",
			n:              5,
			expectedOutput: []string{"hello", "hello"},
		},
		{
			name:           "n does not divide string length, even distribution",
			input:          "hello world",
			n:              4, // Adjusting to ensure test reflects new requirement
			expectedOutput: []string{"hell", "o wo", "rld"},
		},
		{
			name:           "very long string",
			input:          "abcdefghijklmnopqrstuvwxyz",
			n:              4,
			expectedOutput: []string{"abcd", "efgh", "ijkl", "mnop", "qrst", "uvw", "xyz"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the updated splitString function that ensures even distribution
			result := ChunkStringEqualUntilN(tt.input, tt.n)
			if !reflect.DeepEqual(result, tt.expectedOutput) {
				t.Errorf("ChunkStringEqualUntilN(%q, %d) got %v, want %v", tt.input, tt.n, result, tt.expectedOutput)
			}
		})
	}
}

func TestConvertSlice(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	output := []float32{1, 2, 3, 4, 5}

	converted := ConvertSlice(input, func(i int) float32 { return float32(i) })

	require.Equal(t, output, converted)
}

func TestUUIDConvertion(t *testing.T) {
	gid := uuid.New()
	pid := GoogleUUIDToPGXUUID(gid)
	gid2 := PGXUUIDToGoogleUUID(pid)
	require.Equal(t, gid.String(), gid2.String())
}
