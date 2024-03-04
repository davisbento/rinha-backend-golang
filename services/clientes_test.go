package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientHasEnoughLimit(t *testing.T) {
	tests := []struct {
		clientLimit int
		valor       int
		tipo        string
		expected    bool
	}{
		{clientLimit: 100, valor: 50, tipo: "c", expected: true},
		{clientLimit: 100, valor: 150, tipo: "d", expected: true},
		{clientLimit: 100, valor: 250, tipo: "d", expected: false},
		{clientLimit: 100, valor: 200, tipo: "d", expected: true},
		{clientLimit: 0, valor: 50, tipo: "d", expected: false},
		{clientLimit: -100, valor: 150, tipo: "d", expected: false},
		{clientLimit: 1000, valor: 2001, tipo: "d", expected: false},
		{clientLimit: 1000, valor: 2000, tipo: "d", expected: true},
	}

	for _, test := range tests {
		result, err := ClientHasEnoughLimit(test.clientLimit, test.valor, test.tipo)
		assert.NoError(t, err, "unexpected error")
		assert.Equal(t, test.expected, result,
			"for clientLimit=%d, valor=%d, tipo=%s, expected %t, but got %t",
			test.clientLimit, test.valor, test.tipo, test.expected, result)
	}
}
