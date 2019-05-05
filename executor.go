package sagamaton

import "github.com/looplab/fsm"

// ExecutionID provides identifier for Executor.
type ExecutionID string

// Executor provides state machine runner.
type Executor struct {
	id  ExecutionID
	fsm *fsm.FSM
}
