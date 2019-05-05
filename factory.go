package sagamaton

import (
	"github.com/google/uuid"
	"github.com/looplab/fsm"
)

type defaultIDDispenser struct{}

func (d *defaultIDDispenser) DispenseID() ExecutionID {
	u := uuid.New()
	return ExecutionID(u.String())
}

// ExecutorFactory provides building Executor from sub-transactions with Finite State Machine.
type ExecutorFactory struct {
	registry    SagaRegistry
	idDispenser SagaExecutionIDDispenser
}

// NewDefaultExecutorFactory returns new ExecutorFactory object, which uses default id dispenser.
func NewDefaultExecutorFactory(reg SagaRegistry) *ExecutorFactory {
	return NewExecutorFactory(reg, &defaultIDDispenser{})
}

// NewExecutorFactory returns new ExecutorFactory object.
func NewExecutorFactory(reg SagaRegistry, dispenser SagaExecutionIDDispenser) *ExecutorFactory {
	return &ExecutorFactory{reg, dispenser}
}

// NewExecutor returns Executor object.
func (f *ExecutorFactory) NewExecutor(sig Signature) (*Executor, error) {
	saga, err := f.registry.Load(sig)
	if err != nil {
		return nil, err
	}
	id := f.idDispenser.DispenseID()
	stepper := newTxnStepper()
	for _, txn := range saga.SubTransactions() {
		stepper.Push(txn)
	}

	stateMachine := fsm.NewFSM(stepper.Initial(), stepper.Events(), stepper.Callbacks())
	return &Executor{id, stateMachine}, nil
}
