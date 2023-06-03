package main

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
		err  error
	}{
		{
			name: "err invalid number of args",
			err:  errInvalidNumberOfArguments,
		},
		{
			name: "host is empty",
			args: []string{"", "1234"},
			err:  errHostIsEmpty,
		},
		{
			name: "port is empty",
			args: []string{"localhost", ""},
			err:  errPortIsEmpty,
		},
		{
			name: "ok",
			args: []string{"localhost", "1234"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address, err := parseArgs(tt.args)
			require.ErrorIs(t, err, tt.err)

			if tt.err == nil {
				require.Equal(t, net.JoinHostPort(tt.args[0], tt.args[1]), address)
			}
		})
	}
}
