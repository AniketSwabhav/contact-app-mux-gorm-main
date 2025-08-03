package util

import (
	"contact_app_mux_gorm_main/components/apperror"
	"contact_app_mux_gorm_main/components/security/middleware/authorization"
	"net/http"
)

func MiddlewareAdmin(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claim, err := authorization.ValidateToken(w, r)
		if err != nil {
			RespondError(w, apperror.NewInValidTokenError("Invalid or missing token"))
			return
		}

		if !claim.IsAdmin {
			RespondError(w, apperror.NewUnauthorizedUserError("Current user not an admin"))
			return
		}

		if !claim.IsActive {
			RespondError(w, apperror.NewInActiveUserError("current user is not active"))
			return
		}

		next.ServeHTTP(w, r)

	})
}

func MiddlewareContact(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claim, err := authorization.ValidateToken(w, r)
		if err != nil {
			RespondError(w, apperror.NewInValidTokenError("Invalid or missing token"))
			return
		}
		if claim.IsAdmin {
			RespondError(w, apperror.NewUnauthorizedUserError("admin can't access user details"))
			return
		}
		if !claim.IsActive {
			RespondError(w, apperror.NewInActiveUserError("current user is not active"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
