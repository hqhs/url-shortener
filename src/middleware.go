package main

import (
	"net/http"
	"context"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/gomodule/redigo/redis"
)

// ContextKey represent key used for values in r.Context()
type ContextKey int

const (
	// DatabaseKey used for db connection
	DatabaseKey ContextKey = iota
)

// GetDatabaseMiddleware returns middleware which add redis conn to every
// request's context
func GetDatabaseMiddleware(p *redis.Pool) (func(http.Handler) http.Handler) {
	// NOTE: Generally speaking, this one is arguable, but better then
	// global pool variable and easier then struct embedding
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			conn := p.Get()
			ctx := context.WithValue(r.Context(), DatabaseKey, conn)
			next.ServeHTTP(w, r.WithContext(ctx))
			conn.Close()
		})
	}
}

// Paginate is a stub, but very possible to implement middleware logic
// to handle the request params for handling a paginated request.
func Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// just a stub.. some ideas are to look at URL query params for something like
		// the page number, or the limit, and send a query cursor down the chain
		next.ServeHTTP(w, r)
	})
}

// AdminOnly middleware restricts access to just administrators.
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin, ok := r.Context().Value("acl.admin").(bool)
		if !ok || !isAdmin {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ArticleCtx middleware is used to load an Article object from
// the URL parameters passed through as the request. In case
// the Article could not be found, we stop here and return a 404.
func ArticleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var article *Article
		var err error

		if articleID := chi.URLParam(r, "articleID"); articleID != "" {
			article, err = dbGetArticle(articleID)
		} else if articleSlug := chi.URLParam(r, "articleSlug"); articleSlug != "" {
			article, err = dbGetArticleBySlug(articleSlug)
		} else {
			render.Render(w, r, ErrNotFound)
			return
		}
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "article", article)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
