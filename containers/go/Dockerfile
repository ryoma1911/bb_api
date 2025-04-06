# Goのイメージを指定
FROM golang:1.23-alpine
# コンテナのワークディレクトリを指定
WORKDIR /code
# ソースコードをコンテナにコピー
COPY . /code/
COPY entrypoint.sh /entrypoint.sh

#entrypoint.shに実行権限を付与
RUN chmod +x entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]


