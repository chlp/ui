package rest

import (
	"embed"
	"fmt"
	"github.com/chlp/ui/pkg/logger"
	"io/fs"
	"net/http"
	"path"
	"path/filepath"
)

//go:embed swagger/*
var staticFiles embed.FS

func serveSwaggerFiles(mux *http.ServeMux) {
	staticFS, _ := fs.Sub(staticFiles, "swagger")
	indexHtmlFound := false
	swaggerYamlFound := false
	_ = fs.WalkDir(staticFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		switch path {
		case "index.html":
			indexHtmlFound = true
		case "swagger.yaml":
			swaggerYamlFound = true
		}
		return nil
	})
	if !indexHtmlFound || !swaggerYamlFound {
		logger.Fatalf("important static http file not found")
	}
	indexHtmlData, err := staticFiles.ReadFile("swagger/index.html")
	if err != nil {
		logger.Fatalf("index.html not found")
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "" && r.URL.Path != "/" && r.URL.Path != "/index.html" {
			requestedFile := path.Clean(r.URL.Path)
			fileToServe := fmt.Sprintf("swagger%s", requestedFile)
			if _, err := staticFS.Open(requestedFile[1:]); err == nil {
				w.Header().Set("Content-Type", getMimeType(requestedFile))
				data, _ := staticFiles.ReadFile(fileToServe)
				_, _ = w.Write(data)
				return
			}
		}
		_, _ = w.Write(indexHtmlData)
		return
	})
}

var mimeTypes = map[string]string{
	".html":  "text/html; charset=utf-8",
	".css":   "text/css",
	".js":    "application/javascript",
	".json":  "application/json",
	".xml":   "application/xml",
	".svg":   "image/svg+xml",
	".png":   "image/png",
	".jpg":   "image/jpeg",
	".jpeg":  "image/jpeg",
	".gif":   "image/gif",
	".ico":   "image/x-icon",
	".woff":  "font/woff",
	".woff2": "font/woff2",
	".ttf":   "font/ttf",
	".otf":   "font/otf",
	".eot":   "application/vnd.ms-fontobject",
	".mp4":   "video/mp4",
	".webm":  "video/webm",
	".ogg":   "audio/ogg",
	".mp3":   "audio/mpeg",
	".wav":   "audio/wav",
	".zip":   "application/zip",
	".pdf":   "application/pdf",
	".txt":   "text/plain; charset=utf-8",
}

func getMimeType(filename string) string {
	ext := filepath.Ext(filename)
	if mime, found := mimeTypes[ext]; found {
		return mime
	}
	return "application/octet-stream"
}
