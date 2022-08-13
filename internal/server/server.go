package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/tim3-p/gophkeeper/internal/store"
)

const (
	defaultListenAddress = ":8080"
	defaultStoreFile     = "server_store.db"
	// registerPath is the path to serve requests to register new users
	registerPath = "/users"
)

func writeStatus(w http.ResponseWriter, code int, status string) {
	w.WriteHeader(code)
	w.Write([]byte(`{"Status":"` + status + `"}`))
}

func checkSetContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		cType := r.Header.Get("Content-Type")
		if cType != "application/json" {
			log.Print("checkContentType: bad content type " + cType)
			writeStatus(w, http.StatusBadRequest, "Bad Content Type")
			return
		}
		next.ServeHTTP(w, r)
	})
}

var serverStore *store.Store

// InitStore initialiese the server store
func InitStore(storeFile string) error {
	if storeFile == "" {
		storeFile = defaultStoreFile
	}
	var err error
	serverStore, err = store.NewStore(storeFile)
	if err != nil {
		return err
	}
	return nil
}

// DropServerStore drops server storage
func DropServerStore(storeFile string) error {
	if storeFile == "" {
		storeFile = defaultStoreFile
	}
	err := store.DropStore(storeFile)
	if err != nil {
		return err
	}
	return nil
}

// StartServer starts the server
func StartServer(listenPort int, storeFile, keyFile, crtFile string) error {
	err := InitStore(storeFile)
	if err != nil {
		return err
	}

	listenAddress := fmt.Sprintf(":%d", listenPort)
	if listenPort == 0 {
		listenAddress = defaultListenAddress
	}

	r := NewRouter()
	c := make(chan error)
	go func() {
		log.Printf("Listening on %v...", listenAddress)
		err := http.ListenAndServeTLS(listenAddress, crtFile, keyFile, r)
		c <- err
	}()

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)

	select {
	case sig := <-signalChannel:
		switch sig {
		case os.Interrupt:
			log.Print("sigint")
		case syscall.SIGTERM:
			log.Print("sigterm")
		case syscall.SIGINT:
			log.Print("sigint")
		case syscall.SIGQUIT:
			log.Print("sigquit")
		}
	case err := <-c:
		log.Print(err)
		return err
	}

	log.Print("Server finished")
	err = serverStore.CloseDB()
	if err != nil {
		return err
	}
	return nil
}

// NewRouter returns new Router
func NewRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(checkSetContentType)
	r.Use(authUser)

	r.Post("/users", createUser)
	r.Put("/password", changePassword)
	r.Get("/ping", pingHandler)
	r.Post("/records", storeRecord)
	r.Get("/records", listRecords)
	r.Get("/records/by_type/{record_type}", listRecordsByType)
	r.Get("/records/{id}", getRecordByID)
	r.Get("/records/{record_type}/{record_name}", getRecordID)
	r.Put("/records/{id}", updateRecordByID)
	r.Delete("/records/{id}", deleteRecordByID)

	return r
}
