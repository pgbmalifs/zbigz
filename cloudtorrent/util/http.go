package util

import (
	"net/http"
	"strings"

	"golang.org/x/net/context"

	goji "goji.io"
)

type PrefixPattern struct {
	Prefix string
}

func (p *PrefixPattern) Match(ctx context.Context, r *http.Request) context.Context {
	if strings.HasPrefix(r.URL.Path, p.Prefix) {
		return ctx
	}
	return nil
}

func ErrHandler(fn func(ctx context.Context, w http.ResponseWriter, r *http.Request) error) goji.HandlerFunc {
	return goji.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		if err := fn(ctx, w, r); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})
}
