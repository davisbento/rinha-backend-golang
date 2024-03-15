package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientHasEnoughLimit(t *testing.T) {
	tests := []struct {
		clientLimit int
		saldo       int
		valor       int
		tipo        string
		expected    bool
	}{
		{100, 100, 100, "c", true},
		{100, 100, -100, "d", false},
		{100, 100, -101, "d", false},
		{500, 600, -400, "d", false},
		{500, 600, 400, "c", true},
		{1000, 900, -100, "d", false},
	}

	for _, test := range tests {
		result := ClientHasEnoughLimit(test.clientLimit, test.saldo, test.valor, test.tipo)
		assert.Equal(t, test.expected, result,
			"for clientLimit=%d, valor=%d, tipo=%s, expected %t, but got %t",
			test.clientLimit, test.valor, test.tipo, test.expected, result)
	}
}
