package block_stm

import (
	"context"
	"fmt"
	"time"
)

// Executor fields are not mutated during execution.
type Executor struct {
	ctx        context.Context // context for cancellation
	scheduler  *Scheduler      // scheduler for task management
	txExecutor TxExecutor      // callback to actually execute a transaction
	mvMemory   *MVMemory       // multi-version memory for the executor

	// index of the executor, used for debugging output
	i int
}

func NewExecutor(
	ctx context.Context,
	scheduler *Scheduler,
	txExecutor TxExecutor,
	mvMemory *MVMemory,
	i int,
) *Executor {
	return &Executor{
		ctx:        ctx,
		scheduler:  scheduler,
		txExecutor: txExecutor,
		mvMemory:   mvMemory,
		i:          i,
	}
}

// Invariant `num_active_tasks`:
//   - `NextTask` increases it if returns a valid task.
//   - `TryExecute` and `NeedsReexecution` don't change it if it returns a new valid task to run,
//     otherwise it decreases it.
func (e *Executor) Run() {
	var kind TaskKind
	version := InvalidTxnVersion
	for !e.scheduler.Done() {
		if !version.Valid() {
			// check for cancellation
			select {
			case <-e.ctx.Done():
				return
			default:
			}

			version, kind = e.scheduler.NextTask()
			fmt.Printf("mm-Run[%d][%s]-NextTask:[%v]\n", e.i, time.Now(), version.Index)
			continue
		}

		switch kind {
		case TaskKindExecution:
			fmt.Printf("mm-Run[%d][%s]-TryExecute-bf:[%v]\n", e.i, time.Now(), version.Index)
			version, kind = e.TryExecute(version)
			fmt.Printf("mm-Run[%d][%s]-TryExecute-af:[%v]\n", e.i, time.Now(), version.Index)

		case TaskKindValidation:
			fmt.Printf("mm-Run[%d][%s]-NeedsReexecution-bf:[%v]\n", e.i, time.Now(), version.Index)
			version, kind = e.NeedsReexecution(version)
			fmt.Printf("mm-Run[%d][%s]-NeedsReexecution-af:[%v]\n", e.i, time.Now(), version.Index)
		}
	}
}

func (e *Executor) TryExecute(version TxnVersion) (TxnVersion, TaskKind) {
	e.scheduler.executedTxns.Add(1)
	view := e.execute(version.Index)
	wroteNewLocation := e.mvMemory.Record(version, view)
	return e.scheduler.FinishExecution(version, wroteNewLocation)
}

func (e *Executor) NeedsReexecution(version TxnVersion) (TxnVersion, TaskKind) {
	e.scheduler.validatedTxns.Add(1)
	valid := e.mvMemory.ValidateReadSet(version.Index)
	aborted := !valid && e.scheduler.TryValidationAbort(version)
	if aborted {
		e.mvMemory.ConvertWritesToEstimates(version.Index)
	}
	return e.scheduler.FinishValidation(version.Index, aborted)
}

func (e *Executor) execute(txn TxnIndex) *MultiMVMemoryView {
	view := e.mvMemory.View(txn)
	e.txExecutor(txn, view)
	return view
}
