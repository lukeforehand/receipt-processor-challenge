package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Serve(port int) {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	// Health check
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	// API Routes
	router.Mount("/receipts", ReceiptRoutes())
	http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}

func ReceiptRoutes() chi.Router {
	router := chi.NewRouter()
	handler := NewReceiptHandler()
	router.Post("/process", handler.PostReceiptsProcess)
	router.Get("/{id}/points", handler.GetReceiptsIdPoints)
	return router
}
