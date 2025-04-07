package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun_StartsServer(t *testing.T) {

	// ヘルスチェックエンドポイントにリクエストしてみる
	resp, err := http.Get("http://localhost:8080/health")
	assert.NoError(t, err)
	defer resp.Body.Close()

	// 期待するステータスコード
	assert.Equal(t, resp.StatusCode, http.StatusOK)
}
