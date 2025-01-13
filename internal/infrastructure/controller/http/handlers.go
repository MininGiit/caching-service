package http

import (
	"cachingService/internal/usecase"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	ctx context.Context
	uc 	usecase.IUseCase
}

func NewHandler(ctx context.Context, uc usecase.IUseCase) *Handler {
	return &Handler{
		ctx: ctx,
		uc: uc,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/lru/{key}", h.get).Methods("GET")
	router.HandleFunc("/api/lru", h.getAll).Methods("GET")
	router.HandleFunc("/api/lru", h.post).Methods("POST")
	router.HandleFunc("/api/lru/{key}", h.delete).Methods("DELETE")
	router.HandleFunc("/api/lru", h.deleteAll).Methods("DELETE")
	return router
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
    key := string(params["key"])
	value, expiresAt, err := h.uc.Get(h.ctx, key)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	response := NewResponseItem(key, value, expiresAt)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request){
	keys, values, err := h.uc.GetAll(h.ctx)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	response := NewResponseItems(keys, values)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) post(w http.ResponseWriter, r *http.Request){
	var request RequestItem
	json.NewDecoder(r.Body).Decode(&request)

	err := h.uc.Put(h.ctx, request.Key, request.Value, time.Second * time.Duration(request.TtlSeconds))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
    key := string(params["key"])
	_, err := h.uc.Evict(h.ctx, key)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) deleteAll(w http.ResponseWriter, r *http.Request){
	h.uc.EvictAll(h.ctx)
	w.WriteHeader(http.StatusNoContent)
}