module api

go 1.22

require lru v0.0.0-00010101000000-000000000000

require (
	common v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.1
	github.com/sirupsen/logrus v1.9.3
)

require golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect

replace lru => ../pkg/lru

replace common => ../internal/common
