package goelog

import (
	"bytes"
	"encoding/binary"
	"log"
	"time"
)

type EventWriter struct {
	buf *bytes.Buffer
}

// WriteString writes the supplied string to the storage
func (w *EventWriter) WriteString(val string) {
	bb := []byte(val)
	err := binary.Write(w.buf, binary.LittleEndian, int32(len(bb)))
	if err != nil {
		log.Panicf("failed to write binary data. %e", err)
	}
	_, err = w.buf.Write(bb)
	if err != nil {
		log.Panicf("failed to write binary data. %e", err)
	}
}

// WriteInt32 writes the supplied 32-bit integer to the storage
func (w *EventWriter) WriteInt32(val int32) {
	err := binary.Write(w.buf, binary.LittleEndian, val)
	if err != nil {
		log.Panicf("failed to write binary data. %e", err)
	}
}

// WriteFloat32 writes the supplied 32-bit decimal to the storage
func (w *EventWriter) WriteFloat32(val float32) {
	err := binary.Write(w.buf, binary.LittleEndian, val)
	if err != nil {
		log.Panicf("failed to write binary data. %e", err)
	}
}

// WriteFloat64 writes the supplied 64-bit decimal to the storage
func (w *EventWriter) WriteFloat64(val float64) {
	err := binary.Write(w.buf, binary.LittleEndian, val)
	if err != nil {
		log.Panicf("failed to write binary data. %e", err)
	}
}

// WriteTime writes the supplied time.Time object to the storage. The value is guaranteed to be in UTC before being saved
func (w *EventWriter) WriteTime(val time.Time) {
	err := binary.Write(w.buf, binary.LittleEndian, val.UTC().UnixMilli())
	if err != nil {
		log.Panicf("failed to write binary data. %e", err)
	}
}
