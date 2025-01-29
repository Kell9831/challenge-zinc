package workers

import (
	"Kell9831/challenge-zinc/enron_email"
	"Kell9831/challenge-zinc/zinc"
	"fmt"
	"sync"
	"sync/atomic"
)

const batchSize = 100
var indexedCount uint64

func Worker(emailFiles chan string, wg *sync.WaitGroup, batchChan chan []*zinc.Email) {
	defer wg.Done()
	batch := make([]*zinc.Email, 0, batchSize)

	for filePath := range emailFiles {
		email, err := enron_email.ParseEmail(filePath)
		if err != nil {
			fmt.Printf("Error procesando archivo %s: %v\n", filePath, err)
			continue
		}
		batch = append(batch, email)
		if len(batch) == batchSize {
			batchChan <- batch
			for i := range batch {
				batch[i] = nil // Evita referencias pendientes
			}
			batch = batch[:0]
		}
	}

	// Enviar cualquier correo restante en el último batch
	if len(batch) > 0 {
		batchChan <- batch
	}
}

func BatchIndexer(batchChan chan []*zinc.Email, wg *sync.WaitGroup) {
	defer wg.Done()
	for batch := range batchChan {
		err := zinc.IndexEmails(batch)
		if err != nil {
			fmt.Printf("Error indexando batch: %v\n", err)
			continue
		}

		count := atomic.AddUint64(&indexedCount, uint64(len(batch)))
		if count%100 == 0 {
			fmt.Printf("Correos indexados: %d\n", count)
		}

		for i := range batch {
			enron_email.ReleaseEmail(batch[i])
			batch[i] = nil
		}
	}
}

