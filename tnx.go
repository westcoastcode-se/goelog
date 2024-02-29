package goelog

import (
	"bytes"
	"encoding/binary"
	"log"
	"reflect"
)

const (
	MarkerTransactionEnd = byte(4)
)

// Tnx represents a transaction against the actual EventLog storage
type Tnx struct {
	data        bytes.Buffer
	startOffset int32
	store       *EventLog
}

type EventItem struct {
	i    interface{}
	name string
}

// EventFor converts an event into an item that can be saved int the EventLog storage
func EventFor[T any](t *T) *EventItem {
	return &EventItem{
		i:    t,
		name: reflect.TypeFor[T]().Name(),
	}
}

// Append the supplied event to the transaction
func (t *Tnx) Append(e *EventItem) {
	nameLen := int32(len(e.name))
	err := binary.Write(&t.data, binary.LittleEndian, nameLen)
	if err != nil {
		log.Panicf("failed to write binary data. %e", err)
	}
	t.data.Write([]byte(e.name))

	factory := t.store.factories[e.name]
	err = binary.Write(&t.data, binary.LittleEndian, int32(factory.Version()))
	if err != nil {
		log.Panicf("failed to write binary data. %e", err)
	}
	factory.Write(e.i, &EventWriter{&t.data})
}

// Commit the transaction and save the content to the disk. This will return ErrorChangedOutsideTnx if another
// transaction managed to commit it's content before you.
//
// TODO add support for "events" where the order doesn't really matter. Those cases should be allowed to be committed
func (t *Tnx) Commit() error {
	// nothing to commit
	if t.data.Len() == 0 {
		return nil
	}

	// put the transaction end marker at the end of the transaction binary data. This is used by the
	// database when verifying that the transaction isn't saved half-done. Any transaction, at the end, that doesn't have
	// this marker is automatically discarded
	t.data.WriteByte(MarkerTransactionEnd)

	// write the data back to the EventLog storage
	return t.store.writeTransaction(t)
}

// newTnx open a new transaction
func newTnx(s *EventLog) *Tnx {
	t := &Tnx{
		data:        bytes.Buffer{},
		startOffset: s.offset.Load(),
		store:       s,
	}
	return t
}
