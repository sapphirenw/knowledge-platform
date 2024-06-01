package websearch

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWeb(t *testing.T) {
	query := "what is the most common rock in the world"

	response, err := Web(query)
	require.NoError(t, err)
	require.NotEmpty(t, response.Results)

	fmt.Println(*response.Results[0])
}
