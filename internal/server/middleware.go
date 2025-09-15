package server

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ninox14/gore-codenames/internal/response"
	"github.com/pascaldekloe/jwt"

	"github.com/tomasen/realip"
)

type Middleware func(http.Handler) http.Handler

func (s *Server) CreateMWStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}

		return next
	}
}

func (s *Server) corsMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// FIXME: change cors origin on deploy
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

func (s *Server) recoverPanicMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			pv := recover()
			if pv != nil {
				s.serverError(w, r, fmt.Errorf("%v", pv))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (s *Server) logAccessMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mw := response.NewMetricsResponseWriter(w)
		next.ServeHTTP(mw, r)

		var (
			ip     = realip.FromRequest(r)
			method = r.Method
			url    = r.URL.String()
			proto  = r.Proto
		)

		userAttrs := slog.Group("user", "ip", ip)
		requestAttrs := slog.Group("request", "method", method, "url", url, "proto", proto)
		responseAttrs := slog.Group("response", "status", mw.StatusCode, "size", mw.BytesCount)

		s.logger.Info("access", userAttrs, requestAttrs, responseAttrs)
	})
}

func (s *Server) validateTokenClaims(w http.ResponseWriter, r *http.Request, claims *jwt.Claims) error {
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		s.serverError(w, r, err)
		return err
	}

	if !claims.Valid(time.Now()) {
		s.invalidAuthenticationTokenWithUserId(w, r, userID)
		return errors.New("expired token")
	}

	if claims.Issuer != s.config.baseURL {
		s.invalidAuthenticationTokenWithUserId(w, r, userID)
		return errors.New("invalid issuer")
	}

	if !claims.AcceptAudience(s.config.baseURL) {
		s.invalidAuthenticationTokenWithUserId(w, r, userID)
		return errors.New("unacceptable audience")
	}
	return nil
}

func (s *Server) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scheme := r.URL.Scheme
		if scheme == "ws" || scheme == "wss" {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			next.ServeHTTP(w, r)
			return
		}
		headerParts := strings.Split(authorizationHeader, " ")

		if len(headerParts) == 2 && headerParts[0] == "Bearer" {
			token := headerParts[1]

			claims, err := jwt.HMACCheck([]byte(token), []byte(s.config.jwt.secretKey))

			if err != nil {
				s.invalidAuthenticationToken(w, r)
				return
			}
			err = s.validateTokenClaims(w, r, claims)
			if err != nil {
				// Error response gets handled by validateTokenClaims
				return
			}

			userID, _ := uuid.Parse(claims.Subject)

			user, err := s.db.Queries.GetUserByID(r.Context(), userID)
			if err != nil {
				s.serverError(w, r, err)
				return
			}
			r = contextSetAuthenticatedUser(r, user)
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, found := contextGetAuthenticatedUser(r)

		if !found {
			s.authenticationRequired(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
