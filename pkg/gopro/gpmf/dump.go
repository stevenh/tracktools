package gpmf

import (
	"encoding/json"
	"fmt"
	"os"
)

func dump(data []*Element) error {
	d := struct {
		Data []*Element
	}{Data: data}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(d); err != nil {
		return fmt.Errorf("json encode: %w", err)
	}

	fmt.Println("")

	return nil
}
