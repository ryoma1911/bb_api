package api

import (
	"baseball_report/internal/repository"
	"encoding/json"

	"github.com/gorilla/mux"

	"net/http"
)

func GetScoreHandler(w http.ResponseWriter, r *http.Request) {
	//パスパラメータを取得
	vars := mux.Vars(r)
	id := vars["id"]

	//DB接続
	db, err := connect.ConnectOnly()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	//スコア情報を取得
	repo := &repository.DefaultRepository{}

	score, err := repo.GetScore(db, id)
	if err != nil {
		http.Error(w, "Error executing query: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	//結果を返却
	if len(score) != 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(score)
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "No score found"})
	}
}
