package enron_email

import (
	"Kell9831/challenge-zinc/zinc"
	"bufio"
	"os"
	"strings"
)

func ParseEmail(filePath string) (*zinc.Email, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	email := &zinc.Email{}
	var body []string
	parsingBody := false

	for scanner.Scan() {
		line := scanner.Text()
		if parsingBody {
			body = append(body, line)
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
	email.Body = strings.Join(body, "\n")
	return email, nil
}
