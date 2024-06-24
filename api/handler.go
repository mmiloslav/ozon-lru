// пакет с api
package api

import (
	"common"
	"encoding/json"
	"io"
	"lru"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// errors
const (
	errInternal    = "INTERNAL_ERROR"
	errWrongParams = "WRONG_PARAMS"
	errNotFound    = "NOT_FOUND"
)

// базовая структура ответа, содержится во всех структурах ответа
type baseResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// метод устанавливает успешное значение ответов api-методов
func (resp *baseResponse) SetSuccess() {
	resp.Error = ""
	resp.Success = true
}

// метод устанавливает ошибку в ответах api-методов
func (resp *baseResponse) SetError(err string) {
	resp.Error = err
	resp.Success = false
}

// ф-я записи ответа api-методов
func writeResponse(w http.ResponseWriter, response interface{}, statusCode int) {
	byteBody, err := json.Marshal(response)
	if err != nil {
		log.Errorf("couldnt marshal during writeAnswer of object [%+v] with error [%s]", response, err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	io.WriteString(w, string(byteBody))
}

// структура ответа метода api/Ping
type pingResponse struct {
	baseResponse
}

// хендер пинга
func pingHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("health check")

	resp := pingResponse{}
	resp.SetSuccess()

	writeResponse(w, resp, http.StatusOK)
}

// CacheHandler содержит LRU-кеш и методы для работы с ним
type CacheHandler struct {
	cache           lru.ILRUCache
	defaultCacheTTL time.Duration
}

// конструктор CacheHandler
func NewCacheHandler(cache lru.ILRUCache, defaultCacheTTL time.Duration) *CacheHandler {
	return &CacheHandler{
		cache:           cache,
		defaultCacheTTL: defaultCacheTTL,
	}
}

// структура запроса на добавление данных
type addDataRequest struct {
	Key        string      `json:"key"`
	Value      interface{} `json:"value"`
	TTLseconds int         `json:"ttl_seconds"`
}

// структура ответа метода на добавление данных
type addDataResponse struct {
	baseResponse
}

// putHandler HTTP-обработчик для добавления элемента в кеш
func (h *CacheHandler) putHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqData := addDataRequest{}
	resp := addDataResponse{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqData)
	if err != nil {
		log.Errorf("failed to decode put rq body with error [%s]", err.Error())
		resp.SetError(errInternal)
		writeResponse(w, resp, http.StatusBadRequest)
		return
	}

	ttl := h.defaultCacheTTL
	if reqData.TTLseconds > 0 {
		ttl = time.Second * time.Duration(reqData.TTLseconds)
	}

	err = h.cache.Put(ctx, reqData.Key, reqData.Value, ttl)
	if err != nil {
		log.Errorf("failed to put data in cache with error [%s]", err.Error())
		resp.SetError(errInternal)
		writeResponse(w, resp, http.StatusInternalServerError)
		return
	}

	resp.SetSuccess()
	writeResponse(w, resp, http.StatusCreated)
}

// структура ответа метода на получение одного элемента
type getDataResponse struct {
	baseResponse
	Key       string      `json:"key"`
	Value     interface{} `json:"value"`
	ExpiresAt int64       `json:"expires_at"`
}

// getHandler HTTP-обработчик для получения элемента из кеша
func (h *CacheHandler) getHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp := getDataResponse{}

	vars := mux.Vars(r)
	key := vars["key"]

	if !common.ValidString(key) {
		log.Error("key is empty")
		resp.SetError(errWrongParams)
		writeResponse(w, resp, http.StatusBadRequest)
		return
	}

	value, expiresAt, err := h.cache.Get(ctx, key)
	if err != nil {
		log.Errorf("failed to get data by key [%s] with error [%s]", key, err.Error())
		resp.SetError(errNotFound)
		writeResponse(w, resp, http.StatusNotFound)
		return
	}

	resp.Key = key
	resp.Value = value
	resp.ExpiresAt = expiresAt.Unix()

	resp.SetSuccess()
	writeResponse(w, resp, http.StatusOK)
}

// структура ответа метода на получение всех элементов
type getAllDataResponse struct {
	baseResponse
	Keys   []string      `json:"keys"`
	Values []interface{} `json:"values"`
}

// getAllHandler HTTP-обработчик для получения всех элементов из кеша
func (h *CacheHandler) getAllHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp := getAllDataResponse{}

	keys, values, err := h.cache.GetAll(ctx)
	if err != nil {
		log.Errorf("failed to get all data with error [%s]", err.Error())
		resp.SetError(errInternal)
		writeResponse(w, resp, http.StatusInternalServerError)
		return
	}

	resp.SetSuccess()

	if len(keys) == 0 {
		log.Info("cache is empty")
		writeResponse(w, resp, http.StatusNoContent)
		return
	}

	resp.Keys = keys
	resp.Values = values

	writeResponse(w, resp, http.StatusOK)
}

// структура ответа метода на удаление одного элемента
type evictDataResponse struct {
	baseResponse
}

// evictHandler HTTP-обработчик для удаления элемента из кеша
func (h *CacheHandler) evictHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp := evictDataResponse{}

	vars := mux.Vars(r)
	key := vars["key"]

	if !common.ValidString(key) {
		log.Error("key is empty")
		resp.SetError(errWrongParams)
		writeResponse(w, resp, http.StatusBadRequest)
		return
	}

	_, err := h.cache.Evict(ctx, key)
	if err != nil {
		log.Errorf("failed to evict data by key [%s] with error [%s]", key, err.Error())
		resp.SetError(errNotFound)
		writeResponse(w, resp, http.StatusNotFound)
		return
	}

	resp.SetSuccess()
	writeResponse(w, resp, http.StatusNoContent)
}

// структура ответа метода на удаление всех элементов
type evictAllDataResponse struct {
	baseResponse
}

// evictAllHandler HTTP-обработчик для очистки кеша
func (h *CacheHandler) evictAllHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp := evictAllDataResponse{}

	err := h.cache.EvictAll(ctx)
	if err != nil {
		log.Errorf("failed to evict all data  with error [%s]", err.Error())
		resp.SetError(errInternal)
		writeResponse(w, resp, http.StatusInternalServerError)
		return
	}

	resp.SetSuccess()
	writeResponse(w, resp, http.StatusNoContent)
}
