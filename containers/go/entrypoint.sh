#!/bin/sh
cd /code/cmd

# ビルド済みバイナリを起動（ビルドはDockerfileの中でやる）
if [ -f "./main" ]; then
    exec ./main
else
    echo "mainバイナリが存在しません。ビルドに失敗しました。"
    exit 1
fi
