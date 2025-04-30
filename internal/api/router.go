package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	//エンドポイントを設定
	r.HandleFunc("/matches", GetMatchesHandler).Methods("GET")
	r.HandleFunc("/scores/{id}", GetScoreHandler).Methods("GET")

	//ヘルスチェックも追加
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	
	return r
}
