package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"AutoDeploy/extract"
	"AutoDeploy/nginxHandler"
	"AutoDeploy/upload"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(100 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	var targetDir string
	if jsonStr := r.FormValue("metadata"); jsonStr != "" {
		var meta struct {
			TargetDir string `json:"targetDir"`
		}
		if err := json.Unmarshal([]byte(jsonStr), &meta); err != nil {
			http.Error(w, "Invalid JSON metadata", http.StatusBadRequest)
			return
		}
		targetDir = meta.TargetDir
	}

	file, fileHeader, err := r.FormFile("zipfile")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	uploadDir := "./uploads"
	filePath, err := upload.UploadZip(file, fileHeader, uploadDir)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to upload file: %v", err), http.StatusInternalServerError)
		return
	}

	if err := extract.ExtractZip(filePath, targetDir); err != nil {
		http.Error(w, fmt.Sprintf("Failed to extract file: %v", err), http.StatusInternalServerError)
		return
	}

	targetDirString := strings.TrimPrefix(targetDir, "./")
	configPath := "/etc/nginx/sites-available/" + targetDirString + ".sample.local.conf"
	rootDir := "/var/www/" + targetDirString
	url := targetDirString + ".sample.local"

	err = nginxHandler.CreateNginxConfig(url, rootDir, configPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create nginx config: %v", err), http.StatusInternalServerError)
		return
	}

	err = nginxHandler.ReloadNginx()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to reload nginx: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded, extracted, and configured successfully! Served at: %s", url)
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
