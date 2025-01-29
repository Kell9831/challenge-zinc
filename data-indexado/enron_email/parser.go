package enron_email

import (
	"Kell9831/challenge-zinc/zinc"
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

var emailPool = sync.Pool{
	New: func() interface{} {
		return &zinc.Email{}
	},
}

var builderPool = sync.Pool{
	New: func() interface{} {
		return &strings.Builder{}
	},
}

func ParseEmail(filePath string) (*zinc.Email, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Obtener un objeto de zinc.Email del pool
	emailInterface := emailPool.Get()
	if emailInterface == nil {
		return nil, fmt.Errorf("error: emailPool.Get() devolvió nil")
	}
	email := emailInterface.(*zinc.Email)
	*email = zinc.Email{} // Resetear la estructura

	builder := builderPool.Get().(*strings.Builder)
	builder.Reset()

	scanner := bufio.NewScanner(file)
	parsingBody := false

	for scanner.Scan() {
		line := scanner.Text()
		if parsingBody {
			builder.WriteString(line)
			builder.WriteString("\n")
		} else {
			if line == "" {
				parsingBody = true
			} else {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					switch strings.ToLower(key) {
					case "from":
						email.From = value
					case "to":
						email.To = value
					case "subject":
						email.Subject = value
					}
				}
			}
		}
	}
	email.Body = builder.String()
	builderPool.Put(builder)
	return email, nil
}


func ReleaseEmail(email *zinc.Email) {
	if email != nil {
		emailPool.Put(email)
	}
}

