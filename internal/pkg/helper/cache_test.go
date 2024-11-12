package helper

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	cache := NewCache()
	cache.Set("username", "johndoe")

	if val, found := cache.Get("username"); found {
		fmt.Println("Found:", val)
	} else {
		fmt.Println("Not found")
	}
}

func TestSet(t *testing.T) {

}
