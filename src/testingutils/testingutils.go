package testingutils

import (
	"context"
	"fmt"

	"github.com/sapphirenw/ai-content-creation-api/src/queries"
)

func CreateTestCustomer(db queries.DBTX) (*queries.Customer, error) {
	model := queries.New(db)
	item, err := model.CreateCustomer(context.TODO(), "TEST CUSTOMER")
	if err != nil {
		return nil, fmt.Errorf("there was an issue creating the customer: %v", err)
	}
	// create the root folder
	_, err = model.CreateFolderRoot(context.TODO(), item.ID)
	if err != nil {
		return nil, fmt.Errorf("there was an issue creating the root folder: %v", err)
	}
	return item, err
}
