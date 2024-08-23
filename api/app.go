package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	oapimiddleware "github.com/oapi-codegen/nethttp-middleware"
)

func Serve(port int) {
	http.ListenAndServe(fmt.Sprintf(":%d", port), GetRouter())
}

func GetRouter() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// Health check
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Swagger UI
	router.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/swagger-ui/index.html")
	})
	router.Get("/docs/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./api.yml")
	})

	// API Routes
	router.Mount("/receipts", ReceiptRoutes())
	return router
}

func ReceiptRoutes() chi.Router {
	router := chi.NewRouter()
	// request validator
	spec, _ := GetSwagger()
	router.Use(oapimiddleware.OapiRequestValidator(spec))
	handler := NewReceiptHandler()
	router.Post("/process", handler.PostReceiptsProcess)
	router.Get("/{id}/points", handler.GetReceiptsIdPoints)
	return router
}
