package datastore

import (
	"fmt"
	"testing"
)

func TestDatastore(t *testing.T) {
	// ensure all adhere to interface
	objects := []Object{
		&Document{},
		&WebsitePage{},
	}

	fmt.Println(len(objects))
}
