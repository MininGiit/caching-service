package http

import (
	"cachingService/internal/logger"
	"cachingService/internal/usecase"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type handler struct {
	ctx    context.Context
	uc     usecase.IUseCase
	logger logger.Logger
}

func newHandler(ctx context.Context, uc usecase.IUseCase, logger logger.Logger) *handler {
	return &handler{
		ctx:    ctx,
		uc:     uc,
		logger: logger,
	}
}

func (h *handler) initRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/lru/{key}", h.get).Methods("GET")
	router.HandleFunc("/api/lru", h.getAll).Methods("GET")
	router.HandleFunc("/api/lru", h.post).Methods("POST")
	router.HandleFunc("/api/lru/{key}", h.delete).Methods("DELETE")
	router.HandleFunc("/api/lru", h.deleteAll).Methods("DELETE")
	return router
}

func (h *handler) get(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	key := string(params["key"])
	h.logger.Debug("Handling GET request, key:", key)
	value, expiresAt, err := h.uc.Get(h.ctx, key)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		h.logger.Warn("Key not found", "key", key, "error", err)
		return
	}
	response := newResponseItem(key, value, expiresAt)
	h.logger.Debug("GET request succeed", "key", key)
	json.NewEncoder(w).Encode(response)
}

func (h *handler) getAll(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Handling GET request for all data")
	keys, values, err := h.uc.GetAll(h.ctx)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		h.logger.Warn("Cache is empty", "error", err)
		return
	}
	response := newResponseItems(keys, values)
	h.logger.Debug("GET request for all succeed")
	json.NewEncoder(w).Encode(response)
}

func (h *handler) post(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Handling POST request")
	var request requestItem
	json.NewDecoder(r.Body).Decode(&request)

	err := h.uc.Put(h.ctx, request.Key, request.Value, time.Second*time.Duration(request.TtlSeconds))
	if err != nil {
		h.logger.Warn("Invalid request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h.logger.Debug("POST request succeed")
	w.WriteHeader(http.StatusCreated)
}

func (h *handler) delete(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Handling DELETE request")
	params := mux.Vars(r)
	key := string(params["key"])
	_, err := h.uc.Evict(h.ctx, key)
	if err != nil {
		h.logger.Warn("Key not found", "key", key, "error", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	h.logger.Debug("DELETE request succeed")
	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) deleteAll(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Handling DELETE request for all data")
	h.uc.EvictAll(h.ctx)
	w.WriteHeader(http.StatusNoContent)
	h.logger.Debug("DELETE request for all succeed")
}
