package sagamaton

// SagaRegistry provides interface for Saga transactions storage.
type SagaRegistry interface {
	Store(*Saga) error
	Load(Signature) (*Saga, error)
}

// SagaExecutionIDDispenser provides interface for Executor id dispenser.
type SagaExecutionIDDispenser interface {
	DispenseID() ExecutionID
}
