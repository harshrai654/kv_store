package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"kvstore.ttharsh.net/internal"
)

type requestBodyTTLUpdate struct {
	Key      string  `json:"key"`
	ExpireIn int64  `json:"expireIn,omitempty"`
}

type requestBodyUpdateKey struct {
	Key      string  `json:"key"`
	Value    string `json:"value,omitempty"`
}

type requestBodyUpsertKey struct {
	Key      string  `json:"key"`
	Value    string `json:"value,omitempty"`
	ExpireIn int64  `json:"expireIn,omitempty"`
}


func handleIndexPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.Header().Set("Allow", http.MethodGet)
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	res.Write([]byte("Welcome to key-value store"))
}


func handleGetValue(res http.ResponseWriter, req *http.Request) {
	key := req.URL.Query().Get("key")

	if key == "" {
		http.Error(res, "parameter key is required", http.StatusBadRequest)
		return
	}

	value := internal.GetValueFromKey(key)

	if value != nil {
		res.Write(value)
	} else {
		http.Error(res, "Key not found", http.StatusNotFound)
	}

}

func handleUpdate(res http.ResponseWriter, req *http.Request) {
	var requestData requestBodyUpdateKey
	err := json.NewDecoder(req.Body).Decode(&requestData)

	if err != nil {
		http.Error(res, "Failed to parse JSON request body", http.StatusBadRequest)
		return
	}

	key := requestData.Key

	if key == "" {
		http.Error(res, "Key required in request body", http.StatusBadRequest)
		return
	}

	value := requestData.Value

	if value == "" {
		res.WriteHeader(http.StatusNotModified)
		res.Write([]byte("Nothing to modify!"))
		return;
	}
	 
	result, rowCount := internal.UpdateKey(key, value)

	if result {
		if rowCount == 0 {
			res.Write([]byte(fmt.Sprintf("Key %s not found", key)))
		} else {
			res.Write([]byte(fmt.Sprintf("%s: %s", key, value)))
		}
	} else {
		http.Error(res, "Failed to update/insert key", http.StatusInternalServerError)
	}
}

func handleUpsert(res http.ResponseWriter, req *http.Request) {
	var requestData requestBodyUpsertKey
	err := json.NewDecoder(req.Body).Decode(&requestData)

	if err != nil {
		http.Error(res, "Failed to parse JSON request body", http.StatusBadRequest)
		return
	}

	key := requestData.Key

	if key == "" {
		http.Error(res, "Key required in request body", http.StatusBadRequest)
		return
	}

	value := requestData.Value
	expireIn := requestData.ExpireIn

	if value == "" {
		res.WriteHeader(http.StatusNotModified)
		res.Write([]byte("Nothing to modify!"))
		return;
	}
	
	var expired_at int64
	
	if expireIn > 0 {
		expired_at = time.Now().Unix() + int64(expireIn)
	} else {
		expired_at = 0
	}

	result, rowCount := internal.InsertKey(key, value, expired_at)
	
	if result {
		if rowCount == 0 {
			res.Write([]byte(fmt.Sprintf("Key %s not found", key)))
		} else {
			res.Write([]byte(fmt.Sprintf("%s: %s | expired_at: %d", key, value, expired_at)))
		}
	} else {
		http.Error(res, "Failed to update/insert key", http.StatusInternalServerError)
	}
}

func handleDeleteKey(res http.ResponseWriter, req *http.Request) {
	key := req.URL.Query().Get("key")

	if key == "" {
		http.Error(res, "parameter key is required", http.StatusBadRequest)
		return
	}

	result, rowCount := internal.DeleteKey(key)
	
	if result {
		if rowCount == 0 {
			res.Write([]byte(fmt.Sprintf("Key %s not found", key)))
		} else {
			res.Write([]byte(fmt.Sprintf("Key %s deleted successfully", key)))
		}
	} else {
		http.Error(res, "Failed to delete key", http.StatusInternalServerError)
	}
}

func handleTTLUpdate(res http.ResponseWriter, req *http.Request) {
	var requestData requestBodyTTLUpdate
	err := json.NewDecoder(req.Body).Decode(&requestData)

	if err != nil {
		http.Error(res, "Failed to parse JSON request body", http.StatusBadRequest)
		return
	}

	key := requestData.Key

	if key == "" {
		http.Error(res, "Key required in request body", http.StatusBadRequest)
		return
	}
	expireIn := requestData.ExpireIn

	if expireIn < 0 {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Invlaid expire_in value provide positive number of seconds"))
	}
	
	var expired_at int64
	
	if expireIn > 0 {
		expired_at = time.Now().Unix() + int64(expireIn)
	} else {
		expired_at = 0
	}

	result, rowCount := internal.UpdateTTL(key, expired_at)

	if result {
		if rowCount == 0 {
			res.Write([]byte(fmt.Sprintf("Key %s not found", key)))
		} else {
			res.Write([]byte(fmt.Sprintf("%s | expired_at: %d", key, expired_at)))
		}
	} else {
		http.Error(res, fmt.Sprintf("Failed to update TTL for key: %s", key), http.StatusInternalServerError)
	}
}