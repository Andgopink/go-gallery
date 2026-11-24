package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var galleryTemplate *template.Template
var uploadsDir = "./uploads"

func init() {
	galleryTemplate = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/api/images", handleAPIImages)
	http.HandleFunc("/api/images/upload", handleUpload)
	http.HandleFunc("/uploads/", handleUploads)

	fmt.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	galleryTemplate.Execute(w, nil)
}

func handleAPIImages(w http.ResponseWriter, r *http.Request) {
	images, err := getUploadedImages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(images); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filePath := header.Filename
	validExtensions := []string{".jpg", ".jpeg", ".png", ".webp"}
	validExt := false
	for _, ext := range validExtensions {
		if strings.HasSuffix(strings.ToLower(filePath), ext) {
			validExt = true
			break
		}
	}

	if !validExt {
		http.Error(w, "Invalid file extension. Only .jpg, .jpeg, .png, .webp allowed", http.StatusBadRequest)
		return
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	baseName := filepath.Base(filepath.Dir(filePath))
	dirName := filepath.Dir(baseName)

	timestamp := time.Now().UnixNano()
	filename := filepath.Join(dirName, fmt.Sprintf("%s_%s%s", timestamp, baseName, ext))

	savedPath := filepath.Join(uploadsDir, filename)
	if err := os.MkdirAll(savedPath, 0755); err != nil {
		http.Error(w, "Failed to create directory", http.StatusInternalServerError)
		return
	}

	dst, err := os.Create(savedPath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Uploaded: %s\n", filename)
}

func handleUploads(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/uploads/") {
		fileServer := http.FileServer(http.Dir("./uploads"))
		fileServer.ServeHTTP(w, r)
		return
	}
	http.NotFound(w, r)
}

func getUploadedImages() ([]string, error) {
	images := make([]string, 0)
	err := filepath.Walk(uploadsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			images = append(images, info.Name())
		}
		return nil
	})
	return images, err
}