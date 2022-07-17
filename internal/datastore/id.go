package datastore

import (
	"github.com/google/uuid"

	"github.com/lithammer/shortuuid"
)

func newKey(length int) string {
	return shortuuid.New()[:length]
}

func newUUID() string {
	return uuid.NewString()
}
