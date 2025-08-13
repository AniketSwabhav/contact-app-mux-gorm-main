package util

import (
	"contact_app_mux_gorm_main/components/apperror"
	"contact_app_mux_gorm_main/components/security/middleware/authorization"
	"fmt"
	"net/http"
)

func MiddlewareAdmin(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		claim := authorization.Claims{}

		err := authorization.ValidateToken(w, r, &claim)
		if err != nil {
			RespondError(w, apperror.NewInValidTokenError("Invalid or missing token"))
			return
		}

		if !claim.IsAdmin {
			fmt.Println("User is not Admin")
			RespondError(w, apperror.NewUnauthorizedUserError("Current user not an admin"))
			return
		}

		if !claim.IsActive {
			fmt.Println("User is not Active")
			RespondError(w, apperror.NewInActiveUserError("current user is not active"))
			return
		}

		next.ServeHTTP(w, r)

	})
}

func MiddlewareContact(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		claim := authorization.Claims{}

		err := authorization.ValidateToken(w, r, &claim)
		if err != nil {
			RespondError(w, apperror.NewInValidTokenError("Invalid or missing token"))
			return
		}
		if claim.IsAdmin {
			fmt.Println("User is Admin")
			RespondError(w, apperror.NewUnauthorizedUserError("admin can't access user details"))
			return
		}
		if !claim.IsActive {
			fmt.Println("User is not Active")
			RespondError(w, apperror.NewInActiveUserError("current user is not active"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
