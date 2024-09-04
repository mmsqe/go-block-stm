package block_stm

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"

	storetypes "cosmossdk.io/store/types"
)

func ExecuteBlock(
	ctx context.Context,
	blockSize int,
	stores map[storetypes.StoreKey]int,
	storage MultiStore,
	executors int,
	txExecutor TxExecutor,
	t time.Time,
	h int64,
) error {
	if executors < 0 {
		return fmt.Errorf("invalid number of executors: %d", executors)
	}
	if executors == 0 {
		executors = runtime.NumCPU()
	}

	// Create a new scheduler
	scheduler := NewScheduler(blockSize)
	mvMemory := NewMVMemory(blockSize, stores, storage, scheduler)

	var wg sync.WaitGroup
	wg.Add(executors)
	for i := 0; i < executors; i++ {
		e := NewExecutor(ctx, scheduler, txExecutor, mvMemory, i)
		go func() {
			defer wg.Done()
			e.Run(t, h)
		}()
	}
	wg.Wait()

	if !scheduler.Done() {
		if ctx.Err() != nil {
			// cancelled
			return ctx.Err()
		}

		return errors.New("scheduler did not complete")
	}

	// Write the snapshot into the storage
	mvMemory.WriteSnapshot(storage)
	return nil
}
