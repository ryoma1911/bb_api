#!/bin/sh
go mod init baseball_report

go get -u github.com/labstack/echo/v4 #echoのライブラリ
go get -u github.com/PuerkitoBio/goquery # スクレイピングのライブラリ
go get -u github.com/saintfish/chardet # 文字コードの判定用
go get -u golang.org/x/net/html/charset # 文字コードの変換用
go get -u github.com/go-sql-driver/mysql # Mysql用のドライバ
go get -u github.com/stretchr/testify #テスト用のライブラリ
go get -u github.com/DATA-DOG/go-sqlmock
go get -u github.com/joho/godotenv #環境変数用のライブラリ
go get -u github.com/robfig/cron/v3 #スケジュール実行ライブラリ
go get -u github.com/golang/mock/gomock
go get -u github.com/patrickmn/go-cache
go get -u github.com/gorilla/mux
go install -v github.com/go-delve/delve/cmd/dlv@latest


# cmdディレクトリに移動してビルド
cd /code/cmd

# ビルド（main.goをビルド）
go build -o main .

# mainが存在するか確認してから実行
if [ -f "./main" ]; then
    exec ./main
else
    echo "mainバイナリが存在しません。ビルドに失敗しました。"
    exit 1
fi

#tail -f /dev/null