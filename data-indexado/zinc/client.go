package zinc

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
)

const (
	zincURL      = "http://localhost:4080/api/default/_bulk"
	zincUser     = "admin"
	zincPassword = "Complexpass#123"
)

var client = &http.Client{}

var bufferPool = sync.Pool{
	New: func ()  interface{}	{
		return new(bytes.Buffer)
	},
}

func IndexEmails(emails []*Email) error {
	buffer := bufferPool.Get().(*bytes.Buffer)
	buffer.Reset()

	for _, email := range emails {
		meta := `{"index": {"_index": "enron"}}`
		buffer.WriteString(meta + "\n")
		encoder := json.NewEncoder(buffer)
		err := encoder.Encode(email)
		if err != nil {
			bufferPool.Put(buffer)
			return err
		}
	}

	req, err := http.NewRequest("POST", zincURL, buffer)
	if err != nil {
		bufferPool.Put(buffer)
		return err
	}

	req.SetBasicAuth(zincUser, zincPassword)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		bufferPool.Put(buffer)
		return err
	}
	defer resp.Body.Close()

	bufferPool.Put(buffer) // Reutilizar buffer

	if resp.StatusCode != http.StatusOK {
		return errors.New("error en respuesta de ZincSearch: " + resp.Status)
	}
	return nil
}

