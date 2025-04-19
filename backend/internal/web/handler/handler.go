package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"nostalgia/internal/app"
	"nostalgia/internal/common/rbac"
	"nostalgia/internal/web/middleware"
)

type HttpServer struct {
	app       app.Application
	logger    *logrus.Entry
	tokenAuth *jwtauth.JWTAuth
	rbac      rbac.Context
}

func NewHttpServer(application app.Application, logger *logrus.Entry, tokenAuth *jwtauth.JWTAuth) HttpServer {
	return HttpServer{
		app:       application,
		logger:    logger,
		tokenAuth: tokenAuth,
		rbac:      rbac.New(application),
	}
}

func (s HttpServer) RegisterServerHandlers(r chi.Router) http.Handler {
	rbacMiddleware := middleware.NewRbacMiddleware(s.rbac)
	guard := rbacMiddleware.Guard
	r.Group(func(r chi.Router) {
		r.Use(middleware.Verifier(s.tokenAuth))
		r.Use(jwtauth.Authenticator(s.tokenAuth))
		r.Use(middleware.SetUserIdFromJwtMiddleware)
		r.Use(middleware.SetFileIdFromJwt)
		r.Group(func(r chi.Router) {
			r.Use(guard(rbac.All(rbac.RoleWhitelisted)))

			r.Get("/users/@me", s.GetCurrentUser)
			r.With(guard(rbac.All(rbac.RoleCanUploadMedia))).Post("/media", s.UploadMedia)
			//r.Put("/media/{mediaID}", s.UpdateMedia)
			//r.Delete("/media/{mediaID}", s.DeleteMedia)
			//r.Post("/media/{mediaID}/tags", s.AddTagToMedia)
			r.Group(func(r chi.Router) {
				r.Use(guard(rbac.All(rbac.RoleCanViewMedia)))
				r.Get("/media/{mediaID}", s.GetMediaById)
				r.Get("/media", s.GetMedia)
				r.Get("/media/tags", s.GetAllTags)
				r.Get("/media/{mediaID}/tags", s.GetTagsForMedia)
			})
			r.Group(func(r chi.Router) {
				r.Use(guard(rbac.All(rbac.RoleCanManagePermissions)))
				r.Get("/users", s.GetUsers)
				r.Get("/users/{userID}/permissions", s.GetPermissions)
				r.Post("/users/{userID}/permissions", s.SetPermissions)
			})
		})

		r.Get("/files/{fileID}", s.GetFileById)
	})

	r.Post("/auth/discord", s.AuthDiscord)
	r.Get("/auth/discord", s.GetDiscordInfo)
	r.Get("/info", s.GetInfo)

	return r
}
