package sagamaton

import (
	"fmt"

	"github.com/looplab/fsm"
)

type txnIndex int

const (
	prev txnIndex = iota
	current
	next
)

type txnStepper struct {
	txns         [3]*SubTransaction
	events       fsm.Events
	callbacks    fsm.Callbacks
	initialized  bool
	initialState string
}

func newTxnStepper() *txnStepper {
	return &txnStepper{
		[3]*SubTransaction{nil, nil, nil},
		fsm.Events{},
		fsm.Callbacks{},
		false,
		"",
	}
}

func (s *txnStepper) Push(txn *SubTransaction) {
	s.txns = [3]*SubTransaction{
		s.txns[current],
		s.txns[next],
		txn,
	}
	if !s.initialized {
		s.initialize(txn)
		return
	}
	s.forwardState()
	s.compensatingState()
}

func (s *txnStepper) initialize(txn *SubTransaction) {
	s.initialState = s.forwardName(txn)
	s.initialized = true
}

func (s *txnStepper) Initial() string {
	return s.initialState
}

func (s *txnStepper) forwardName(txn *SubTransaction) string {
	return fmt.Sprintf("%s_forward", txn.Name())
}

func (s *txnStepper) compensatingName(txn *SubTransaction) string {
	return fmt.Sprintf("%s_compensating", txn.Name())
}

func (s *txnStepper) forwardState() {
	txn := s.txns[current]
	action := txn.ForwardAction()
	s.events = append(s.events, s.forwardEventDesc())
	s.callbacks[fmt.Sprintf("enter_%s", s.forwardName(txn))] = func(ev *fsm.Event) {
		err := action(ev.Args)
		if err == nil {
			return
		}
		ev.Dst = s.compensatingName(txn)
		ev.Err = err
	}
}

func (s *txnStepper) compensatingState() {
	txn := s.txns[current]
	action := txn.CompensatingAction()
	s.events = append(s.events, s.compensatingEventDesc())
	s.callbacks[fmt.Sprintf("enter_%s", s.compensatingName(txn))] = func(ev *fsm.Event) {
		err := action(ev.Args)
		ev.Err = err
	}
}

func (s *txnStepper) forwardEventDesc() fsm.EventDesc {
	return fsm.EventDesc{
		Name: s.forwardName(s.txns[current]),
		Src:  s.forwardSrc(),
		Dst:  s.forwardName(s.txns[next]),
	}
}

func (s *txnStepper) forwardSrc() []string {
	if txn := s.txns[prev]; txn != nil {
		return []string{txn.Name()}
	}
	return []string{}
}

func (s *txnStepper) compensatingEventDesc() fsm.EventDesc {
	return fsm.EventDesc{
		Name: s.compensatingName(s.txns[current]),
		Src:  []string{s.forwardName(s.txns[current])},
		Dst:  s.compensatingName(s.txns[prev]), // TODO nilの時対応
	}
}

func (s *txnStepper) Events() fsm.Events {
	return s.events
}

func (s *txnStepper) Callbacks() fsm.Callbacks {
	return s.callbacks
}
