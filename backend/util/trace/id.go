package trace

import uuid "github.com/satori/go.uuid"

func NewTraceID() string {
	return uuid.NewV4().String()
}