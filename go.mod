module github.com/hqhs/url-shortener

replace example.com/url-shortener/service => ./service

replace example.com/url-shortener/redis => ./redis

require example.com/url-shortener/service v0.0.0
