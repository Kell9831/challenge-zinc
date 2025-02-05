package zinc

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

func IndexEmail(email *Email) error {

	zincURL := os.Getenv("ZINC_URL")
	zincUser := os.Getenv("ZINC_USER")
	zincPassword := os.Getenv("ZINC_PASSWORD")
	
	emailJSON, err := json.Marshal(email)
	if err != nil {
		return err
	}

	payload := bytes.NewBuffer([]byte(`{"index": {"_index": "enron"}}` + "\n" + string(emailJSON) + "\n"))

	//Crea una solicitud HTTP de tipo POST con el payload para enviar los datos a ZincSearch.
	req, err := http.NewRequest("POST", zincURL, payload)
	if err != nil {
		return err
	}

	req.SetBasicAuth(zincUser, zincPassword)
	req.Header.Set("Content-Type", "application/json")

	//Crea un cliente HTTP y envía la solicitud (client.Do(req)).
	client := &http.Client{}
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
