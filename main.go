package main

import (
	"cloud.google.com/go/firestore"
	"context"
	json2 "encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/api/iterator"
	"log"
	"net/http"
)

func createClient(ctx context.Context) *firestore.Client {
	// Sets your Google Cloud Platform project ID.
	projectID := "dummy-project"

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	// Close client when done with
	// defer client.Close()
	return client
}

type User struct {
	First string `json:"first"`
	Last  string `json:"last"`
	Born  string `json:"born"`
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	ctx := context.Background()
	client := createClient(ctx)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		iter := client.Collection("users").Documents(ctx)
		result := make([]User, 0)

		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to iterate: %v", err)
			}
			userRaw := doc.Data()
			user := User{
				First: userRaw["first"].(string),
				Last:  userRaw["last"].(string),
				Born:  userRaw["born"].(string),
			}
			result = append(result, user)
		}
		json, err := json2.Marshal(result)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	})

	r.Post("/users", func(w http.ResponseWriter, r *http.Request) {
		parseErr := r.ParseForm()
		if parseErr != nil {
			log.Fatalf("Failed to parse form: %v", parseErr)
			w.Write([]byte("failed to parse form"))
			return
		}
		postForm := r.PostForm
		fmt.Printf("%v", postForm)
		first := postForm.Get("first")
		last := postForm.Get("last")
		born := postForm.Get("born")

		_, _, err := client.Collection("users").Add(ctx, map[string]interface{}{
			"first": first,
			"last":  last,
			"born":  born,
		})
		if err != nil {
			log.Fatalf("Failed adding alovelace: %v", err)
		}

		user := User{
			First: first,
			Last:  last,
			Born:  born,
		}
		json, err := json2.Marshal(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		r.Header.Set("Content-Type", "application/json")
		w.Write(json)
	})

	http.ListenAndServe(":3000", r)
}
