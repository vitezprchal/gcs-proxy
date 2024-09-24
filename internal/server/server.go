package server

import (
	"gcs-proxy/config"
	"io"
	"mime"
	"net/http"
	"path"
	"strings"

	"cloud.google.com/go/storage"
)

type Server struct {
	httpServer *http.Server
	gcsClient  *storage.Client
	config     *config.Config
}

func NewServer(gcsClient *storage.Client, cfg *config.Config) *Server {
	return &Server{
		gcsClient: gcsClient,
		config:    cfg,
	}
}

func (s *Server) HandleFile(w http.ResponseWriter, r *http.Request) {
	fileName := strings.TrimPrefix(r.URL.Path, "/")

	if fileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	bucket := s.gcsClient.Bucket(s.config.GoogleBucketName)

	obj := bucket.Object(fileName)

	_, err := obj.Attrs(r.Context())
	if err == storage.ErrObjectNotExist {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error accessing file", http.StatusInternalServerError)
		return
	}

	reader, err := obj.NewReader(r.Context())
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	contentType := mime.TypeByExtension(path.Ext(fileName))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", contentType)

	_, err = io.Copy(w, reader)
	if err != nil {
		http.Error(w, "Error streaming file", http.StatusInternalServerError)
	}
}

func (s *Server) Start(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.HandleFile)

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return s.httpServer.ListenAndServe()
}

func InitServer(gcsClient *storage.Client, cfg *config.Config) *Server {
	return NewServer(gcsClient, cfg)
}
