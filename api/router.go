// пакет с api
package api

import (
	"net/http"
	"reflect"
	"runtime"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// структура путей хендлеров
type route struct {
	Name               string
	Method             string
	Pattern            string
	HandlerFunc        http.HandlerFunc
	MiddlewareAuthFunc func(http.Handler) http.Handler
}

// Инициализация маршрутов и создание роутера
func NewRouter(ch *CacheHandler) *mux.Router {
	var routes = []route{
		{Name: "Ping", Method: http.MethodGet, Pattern: "/api/ping", HandlerFunc: pingHandler, MiddlewareAuthFunc: logMiddleware},
		{Name: "Put", Method: http.MethodPost, Pattern: "/api/lru", HandlerFunc: ch.putHandler, MiddlewareAuthFunc: logMiddleware},
		{Name: "Get", Method: http.MethodGet, Pattern: "/api/lru/{key}", HandlerFunc: ch.getHandler, MiddlewareAuthFunc: logMiddleware},
		{Name: "GetAll", Method: http.MethodGet, Pattern: "/api/lru", HandlerFunc: ch.getAllHandler, MiddlewareAuthFunc: logMiddleware},
		{Name: "Evict", Method: http.MethodDelete, Pattern: "/api/lru/{key}", HandlerFunc: ch.evictHandler, MiddlewareAuthFunc: logMiddleware},
		{Name: "EvictAll", Method: http.MethodDelete, Pattern: "/api/lru", HandlerFunc: ch.evictAllHandler, MiddlewareAuthFunc: logMiddleware},
	}

	router := mux.NewRouter().StrictSlash(false)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.MiddlewareAuthFunc(route.HandlerFunc))
	}
	return router
}

// middleware с логгированием данныз запроса, без авторизации
func logMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Вызов следующего обработчика
		handler.ServeHTTP(w, r)

		// Логирование запроса
		logrus.WithFields(logrus.Fields{
			"handler":      runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name(),
			"method":       r.Method,
			"url":          r.URL.Path,
			"responseTime": time.Since(startTime),
		}).Debug("Handled request")
	})
}
