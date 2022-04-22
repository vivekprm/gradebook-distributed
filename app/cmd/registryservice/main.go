package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"app/registry"
)

func main () {
	http.Handle("/services", &registry.RegistryService{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var srv http.Server
	srv.Addr = registry.ServerPort

	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()

	go func ()  {
		fmt.Println("Registry Service started. Press any key to stop")
		var s string
		fmt.Scanln(&s)	// Blocked here until we receive a keystroke
		srv.Shutdown(ctx)
		cancel()
	}()

	<-ctx.Done()
	fmt.Println("Shutting down registry service.")
}