package goelog

import "reflect"

type TFactoryWriter[T any] func(t *T, w *EventWriter)

type TFactoryReader[T any] func(t *T, version int, r *EventReader)

// EventFactory is, primarily used during the serialization and deserialization process
type EventFactory interface {
	New() interface{}
	Name() string
	Version() int
	Write(i interface{}, w *EventWriter)
	Read(i interface{}, version int, r *EventReader)
}

// TEventFactory is a generics based implementation of the EventFactory interface
type TEventFactory[T any] struct {
	name    string
	version int
	writer  TFactoryWriter[T]
	reader  TFactoryReader[T]
}

func (f *TEventFactory[T]) New() interface{} {
	return new(T)
}

func (f *TEventFactory[T]) Name() string {
	return f.name
}

func (f *TEventFactory[T]) Version() int {
	return f.version
}

func (f *TEventFactory[T]) Write(i interface{}, w *EventWriter) {
	c, _ := any(i).(*T)
	f.writer(c, w)
}

func (f *TEventFactory[T]) Read(i interface{}, version int, r *EventReader) {
	c, _ := any(i).(*T)
	f.reader(c, version, r)
}

// nameOf helps us figure out the type name based on the T
func nameOf[T any](p T) string {
	tp := reflect.TypeOf(p)
	if tp.Kind() == reflect.Struct {
		return tp.Name()
	} else {
		return tp.Elem().Name()
	}
}

// FactoryOf creates an event factory based the supplied first type.
//
// The first argument v is an instance of the struct we want to create the factory for.
//
// The second argument version is used during the serialization process so that we can keep track of
// multiple versions of the same model.
//
// The third and fourth arguments are used for reading and writing the actual binary data from the storage file
func FactoryOf[T any](v T, version int, reader TFactoryReader[T], writer TFactoryWriter[T]) EventFactory {
	return &TEventFactory[T]{
		name:   nameOf[T](v),
		reader: reader,
		writer: writer,
	}
}
