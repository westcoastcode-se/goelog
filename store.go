package goelog

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"os"
	"sync/atomic"
)

// EventLog storage where one or more events can be found
type EventLog struct {
	Path      string
	offset    atomic.Int32
	factories factories
}

func (s *EventLog) writeTransaction(t *Tnx) error {
	if !s.offset.CompareAndSwap(t.startOffset, t.startOffset+int32(t.data.Len())) {
		return ErrorChangedOutsideTnx
	}

	f, err := os.OpenFile(s.Path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			f, err = os.OpenFile(s.Path, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			// write header
			_, err = f.Write([]byte{'G', 'E', 'L', '1'})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	defer f.Close()

	err = binary.Write(f, binary.LittleEndian, int32(t.data.Len()))
	if err != nil {
		log.Panicf("%e happened when trying to save transaction. The next time the application is starting the event log should repair itself automatically if needed", err)
	}

	_, err = t.data.WriteTo(f)
	if err != nil {
		log.Panicf("%e happened when trying to save transaction. The next time the application is starting the event log should repair itself automatically if needed", err)
	}
	return nil
}

// NewTransaction creates a new transaction that we can use to store events with
func (s *EventLog) NewTransaction() *Tnx {
	return newTnx(s)
}

// Subscribe for events and puts the result into the supplied channel
func (s *EventLog) Subscribe(ch EventStream) error {
	// TODO implement
	return nil
}

func readAllFileAsBytes(path string) ([]byte, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		// if the file doesn't exist then assume that this is an empty event log for now
		// todo should force the user to actually create the event log if it doesn't exist?
		if err == os.ErrNotExist {
			return []byte{}, nil
		}
		return nil, err
	}

	b, err := io.ReadAll(f)
	if err != nil {
		_ = f.Close()
		return nil, err
	}

	_ = f.Close()
	return b, nil
}

// SubscribeAndLoad subscribes for events and then load the event log immediately
// this method closes the supplied channel when all events are read, unless an error occurs
//
// TODO figure out if we want to keep the event stream and stream if we receive events from committed transactions
// TODO replace channel with a callback function?
func (s *EventLog) SubscribeAndLoad(ch EventStream) error {
	b, err := readAllFileAsBytes(s.Path)
	if err != nil {
		return err
	}
	if len(b) <= 4 {
		close(ch)
		return nil
	}

	// verify header
	if b[0] != 'G' || b[1] != 'E' || b[2] != 'L' || b[3] != '1' {
		return errors.New("unknown event log header")
	}
	b = b[4:]
	go func() {
		for len(b) > 0 {
			// Each transaction contains <i32 LEN><bytes DATA><byte END-TOKEN>
			b = readTransaction(s.factories, ch, b[:])
		}
		close(ch)
	}()

	return nil
}

func readTransaction(factories map[string]EventFactory, result EventStream, b []byte) []byte {
	// first, skip the header
	bytesReader := bytes.NewReader(b)

	// read transaction header
	var transactionLength int32
	_ = binary.Read(bytesReader, binary.LittleEndian, &transactionLength)
	b = b[4:]

	// check transaction end
	if b[transactionLength-1] != MarkerTransactionEnd {
		log.Panicf("missing transaction marker")
	}

	// read all events
	transactionReader := EventReader{
		r: bytes.NewBuffer(b[:transactionLength]),
	}
	for {
		name := transactionReader.ReadString()
		parser := factories[name]
		version := int(transactionReader.ReadInt32())
		val := parser.New()
		parser.Read(val, version, &transactionReader)
		result <- Event{
			Name:    name,
			Version: version,
			Data:    val,
		}

		if transactionReader.bytesLeft() == 1 {
			break
		}
	}

	return b[transactionLength:]
}
