package server

import (
	"github.com/gorilla/mux"
	"stocks-api/module/controllers"
)

type Serve struct {
	Server          *mux.Router
	stockController *controllers.StockController
}

func NewServe(stockController *controllers.StockController) *Serve {
	return &Serve{
		Server:          mux.NewRouter().StrictSlash(true),
		stockController: stockController,
	}
}

func (s *Serve) RegisterHandlers() {
	s.Server.HandleFunc("/", s.stockController.GetAll).Methods("GET")
	s.Server.HandleFunc("/{id}", s.stockController.GetOne).Methods("GET")
	s.Server.HandleFunc("/{id}", s.stockController.GetOne).Methods("POST")
	s.Server.HandleFunc("/{id}", s.stockController.DeleteOne).Methods("PUT")
}
