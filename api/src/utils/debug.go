package utils

import (
	"encoding/json"
	"fmt"
)

func DebugPrint(input any) {
	enc, err := json.MarshalIndent(input, "", "    ")
	if err != nil {
		return
	}
	fmt.Println(string(enc))
}
