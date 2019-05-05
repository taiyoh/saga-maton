package sagamaton

// Action provides function for each state.
type Action func([]interface{}) error

// SubTransaction provides Action binding when forward and compensating.
type SubTransaction struct {
	name         string
	forward      Action
	compensating Action
}

// Name returns sub-transaction name.
func (s SubTransaction) Name() string {
	return s.name
}

// ForwardAction returns Action for forward operation.
func (s SubTransaction) ForwardAction() Action {
	return s.forward
}

// CompensatingAction returns Action for compensating operation.
func (s SubTransaction) CompensatingAction() Action {
	return s.compensating
}
