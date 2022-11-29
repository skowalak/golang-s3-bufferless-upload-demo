package httphandler

import (
	"net/http"

	"go.uber.org/zap"
)

type Handler struct {
	mux *http.ServeMux
	log *zap.SugaredLogger
}

// New Handler
func New(s *http.ServeMux, logger *zap.SugaredLogger) *Handler {
	h := Handler{s, logger}
	h.registerRoutes()

	return &h
}

func (h *Handler) registerRoutes() {
	h.mux.HandleFunc("/", h.hello)
	h.mux.HandleFunc("/api/v1/firmwares/123/image", h.uploadFileToS3)
}

func (h *Handler) hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Hello World!"))
}

func (h *Handler) uploadFileToS3(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(201)
	w.Write([]byte("ayy lmao\n"))
}
