package goelog

import (
	"errors"
	"path"
)

var (
	ErrorChangedOutsideTnx = errors.New("changes made outside the current transaction")
)

// EventStream represents a stream of events
type EventStream chan Event

type Event struct {
	Name    string
	Version int
	Data    interface{}
}

type factories map[string]EventFactory

type Repository struct {
	RootPath  string
	factories factories
}

// OpenEventLog opens an event log with the given name
func (d *Repository) OpenEventLog(name string) *EventLog {
	store := &EventLog{
		Path:      path.Join(d.RootPath, name+".events"),
		factories: d.factories,
	}
	// TODO use the MarkerTransactionEnd to truncate the binary storage until we get the recent and complete
	//      transaction
	return store
}

// NewRepository creates a new repository and configures with a list of factories
func NewRepository(rootPath string, factories []EventFactory) *Repository {
	db := &Repository{
		RootPath:  rootPath,
		factories: make(map[string]EventFactory),
	}
	for _, f := range factories {
		db.factories[f.Name()] = f
	}
	return db
}
