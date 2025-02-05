package main

import (
	"Kell9831/challenge-zinc/enron_email"
	"Kell9831/challenge-zinc/workers"
	"Kell9831/challenge-zinc/zinc"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
	"sync"
	"github.com/joho/godotenv"
	"time"
)

const (
	maildirPath = "./enron_mail_20110402/maildir"
	maxWorkers  = 10
)

func startWorkers(emailFiles chan string, batchChan chan []*zinc.Email, numWorkers int, wg *sync.WaitGroup) {
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go workers.Worker(emailFiles, wg, batchChan)
	}
}

func startBatchIndexers(batchChan chan []*zinc.Email, numIndexers int, wg *sync.WaitGroup) {
	for i := 0; i < numIndexers; i++ {
		wg.Add(1)
		go workers.BatchIndexer(batchChan, wg)
	}
}

func main() {
	
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}


	// Iniciar servidor de profiling
	go func() {
		fmt.Println("Iniciando servidor de profiling en :6060")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			fmt.Printf("Error iniciando servidor de profiling: %v\n", err)
		}
	}()

	// Crear perfiles de CPU y memoria
	cpuProfileFile, err := os.Create("cpu_profile.prof")
	if err != nil {
		fmt.Printf("Error creando archivo de perfil de CPU: %v\n", err)
		return
	}
	defer cpuProfileFile.Close()

	memProfileFile, err := os.Create("heap_profile.prof")
	if err != nil {
		fmt.Printf("Error creando archivo de perfil de memoria: %v\n", err)
		return
	}
	defer memProfileFile.Close()

	pprof.StartCPUProfile(cpuProfileFile)
	defer pprof.StopCPUProfile()

	startTime := time.Now()

	emailFiles := make(chan string, maxWorkers*2)
	batchChan := make(chan []*zinc.Email, maxWorkers)
	var wg sync.WaitGroup
	var wgBatchIndexer sync.WaitGroup

	startWorkers(emailFiles, batchChan, maxWorkers, &wg)
	startBatchIndexers(batchChan, 2, &wgBatchIndexer)

	// Recorrer el directorio
	go func() {
		err := enron_email.Walk(maildirPath, emailFiles)
		if err != nil {
			fmt.Printf("Error procesando maildir: %v\n", err)
		}
		close(emailFiles)
	}()

	go func() {
		wg.Wait()       // Espera a los workers
		close(batchChan) // Cierra el canal de batches
	}()

	wgBatchIndexer.Wait()

	pprof.WriteHeapProfile(memProfileFile)
	totalDuration := time.Since(startTime)
	fmt.Printf("Todos los correos han sido indexados en %s.\n", totalDuration)
}

