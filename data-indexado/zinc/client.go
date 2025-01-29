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

func IndexEmail(email *Email) error {

	emailJSON, err := json.Marshal(email)
	if err != nil {
		return err
	}

	payload := bytes.NewBuffer([]byte(`{"index": {"_index": "enron"}}` + "\n" + string(emailJSON) + "\n"))

	req, err := http.NewRequest("POST", zincURL, payload)
	if err != nil {
		return err
	}

	req.SetBasicAuth(zincUser, zincPassword)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	//asegura que el cuerpo de la respuesta se cierre después de leerla, evitando fugas de recursos.
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("error en respuesta de ZincSearch: " + resp.Status)
	}
	return nil
}
