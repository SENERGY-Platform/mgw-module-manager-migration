package naming

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

func SetManagerId(pth, val string) error {
	if val == "" {
		return errors.New("empty manager id")
	}
	file, err := os.OpenFile(pth, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("create manager id file: %w", err)
	}
	defer file.Close()
	je := json.NewEncoder(file)
	je.SetIndent("", "\t")
	err = je.Encode(struct {
		Id string `json:"manager_id"`
	}{Id: val})
	if err != nil {
		return fmt.Errorf("write manager id file: %w", err)
	}
	return nil
}
