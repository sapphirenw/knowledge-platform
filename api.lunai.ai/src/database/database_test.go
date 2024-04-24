package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseConnection(t *testing.T) {
	p, err := GetPool(nil)
	assert.Nil(t, err)
	if err != nil {
		return
	}
	assert.NotNil(t, p)
	if p == nil {
		return
	}
	err = p.Ping(context.TODO())
	assert.Nil(t, err)
}
