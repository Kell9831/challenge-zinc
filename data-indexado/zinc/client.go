package zinc

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

const (
	zincURL      = "http://localhost:4080/api/default/_bulk"
	zincUser     = "admin"
	zincPassword = "Complexpass#123"
)

var client = &http.Client{}

func IndexEmails(emails []*Email) error {
	var payload bytes.Buffer
	for _, email := range emails {
		meta := `{"index": {"_index": "enron"}}`
		emailJSON, err := json.Marshal(email)
		if err != nil {
			return err
		}
		payload.WriteString(meta + "\n")
		payload.Write(emailJSON)
		payload.WriteString("\n")
	}

	req, err := http.NewRequest("POST", zincURL, &payload)
	if err != nil {
		return err
	}
	req.SetBasicAuth(zincUser, zincPassword)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("error en respuesta de ZincSearch: " + resp.Status)
	}
	return nil
}

