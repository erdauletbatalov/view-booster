package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var (
	viewCounts = make(map[string]int)
	mu         sync.Mutex
)

// Логирование реального просмотра
func logRealView(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "Missing 'url' parameter", http.StatusBadRequest)
		return
	}

	mu.Lock()
	viewCounts[url]++
	count := viewCounts[url]
	mu.Unlock()

	response := map[string]interface{}{
		"url":     url,
		"views":   count,
		"message": "Real view logged",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Получение реального количества просмотров
func getRealViewCount(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "Missing 'url' parameter", http.StatusBadRequest)
		return
	}

	mu.Lock()
	count := viewCounts[url]
	mu.Unlock()

	response := map[string]interface{}{
		"url":   url,
		"views": count,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/log-view", logRealView)        // Логирование просмотров
	http.HandleFunc("/real-views", getRealViewCount) // Получение количества просмотров

	port := ":5001"
	fmt.Println("📊 Трекер запущен на", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
