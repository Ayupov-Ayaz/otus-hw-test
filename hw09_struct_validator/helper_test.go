package hw09structvalidator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_castToString(t *testing.T) {
	const exp = "12345fd"

	type str string

	tests := []interface{}{
		str(exp),
	}

	for _, tt := range tests {
		require.Equal(t, exp, castToString(tt))
	}
}
