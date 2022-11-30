package httphandler

import (
	"file-upload-demo/aws"
	"io"
	"mime"
	"mime/multipart"
	"net/http"

	"go.uber.org/zap"
)

type Handler struct {
	mux *http.ServeMux
	log *zap.SugaredLogger
	s3  *aws.S3
}

// New Handler
func New(s *http.ServeMux, logger *zap.SugaredLogger, s3 *aws.S3) *Handler {
	h := Handler{s, logger, s3}
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
	h.log.Infof("uploading image: %s", "filename")

	_, params, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	multipartReader := multipart.NewReader(r.Body, params["boundary"])
	// buf := make([]byte, 256)
	for {
		part, err := multipartReader.NextPart()
		if err == io.EOF {
			break
		}
		contentType := part.Header.Get("Content-Type")
		fname := part.FileName()
		h.log.Infow("read file part", "Content-Type", contentType, "fname", fname)
		// upload file to aws s3
		h.s3.UploadStream(r.Context(), part, fname, contentType)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("👋 ayy lmao\n"))
}
