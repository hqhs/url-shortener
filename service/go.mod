module example.com/url-shortener/service

require (
	example.com/url-shortener/redis v0.0.0
	github.com/go-chi/chi v4.0.1+incompatible
	github.com/go-chi/docgen v1.0.5
	github.com/go-chi/render v1.0.1
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/rafaeljusto/redigomock v0.0.0-20190202135759-257e089e14a1
)

replace example.com/url-shortener/redis => ../redis