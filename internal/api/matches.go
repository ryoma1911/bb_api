package api

import (
	db "baseball_report/internal/config"
	"baseball_report/internal/repository"
	"baseball_report/utils"
	"encoding/json"

	"net/http"
	"time"
)

var connect db.DBHandler = &db.DBService{}

func GetMatchesHandler(w http.ResponseWriter, r *http.Request) {

	todate := time.Now().Format("2006/01/02")
	query := "SELECT id, date, home, away, league, stadium, starttime, status FROM matches WHERE date ='" + todate + "'"
	//DB接続
	dsn, err := connect.GetDSNFromEnv("/code/.env")
	if err != nil {
		http.Error(w, "Get dsn error", http.StatusInternalServerError)
		return
	}
	db, err := connect.ConnectOnly(dsn)
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	//試合情報を取得
	repo := &repository.DefaultRepository{}

	matches, err := repo.GetMatchAPI(db, query)
	if err != nil {
		http.Error(w, "Error executing query: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	if len(matches) != 0 {
		//リーグをヘッダーとしたJSON形式に変換
		result, err := utils.ConvertToJSON(matches, "league")
		if err != nil {
			http.Error(w, "Error converting to JSON: "+err.Error(), http.StatusInternalServerError)
			return
		}

		//結果を返却
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "No matches found"})
	}

}
