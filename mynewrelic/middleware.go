package mynewrelic

import (
	"context"
	"net/http"

	newrelic "github.com/newrelic/go-agent"
)

type ctxID int

// Ctx context id to extract newrelic transaction
const Ctx ctxID = iota

// Middleware type for DI
type Middleware func(next http.Handler) http.Handler

// NewMiddleware create http middleware to inject newrelic transaction into request context
func NewMiddleware(app newrelic.Application) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			txn := app.StartTransaction(r.Method+" "+r.URL.Path, w, r)
			defer func() {
				_ = txn.End()
			}()
			ctx := context.WithValue(r.Context(), Ctx, txn)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
