package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Sabbir185/gpi/pkg/httpx"
	"github.com/Sabbir185/gpi/pkg/token"
)

const authKey ctxKey = iota + 1

type Auth struct {
	UserID    int64
	SessionID string
}

// SessionValidator rejects revoked or expired sessions and inactive users, so
// signing out or blocking a user takes effect before the access token expires.
type SessionValidator interface {
	ValidateSession(ctx context.Context, sessionID string, userID int64) (bool, error)
}

type Authenticator func(http.HandlerFunc) http.HandlerFunc

func Authenticate(tokens *token.Manager, sessions SessionValidator) Authenticator {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			raw, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
			if !ok || raw == "" {
				sendUnauthorized(w)
				return
			}

			claims, err := tokens.ParseAccessToken(raw)
			if err != nil {
				sendUnauthorized(w)
				return
			}
			userID, err := claims.UserID()
			if err != nil {
				sendUnauthorized(w)
				return
			}

			valid, err := sessions.ValidateSession(r.Context(), claims.SessionID, userID)
			if err != nil {
				httpx.SendAppError(w, r, err)
				return
			}
			if !valid {
				sendUnauthorized(w)
				return
			}

			ctx := context.WithValue(r.Context(), authKey, Auth{UserID: userID, SessionID: claims.SessionID})
			next(w, r.WithContext(ctx))
		}
	}
}

func AuthFromContext(ctx context.Context) Auth {
	auth, _ := ctx.Value(authKey).(Auth)
	return auth
}

func sendUnauthorized(w http.ResponseWriter) {
	httpx.SendError(w, http.StatusUnauthorized, httpx.CodeUnauthorized, "Unauthorized", nil)
}
