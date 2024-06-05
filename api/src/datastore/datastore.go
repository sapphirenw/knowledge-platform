package datastore

import (
	"bytes"
	"context"
)

type Object interface {
	GetRaw(ctx context.Context) (*bytes.Buffer, error)
	GetCleaned(ctx context.Context) (*bytes.Buffer, error)
	GetSha256() (string, error)
}
