package testingutils

import (
	"context"
	"fmt"

	"github.com/sapphirenw/ai-content-creation-api/src/queries"
)

const TEST_CUSTOMER_ID = 7

func CreateTestCustomer(db queries.DBTX) (*queries.Customer, error) {
	model := queries.New(db)
	c, err := model.GetCustomer(context.TODO(), TEST_CUSTOMER_ID)
	if err != nil {
		return nil, fmt.Errorf("there was an issue getting the test customer: %v", err)
	}
	return c, err
}
