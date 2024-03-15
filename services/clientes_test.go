package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientHasEnoughLimit(t *testing.T) {
	tests := []struct {
		clientLimit int
		saldo       int
		tipo        string
		expected    bool
	}{
		{100, 50, "c", true},
		{100, 50, "d", true},
		{100, 100, "d", true},
		{100, 100, "c", true},
		{100, 101, "d", true},
		{200, -201, "d", false},
		{200, -200, "d", true},
	}

	for _, test := range tests {
		result := ClientHasEnoughLimit(test.clientLimit, test.saldo, test.tipo)
		assert.Equal(t, test.expected, result,
			"for clientLimit=%d, valor=%d, tipo=%s, expected %t, but got %t",
			test.clientLimit, test.tipo, test.expected, result)
	}
}
