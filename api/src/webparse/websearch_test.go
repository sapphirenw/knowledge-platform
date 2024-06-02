package webparse

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWeb(t *testing.T) {
	query := "what is the most common rock in the world"

	response, err := WebSearch(query)
	require.NoError(t, err)
	require.NotEmpty(t, response.Results)

	enc, _ := json.MarshalIndent(response.Results[0], "", "    ")
	fmt.Println(string(enc))
}

func TestImage(t *testing.T) {
	query := "!images what is the most common rock in the world"

	response, err := WebSearch(query)
	require.NoError(t, err)
	require.NotEmpty(t, response.Results)
	require.NotEmpty(t, response.Results[0].ImgSrc)

	enc, _ := json.MarshalIndent(response.Results[0], "", "    ")
	fmt.Println(string(enc))
}
