package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
)

var s store

func randomIdentifierOfLength(len int) string {
	buflen := hex.DecodedLen(len)
	buffer := make([]byte, buflen)
	rand.Read(buffer)

	return hex.EncodeToString(buffer)
}

func apiPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	switch r.Method {
	case "GET":
		posts := make([]post, 0)
		if err := s.withRoFile(func(f *os.File) error {
			var err error
			posts, err = loadPosts(f)
			return err
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if err := json.NewEncoder(w).Encode(struct {
			Posts []post `json:"posts"`
		}{posts}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case "OPTIONS":
		w.WriteHeader(http.StatusOK)
	case "POST":

		var newPost post
		if err := json.NewDecoder(r.Body).Decode(&struct {
			Post *post `json:"post"`
		}{&newPost}); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newPost.Pid = randomIdentifierOfLength(20)

		if err := s.withRWFile(func(f *os.File) error {
			posts, err := loadPosts(f)
			if err != nil {
				return err
			}
			posts = append(posts, newPost)
			return putPosts(f, posts)
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(struct {
			Post post `json:"post"`
		}{newPost}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
