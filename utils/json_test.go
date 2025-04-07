package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToJSON(t *testing.T) {
	t.Run("1league2games", func(t *testing.T) {
		//引数を定義
		testdata := []map[string]interface{}{
			{
				"id":        1,
				"date":      "2025-03-28",
				"home":      "Yankees",
				"away":      "Red Sox",
				"league":    "セ・リーグ",
				"stadium":   "Yankee Stadium",
				"starttime": "19:00",
			},
			{
				"id":        2,
				"date":      "2025-03-28",
				"home":      "Dodgers",
				"away":      "Giants",
				"league":    "セ・リーグ",
				"stadium":   "Dodger Stadium",
				"starttime": "18:30",
			},
		}

		// 期待値を設定
		expected := map[string][]map[string]interface{}{
			"セ・リーグ": {
				{
					"id":        1,
					"date":      "2025-03-28",
					"home":      "Yankees",
					"away":      "Red Sox",
					"league":    "セ・リーグ",
					"stadium":   "Yankee Stadium",
					"starttime": "19:00",
				},
				{
					"id":        2,
					"date":      "2025-03-28",
					"home":      "Dodgers",
					"away":      "Giants",
					"league":    "セ・リーグ",
					"stadium":   "Dodger Stadium",
					"starttime": "18:30",
				},
			},
		}

		//テスト実施
		result, err := ConvertToJSON(testdata, "league")
		assert.NoError(t, err)

		//期待値と結果が一致していること（フィールドの順序は問わない）
		assert.Equal(t, expected, result)
	})
	t.Run("2leagues4games", func(t *testing.T) {
		// 引数を定義
		testdata := []map[string]interface{}{
			// セ・リーグの試合
			{
				"id":        1,
				"date":      "2025-03-28",
				"home":      "Yankees",
				"away":      "Red Sox",
				"league":    "セ・リーグ",
				"stadium":   "Yankee Stadium",
				"starttime": "19:00",
			},
			{
				"id":        2,
				"date":      "2025-03-28",
				"home":      "Dodgers",
				"away":      "Giants",
				"league":    "セ・リーグ",
				"stadium":   "Dodger Stadium",
				"starttime": "18:30",
			},
			// パ・リーグの試合
			{
				"id":        3,
				"date":      "2025-03-29",
				"home":      "SoftBank Hawks",
				"away":      "Lions",
				"league":    "パ・リーグ",
				"stadium":   "PayPay Dome",
				"starttime": "14:00",
			},
			{
				"id":        4,
				"date":      "2025-03-29",
				"home":      "Eagles",
				"away":      "Marines",
				"league":    "パ・リーグ",
				"stadium":   "Rakuten Seimei Park",
				"starttime": "13:00",
			},
		}

		// 期待値を設定
		expected := map[string][]map[string]interface{}{
			"セ・リーグ": {
				{
					"id":        1,
					"date":      "2025-03-28",
					"home":      "Yankees",
					"away":      "Red Sox",
					"league":    "セ・リーグ",
					"stadium":   "Yankee Stadium",
					"starttime": "19:00",
				},
				{
					"id":        2,
					"date":      "2025-03-28",
					"home":      "Dodgers",
					"away":      "Giants",
					"league":    "セ・リーグ",
					"stadium":   "Dodger Stadium",
					"starttime": "18:30",
				},
			},
			"パ・リーグ": {
				{
					"id":        3,
					"date":      "2025-03-29",
					"home":      "SoftBank Hawks",
					"away":      "Lions",
					"league":    "パ・リーグ",
					"stadium":   "PayPay Dome",
					"starttime": "14:00",
				},
				{
					"id":        4,
					"date":      "2025-03-29",
					"home":      "Eagles",
					"away":      "Marines",
					"league":    "パ・リーグ",
					"stadium":   "Rakuten Seimei Park",
					"starttime": "13:00",
				},
			},
		}

		// テスト実施
		result, err := ConvertToJSON(testdata, "league")
		assert.NoError(t, err)

		// 期待値と結果が一致していること（フィールドの順序は問わない）
		assert.Equal(t, expected, result)
	})

}
