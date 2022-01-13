package app

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func handleDelete(client *storage.Client, bucketName string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling delete.")


		key, ok := mux.Vars(r)["fileName"]
		if !ok {
			fail(w, http.StatusBadRequest, "Filename missing.")
			return
		}

		if err := client.Bucket(bucketName).Object(key).Delete(context.Background()); err != nil {
			fail(w, http.StatusFailedDependency, "Delete: %v", err)
			return
		}
		log.Println("Blob deleted.")

		w.WriteHeader(http.StatusGone)
	}
}