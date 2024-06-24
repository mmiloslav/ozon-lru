module main

go 1.22

require (
	api v0.0.0-00010101000000-000000000000
	config v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.9.3
	lru v0.0.0-00010101000000-000000000000
)

require (
	common v0.0.0-00010101000000-000000000000 // indirect
	github.com/caarlos0/env/v8 v8.0.0 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
)

replace api => ./api

replace common => ./internal/common

replace config => ./internal/config

replace lru => ./pkg/lru
