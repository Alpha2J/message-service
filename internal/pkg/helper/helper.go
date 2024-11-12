package helper

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func FindProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find project root")
		}
		dir = parent
	}
}

func StringArrToJSONStr(arr []string) (string, error) {
	jsonBytes, err := json.Marshal(arr)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func JSONStrToStringArr(jsonStr string) ([]string, error) {
	if jsonStr == "" {
		return []string{}, nil
	}

	var stringArr []string
	err := json.Unmarshal([]byte(jsonStr), &stringArr)
	if err != nil {
		log.Fatal(err)
	}

	return stringArr, nil
}

func NewUUIDStr() string {
	uuidStr := uuid.New().String()
	return strings.ReplaceAll(uuidStr, "-", "")
}
