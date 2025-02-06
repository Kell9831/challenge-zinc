package workers

import (
	"Kell9831/challenge-zinc/enron_email"
	"Kell9831/challenge-zinc/zinc"
	"fmt"
	"sync"
)

// procesa los correos ejecutandose concurrente
func Worker(emailFiles chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for filePath := range emailFiles {

		email, err := enron_email.ParseEmail(filePath)
		if err != nil {
			fmt.Printf("Error procesando archivo %s: %v\n", filePath, err)
			continue
		}

		err = zinc.IndexEmail(email)
		if err != nil {
			fmt.Printf("Error indexando correo: %v\n", err)
		} else {
			fmt.Printf("Correo indexado: %s\n", filePath)
		}
	}
}
