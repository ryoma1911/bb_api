package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRun_StartsServer(t *testing.T) {
	// 別ゴルーチンでサーバを起動（短時間だけ）
	go func() {
		err := Run()
		assert.NoError(t, err)
	}()

	// サーバ起動を待つ
	time.Sleep(500 * time.Millisecond)

	// ヘルスチェックエンドポイントにリクエストしてみる
	resp, err := http.Get("http://localhost:8080/health")
	assert.NoError(t, err)
	defer resp.Body.Close()

	// 期待するステータスコード
	assert.Equal(t, resp.StatusCode, http.StatusOK)
}
