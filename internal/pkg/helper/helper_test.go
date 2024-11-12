package helper

import (
	"fmt"
	"testing"
)

func TestStringArrToJSON(t *testing.T) {
	js, err := StringArrToJSONStr([]string{"apple", "banana", "cherry"})
	if err != nil {
		return
	}

	fmt.Println("JSON:", js)
}
