package enron_email

import (
	"Kell9831/challenge-zinc/zinc"
	"bufio"
	"strings"
	"sync"
	"os"
)

var emailPool = sync.Pool{
	New: func() interface{} {
		return &zinc.Email{}
	},
}

func ParseEmail(filePath string) (*zinc.Email, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	email := emailPool.Get().(*zinc.Email)
	*email = zinc.Email{} // Resetear la estructura

	scanner := bufio.NewScanner(file)
	var body strings.Builder
	parsingBody := false

	for scanner.Scan() {
		line := scanner.Text()
		if parsingBody {
			body.WriteString(line)
			body.WriteString("\n")
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
	email.Body = body.String()
	return email, nil
}

func ReleaseEmail(email *zinc.Email) {
	emailPool.Put(email)
}
