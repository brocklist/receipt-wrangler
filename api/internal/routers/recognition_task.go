package routers

import (
	"github.com/go-chi/chi/v5"
	"receipt-wrangler/api/internal/handlers"
	"receipt-wrangler/api/internal/middleware"
)

func BuildRecognitionTaskRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.UnifiedAuthMiddleware)
	router.Post("/", handlers.CreateRecognitionTask)
	router.Get("/", handlers.GetRecognitionTasks)
	router.Get("/{id}", handlers.GetRecognitionTask)
	router.Put("/{id}/file", handlers.UploadRecognitionTaskFile)
	router.Post("/{id}/retry", handlers.RetryRecognitionTask)
	return router
}
