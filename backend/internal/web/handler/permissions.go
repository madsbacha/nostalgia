package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"nostalgia/internal/app/command"
	"nostalgia/internal/app/query"
	"nostalgia/internal/common/rbac"
	"nostalgia/internal/common/util"
	"nostalgia/internal/web/middleware"
	"nostalgia/internal/web/models"
)

type GetPermissionsResponse = models.UserPermissions

func (s HttpServer) GetPermissions(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userID")
	roles, err := s.app.Queries.GetRolesForUser.Handle(r.Context(), query.GetRolesForUser{
		UserId: userId,
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	permissions := models.NewUserPermissionsFromRoles(roles)

	// TODO: Can this be made more clean?
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(permissions)
}

type SetPermissionsRequest = models.UserPermissions

func (s HttpServer) SetPermissions(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userID")
	var req SetPermissionsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	addRoles, removeRoles := req.RolesFromUserPermissions()

	currentUserId := middleware.UserIdFromContext(r.Context())
	if userId == currentUserId {
		if util.Contains(removeRoles, rbac.RoleWhitelisted) || util.Contains(removeRoles, rbac.RoleCanManagePermissions) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
	}

	err = s.app.Commands.AddRolesToUser.Handle(r.Context(), command.AddRolesToUser{
		UserId: userId,
		Roles:  addRoles,
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = s.app.Commands.RemoveRolesFromUser.Handle(r.Context(), command.RemoveRolesFromUser{
		UserId: userId,
		Roles:  removeRoles,
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}
