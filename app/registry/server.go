package registry

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

const (
	ServerPort = ":3000"
	ServicesURL = "http://localhost:" + ServerPort + "/services"
)

type registry struct {
	registrations []Registration
	// To make sure we are changing these registration objects in thread-safe way.
	mutex *sync.Mutex
}

func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	// We are not making deferred call as we are going to make other calls 
	// from here, which will use the mutex.
	r.mutex.Unlock()
	return nil
}

var reg = registry{
	registrations: make([]Registration, 0),
	mutex: new(sync.Mutex),
}

type RegistryService struct {}

func (s RegistryService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println("Request received.")
	switch req.Method {
	case http.MethodPost:
		dec := json.NewDecoder(req.Body)
		var r Registration
		err := dec.Decode(&r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Adding service: %v, with URL: %v", r.ServiceName, r.ServiceURL)
		err = reg.add(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}