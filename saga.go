package sagamaton

// Signature represents key for finding saga.
type Signature string

// Saga provides binding sub-transactions for event driven tasks.
type Saga struct {
	signature Signature
	txns      []*SubTransaction
}

// NewSaga returns Saga object.
func NewSaga(sig Signature, txns ...*SubTransaction) *Saga {
	if txns == nil {
		txns = []*SubTransaction{}
	}
	return &Saga{sig, txns}
}

// Signature returns Signature.
func (s *Saga) Signature() Signature {
	return s.signature
}

// SubTransactions returns sub-transactions list.
func (s *Saga) SubTransactions() []*SubTransaction {
	return s.txns
}
