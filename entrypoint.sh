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

#tail -f /dev/null