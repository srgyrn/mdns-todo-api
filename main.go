package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/srgyrn/mdns-todo-api/storage"
)

func main() {
	uri := "/items"
	db := storage.NewDBHandler()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		successfulResponseListener(w, "success")
	})

	http.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		var id int
		id, _ = strconv.Atoi(r.URL.Path[len(uri):])

		switch r.Method {
		case http.MethodGet:
			result, err := db.GetItems()

			if err != nil {
				http.Error(w, err.Error(), http.StatusNoContent)
				return
			}

			successfulResponseListener(w, result)
		case http.MethodPost:
			type req struct  {
				Content string `json:"content"`
			}
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			var bod req
			err = json.Unmarshal(b, &bod)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			result, err := db.AddNewItem(bod.Content)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			successfulResponseListener(w, result)
		case http.MethodDelete:
			if _, err := db.DeleteItem(id); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		default:
			http.NotFound(w, r)
		}
	})

	fmt.Print("Starting server...")
	log.Fatal(http.ListenAndServe(":8090", nil))
}

// successfulResponseListener sets necessary information to the response and prevents code duplication
func successfulResponseListener(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}