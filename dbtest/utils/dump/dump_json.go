package dump

import (
	"encoding/json"
	"fmt"
)

func ToJSON(obj any) (string, error) {
	b, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func PrintJson(obj any) error {
	j, err := ToJSON(obj)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", j)
	return nil
}
