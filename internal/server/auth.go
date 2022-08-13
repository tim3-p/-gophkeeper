package server

import (
	"log"
	"net/http"
)

func verifyUser(user string, pass string) (bool, error) {
	ok, err := serverStore.CheckUserAuth(user, pass)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func authUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, basicSet := r.BasicAuth()

		if basicSet {
			userOK, err := verifyUser(user, pass)
			if err != nil {
				log.Printf("verifyUser error: %v", err)
				writeStatus(w,
					http.StatusForbidden,
					"Access Denied",
				)
				return
			}
			if !userOK {
				log.Printf("verifyUser: password incorrect")
				writeStatus(w,
					http.StatusForbidden,
					"Access Denied",
				)
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		if !basicSet {
			if r.URL.Path == registerPath ||
				r.URL.Path == registerPath+"/" {
				// registering does not require user auth
				next.ServeHTTP(w, r)
				return
			}
			w.Header().Set("WWW-Authenticate", `Basic realm="storeapi"`)
			writeStatus(w,
				http.StatusUnauthorized,
				"Unauthorized",
			)
			return
		}

	})
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("ping")
	user, _, ok := r.BasicAuth()
	if !ok {
		writeStatus(w, http.StatusBadRequest, "no basic auth")
		return
	}
	writeStatus(w, http.StatusOK, "OK. User "+user)
}
