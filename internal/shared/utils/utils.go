package utils

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

func GenUUID() string {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id.String()
}
