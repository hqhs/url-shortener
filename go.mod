module github.com/hqhs/url-shortener

replace github.com/hqhs/url-shortener/service => ./service

replace github.com/hqhs/url-shortener/redis => ./redis

require (
	github.com/go-kit/kit v0.8.0
	github.com/hqhs/url-shortener/redis v0.0.0
	github.com/hqhs/url-shortener/service v0.0.0
)
