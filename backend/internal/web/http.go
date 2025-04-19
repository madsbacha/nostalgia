package web

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

func addCorsMiddleware(router *chi.Mux) {
	allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ";")
	if len(allowedOrigins) == 0 {
		return
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	router.Use(corsMiddleware.Handler)
}

func RunHttpServer(createApiHandler func(router chi.Router) http.Handler) {
	rootRouter := chi.NewRouter()

	rootRouter.Use(middleware.Logger)
	rootRouter.Use(middleware.Recoverer)
	rootRouter.Use(middleware.CleanPath)
	rootRouter.Use(middleware.Heartbeat("/ping"))

	apiRouter := chi.NewRouter()
	addCorsMiddleware(apiRouter)
	// we are mounting all APIs under /api path
	rootRouter.Mount("/api", createApiHandler(apiRouter))

	if strings.TrimSpace(os.Getenv("PORT")) == "" {
		logrus.Panicln("PORT environment variable not set")
	}
	addr := ":" + os.Getenv("PORT")

	logrus.WithField("address", addr).Info("Starting HTTP server")

	err := http.ListenAndServe(addr, rootRouter)
	if err != nil {
		logrus.WithError(err).Panic("Unable to start HTTP server")
	}
}
