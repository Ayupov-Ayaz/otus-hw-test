package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/ayupov-ayaz/otus-hw-test/hw11_telnet_client/mocks"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})
}

type mock struct {
	reader *mocks.MockReadCloser
	writer *mocks.MockWriter
	conn   *mocks.MockConn
}

func newMock(ctrl *gomock.Controller) *mock {
	return &mock{
		reader: mocks.NewMockReadCloser(ctrl),
		writer: mocks.NewMockWriter(ctrl),
		conn:   mocks.NewMockConn(ctrl),
	}
}

func (m *mock) expectConnRead() *gomock.Call {
	return m.conn.EXPECT().Read(gomock.Any()).Times(1)
}

func (m *mock) expectRead() *gomock.Call {
	return m.reader.EXPECT().Read(gomock.Any()).Times(1)
}

func (m *mock) expectWrite() *gomock.Call {
	return m.writer.EXPECT().Write(gomock.Any()).Times(1)
}

func TestTelnetClient_Receive(t *testing.T) {
	errRead := errors.New("read failed")

	tests := []struct {
		name     string
		before   func(m *mock)
		expected string
		err      error
	}{
		{
			name: "failed to conn read",
			err:  errRead,
			before: func(m *mock) {
				m.expectConnRead().Return(0, errRead)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := newMock(ctrl)
			tt.before(m)

			client, ok := NewTelnetClient("", 0, m.reader, m.writer).(*telnetClient)
			require.True(t, ok)
			client.conn = m.conn
			require.ErrorIs(t, client.Receive(), tt.err)
		})
	}
}

func TestTelnetClient_Send(t *testing.T) {
	errRead := errors.New("read failed")

	tests := []struct {
		name     string
		before   func(m *mock)
		expected string
		err      error
	}{
		{
			name: "failed to read",
			err:  errRead,
			before: func(m *mock) {
				m.expectRead().Return(0, errRead)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := newMock(ctrl)
			tt.before(m)

			client, ok := NewTelnetClient("", 0, m.reader, m.writer).(*telnetClient)
			require.True(t, ok)
			client.conn = m.conn
			require.ErrorIs(t, client.Send(), tt.err)
		})
	}
}
