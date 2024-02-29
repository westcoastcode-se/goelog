package goelog

import (
	"bytes"
	"encoding/binary"
	"log"
	"time"
)

type EventReader struct {
	r *bytes.Buffer
}

// ReadString reads a string from the storage
func (r *EventReader) ReadString() string {
	var strLen int32 = 0
	err := binary.Read(r.r, binary.LittleEndian, &strLen)
	if err != nil {
		log.Panicf("failed to read binary data: %e", err)
	}

	var b = make([]byte, strLen)
	_, err = r.r.Read(b)
	if err != nil {
		log.Panicf("failed to read binary data: %e", err)
	}

	return string(b)
}

// ReadInt32 reads a 32-bit integer from the storage
func (r *EventReader) ReadInt32() int32 {
	var value int32
	err := binary.Read(r.r, binary.LittleEndian, &value)
	if err != nil {
		log.Panicf("failed to read binary data: %e", err)
	}

	return value
}

// ReadFloat32 reads a 32-bit decimal from the storage
func (r *EventReader) ReadFloat32() float32 {
	var value float32
	err := binary.Read(r.r, binary.LittleEndian, &value)
	if err != nil {
		log.Panicf("failed to read binary data: %e", err)
	}

	return value
}

// ReadFloat64 reads a 64-bit decimal from the storage
func (r *EventReader) ReadFloat64() float64 {
	var value float64
	err := binary.Read(r.r, binary.LittleEndian, &value)
	if err != nil {
		log.Panicf("failed to read binary data: %e", err)
	}

	return value
}

// ReadTime reads time.Time object from the storage. It's saved in UTC, but converted into time.Local automatically
func (r *EventReader) ReadTime() time.Time {
	var millis int64
	err := binary.Read(r.r, binary.LittleEndian, &millis)
	if err != nil {
		log.Panicf("failed to read binary data: %e", err)
	}
	return time.UnixMilli(millis).In(time.Local)
}

func (r *EventReader) bytesLeft() int {
	return r.r.Len()
}
