package main

import (
	"Kell9831/challenge-zinc/enron_email"
	"Kell9831/challenge-zinc/workers"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
	"sync"
	"time"
)

const (
	maildirPath = "./enron_mail_20110402/maildir"
	maxWorkers  = 10
)

func main() {

	// Iniciar servidor de profiling en una goroutine separada
	go func() {
		fmt.Println("Iniciando servidor de profiling en :6060")
		http.ListenAndServe("localhost:6060", nil)
	}()

	    // Crear archivos para perfiles
		cpuProfileFile, err := os.Create("cpu_profile.prof")
		if err != nil {
			fmt.Println("Error al crear archivo de perfil de CPU:", err)
			return
		}
		defer cpuProfileFile.Close()
	
		memProfileFile, err := os.Create("heap_profile.prof")
		if err != nil {
			fmt.Println("Error al crear archivo de perfil de memoria:", err)
			return
		}
		defer memProfileFile.Close()
	
		// Iniciar profiling de CPU
		fmt.Println("Iniciando perfil de CPU")
		pprof.StartCPUProfile(cpuProfileFile)

		  // Lógica de tu programa
		emailFiles := make(chan string, maxWorkers)
		var wg sync.WaitGroup

		for i := 0; i < maxWorkers; i++ {
			wg.Add(1)
			go workers.Worker(emailFiles, &wg)
		}

		err = enron_email.Walk(maildirPath, emailFiles)
		if err != nil {
		fmt.Printf("Error procesando maildir: %v\n", err)
		}

		close(emailFiles)

		wg.Wait()
		pprof.StopCPUProfile()

	    // Capturar perfil de memoria
		fmt.Println("Capturando perfil de memoria")
		pprof.WriteHeapProfile(memProfileFile)
	
		fmt.Println("Todos los correos han sido indexados.")
		fmt.Println("Programa finalizado. Manteniendo servidor activo...")
		time.Sleep(30 * time.Minute) // Mantener activo para analizar desde /pprof

}

