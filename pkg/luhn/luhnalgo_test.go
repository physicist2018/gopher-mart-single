package luhn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateLuhnNumber(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{
			name:   "length 2",
			length: 2,
		},
		{
			name:   "length 3",
			length: 3,
		},
		{
			name:   "length 16",
			length: 16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateLuhnNumber(tt.length)
			println(got)
			valid := ValidateLuhnNumber(got)
			assert.True(t, valid)
		})
	}
}

func TestValidateLuhnNumber(t *testing.T) {
	tests := []struct {
		name   string
		number string
		want   bool
	}{
		{
			name:   "valid number",
			number: "0552199341654035",
			want:   true,
		},
		{
			name:   "invalid number",
			number: "0552199341654036",
			want:   false,
		},
		{
			name:   "number contains non-digit characters",
			number: "0552199341654035a",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateLuhnNumber(tt.number)
			println(got)
			assert.Equal(t, tt.want, got)
		})
	}
}
