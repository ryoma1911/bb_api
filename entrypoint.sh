#!/bin/sh
go mod init baseball_report

go get -u github.com/PuerkitoBio/goquery # スクレイピングのライブラリ
go get -u github.com/saintfish/chardet # 文字コードの判定用
go get -u golang.org/x/net/html/charset # 文字コードの変換用