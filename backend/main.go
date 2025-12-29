package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Configuration structure
type Config struct {
	Server struct {
		Port    string `json:"port"`
		BaseURL string `json:"base_url"`
	} `json:"server"`
	Storage struct {
		UploadDir            string `json:"upload_dir"`
		MaxFileSizeMB        int64  `json:"max_file_size_mb"`
		FileTTLMinutes       int    `json:"file_ttl_minutes"`
		DefaultDownloadLimit int    `json:"default_download_limit"`
	} `json:"storage"`
}

var (
	config Config
)

type FileMeta struct {
	ID            string    `json:"id"`
	OriginalName  string    `json:"name"`
	Path          string    `json:"-"`
	CreatedAt     time.Time `json:"created_at"`
	ExpiresAt     time.Time `json:"expires_at"`
	DownloadLimit int       `json:"download_limit"`
	DownloadCount int       `json:"download_count"`
}

// Thread-safe store
type Store struct {
	sync.RWMutex
	files map[string]FileMeta
}

var store = Store{
	files: make(map[string]FileMeta),
}

func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Println("config.json not found, using defaults")
		config.Server.Port = ":8080"
		config.Server.BaseURL = "http://localhost:8080"
		config.Storage.UploadDir = "./uploads"
		config.Storage.MaxFileSizeMB = 100
		config.Storage.FileTTLMinutes = 60
		config.Storage.DefaultDownloadLimit = 5
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatal("Failed to decode config:", err)
	}
}

func main() {
	loadConfig()

	// Ensure upload directory exists
	if err := os.MkdirAll(config.Storage.UploadDir, 0755); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	// CORS configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	r.Use(cors.New(corsConfig))

	// Routes
	// Set max multipart memory
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.POST("/upload", handleUpload)
	r.GET("/download/:id", handleDownload)
	r.GET("/info/:id", handleInfo)

	// Start cleanup routine
	go cleanupRoutine()

	log.Printf("Server starting on %s", config.Server.Port)
	r.Run(config.Server.Port)
}

func handleUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	maxSize := config.Storage.MaxFileSizeMB * 1024 * 1024
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("File too large (Max %dMB)", config.Storage.MaxFileSizeMB)})
		return
	}

	id := uuid.New().String()
	ext := filepath.Ext(file.Filename)
	saveFilename := id + ext
	savePath := filepath.Join(config.Storage.UploadDir, saveFilename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	now := time.Now()
	ttl := time.Duration(config.Storage.FileTTLMinutes) * time.Minute
	
	meta := FileMeta{
		ID:            id,
		OriginalName:  file.Filename,
		Path:          savePath,
		CreatedAt:     now,
		ExpiresAt:     now.Add(ttl),
		DownloadLimit: config.Storage.DefaultDownloadLimit,
		DownloadCount: 0,
	}

	store.Lock()
	store.files[id] = meta
	store.Unlock()

	downloadURL := fmt.Sprintf("%s/download/%s", config.Server.BaseURL, id)
	
	c.JSON(http.StatusOK, gin.H{
		"id":           id,
		"download_url": downloadURL,
		"expires_at":   meta.ExpiresAt,
		"download_limit": meta.DownloadLimit,
	})
}

func handleDownload(c *gin.Context) {
	id := c.Param("id")
	
	store.Lock() // Lock for writing download count
	meta, exists := store.files[id]
	
	if !exists {
		store.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found or expired"})
		return
	}

	if time.Now().After(meta.ExpiresAt) {
		delete(store.files, id) // Lazy cleanup
		store.Unlock()
		os.Remove(meta.Path)
		c.JSON(http.StatusNotFound, gin.H{"error": "File expired"})
		return
	}

	if meta.DownloadLimit > 0 && meta.DownloadCount >= meta.DownloadLimit {
		delete(store.files, id)
		store.Unlock()
		os.Remove(meta.Path)
		c.JSON(http.StatusNotFound, gin.H{"error": "Download limit reached"})
		return
	}

	// Increment count
	meta.DownloadCount++
	store.files[id] = meta // Update map
	
	// Check if this was the last allowed download
	shouldDelete := false
	if meta.DownloadLimit > 0 && meta.DownloadCount >= meta.DownloadLimit {
		shouldDelete = true
		delete(store.files, id)
	}
	store.Unlock()

	// If we need to delete, we do it AFTER serving the file?
	// Actually, if we delete the file now, we can't serve it.
	// So we serve it, and if it was the last one, we might need a background cleanup or just delete it now?
	// c.FileAttachment serves the file. If we delete it immediately, it might fail if it's large.
	// Better approach: Keep it in map but mark as expired? Or just delete after serving.
	// Gin's c.FileAttachment writes to response.
	
	c.FileAttachment(meta.Path, meta.OriginalName)

	if shouldDelete {
		// Wait a bit to ensure file is served? 
		// Actually, FileAttachment doesn't block until done? It probably does.
		// But to be safe on Windows (file locking), we might need to be careful.
		// Let's just remove it from store (already done above) so no new downloads can start.
		// The file on disk can be cleaned up by the cleanup routine or here.
		// On Windows, you can't delete an open file.
		// So we will let the cleanup routine handle it or try to delete it later.
		// But we removed it from 'store', so it's effectively gone for new requests.
		// Let's rely on cleanup routine or a deferred delete.
		go func(path string) {
			time.Sleep(10 * time.Second) // Give it time to finish sending
			os.Remove(path)
			log.Printf("Deleted file after max downloads: %s", path)
		}(meta.Path)
	}
}

func handleInfo(c *gin.Context) {
	id := c.Param("id")
	
	store.RLock()
	meta, exists := store.files[id]
	store.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	c.JSON(http.StatusOK, meta)
}

func cleanupRoutine() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		now := time.Now()
		var toDelete []string

		store.RLock()
		for id, meta := range store.files {
			if now.After(meta.ExpiresAt) {
				toDelete = append(toDelete, id)
			}
		}
		store.RUnlock()

		for _, id := range toDelete {
			store.Lock()
			meta, ok := store.files[id]
			if ok {
				delete(store.files, id)
			}
			store.Unlock()

			if ok {
				os.Remove(meta.Path)
				log.Printf("Cleaned up expired file: %s", id)
			}
		}
	}
}
