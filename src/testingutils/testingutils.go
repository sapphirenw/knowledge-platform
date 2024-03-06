package testingutils

import (
	"context"

	"github.com/sapphirenw/ai-content-creation-api/src/queries"
)

func CreateTestCustomer(db queries.DBTX) (*queries.Customer, error) {
	model := queries.New(db)
	item, err := model.GetCustomer(context.TODO(), 14)
	if err != nil {
		return nil, err
	}
	return &item, err
}
