# Golang Event Log

A filebased event log that can be used for, for example, event sourcing

## Example

Based on `example/main.go`

```go
package main

import (
	"log"
)

type CreateUser struct {
	Name string
}

func CreateUserRead(e *CreateUser, version int, r *goelog.EventReader) {
	e.Name = r.ReadString()
}

func CreateUserWrite(e *CreateUser, w *goelog.EventWriter) {
	w.WriteString(e.Name)
}

func main() {
	// Prepare the actual factories used when serializing and deserializing the event structs.
	//
	// We are using a generic-type solution to figure out a factory base
	types := []goelog.EventFactory{
		goelog.FactoryOf(CreateUser{}, 0, CreateUserRead, CreateUserWrite),
	}

	// Create a repository where one or more event logs can be found
	d := goelog.NewRepository("", types)

	// Open the storage so that we can read and write data from it
	s := d.OpenEventLog("users")

	// Open a new transaction and add an event to the event log
	t := s.NewTransaction()
	t.Append(goelog.EventFor(&CreateUser{Name: "admin"}))
	if e := t.Commit(); e != nil {
		log.Panicf("failed to commit transaction. %e", e)
	}

	// Create a stream and load all events from the store that's already saved into it
	ch := make(goelog.EventStream)
	if e := s.SubscribeAndLoad(ch); e != nil {
		log.Panicf("failed to subscribe to stream. %e", e)
	}

	// Read data from the channel and do stuff with each event
	for event := range ch {
		switch data := event.Data.(type) {
		case *CreateUser:
			log.Printf("user %s is created", data.Name)
		}
	}
}
```
