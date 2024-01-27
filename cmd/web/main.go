package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	serve := http.NewServeMux()

	serve.HandleFunc("/", handleIndexPage)
	serve.HandleFunc("/kv", func (res http.ResponseWriter, req *http.Request)  {
		allowedMethods := []string{http.MethodGet, http.MethodDelete, http.MethodPost, http.MethodPut}
		res.Header().Set("Allow", strings.Join(allowedMethods, ","))

		switch req.Method {
			case http.MethodGet:
				handleGetValue(res, req);
			case http.MethodPost:
				handleUpsert(res, req)
			case http.MethodPut:
				handleUpdate(res, req)
			case http.MethodDelete:
				handleDeleteKey(res, req)
			default:
				http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)	
		}
	})
	serve.HandleFunc("/kv/ttl", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPut {
			res.Header().Set("Allow", http.MethodPut)
			http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handleTTLUpdate(res, req);
	})

	port := os.Getenv("PORT")
	log.Printf("Server starting on PORT :%s", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), serve)
	
	if err != nil {
		log.Fatal("Server not started!", err)
	}
}