package http

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"stocks-api/module/controllers"
)

// Serve a server instance.
type Serve struct {
	Server          *mux.Router
	logger          *logrus.Logger
	stockController *controllers.StockController
}

// NewServe a constructor for Serve.
func NewServe(stockController *controllers.StockController, l *logrus.Logger) *Serve {
	return &Serve{
		Server:          mux.NewRouter(),
		logger:          l,
		stockController: stockController,
	}
}

// RegisterHandlers registers the routes available for our API.
func (s *Serve) RegisterHandlers() {
	s.Server.HandleFunc("/", s.stockController.GetAll).Methods("GET")
	s.Server.HandleFunc("/", s.stockController.InsertOne).Methods("POST")
	s.Server.HandleFunc("/{id}", s.stockController.GetOne).Methods("GET")
	s.Server.HandleFunc("/{id}", s.stockController.UpdateOne).Methods("POST")
	s.Server.HandleFunc("/{id}", s.stockController.DeleteOne).Methods("PUT")
}

// Serve starts the server and accepts new calls.
func (s *Serve) Serve() {
	port := os.Getenv("SERVER_PORT")
	address := os.Getenv("SERVER_ADDRESS")
	if port == "" || address == "" {
		s.logger.Fatal(errors.New("failed to load env param"))
	}

	s.RegisterHandlers()
	s.Server.Use(s.loggingMiddleware)

	s.logger.Info(fmt.Sprintf("Serving HTTP on: %s:%s", address, port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), s.Server))
}

func (s *Serve) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info(r.URL, r.Body)
		next.ServeHTTP(w, r)
	})
}
