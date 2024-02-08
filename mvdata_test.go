package block_stm

import (
	"testing"

	"github.com/test-go/testify/require"
)

func TestMVData(t *testing.T) {
	data := NewMVData()

	// read closest version
	data.Write([]byte("a"), []byte("1"), TxnVersion{Index: 1, Incarnation: 1}, nil)
	data.Write([]byte("a"), []byte("2"), TxnVersion{Index: 2, Incarnation: 1}, nil)
	data.Write([]byte("a"), []byte("3"), TxnVersion{Index: 3, Incarnation: 1}, nil)
	data.Write([]byte("b"), []byte("2"), TxnVersion{Index: 2, Incarnation: 1}, nil)

	value, version, err := data.Read([]byte("a"), 4)
	require.NoError(t, err)
	require.Equal(t, Value([]byte("3")), value)
	require.Equal(t, TxnVersion{Index: 3, Incarnation: 1}, version)

	value, version, err = data.Read([]byte("a"), 3)
	require.NoError(t, err)
	require.Equal(t, Value([]byte("2")), value)
	require.Equal(t, TxnVersion{Index: 2, Incarnation: 1}, version)

	value, version, err = data.Read([]byte("b"), 3)
	require.NoError(t, err)
	require.Equal(t, Value([]byte("2")), value)
	require.Equal(t, TxnVersion{Index: 2, Incarnation: 1}, version)

}