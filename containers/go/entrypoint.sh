#!/bin/sh

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