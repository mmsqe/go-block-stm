package block_stm

import (
	"testing"

	"github.com/test-go/testify/require"
)

func BenchmarkBlockSTM(b *testing.B) {
	storage := NewMemDB()
	// blk := testBlock(10000, 10000)
	// blk := noConflictBlock(10000)
	blk := worstCaseBlock(10000)
	b.Run("worker-1", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			require.NoError(b, ExecuteBlock(storage, blk, 1))
		}
	})
	b.Run("worker-5", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			require.NoError(b, ExecuteBlock(storage, blk, 5))
		}
	})
	b.Run("worker-15", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			require.NoError(b, ExecuteBlock(storage, blk, 15))
		}
	})
	b.Run("worker-30", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			require.NoError(b, ExecuteBlock(storage, blk, 30))
		}
	})
}