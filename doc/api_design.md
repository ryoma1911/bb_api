# API設計書

## 🎯 目的
- リアルタイムで試合情報と進捗を取得し、JSON形式で提供するAPI

## 📘 API仕様
### 1. GET /matches
- **説明**: 当日の試合情報を取得
- **リクエストパラメータ**: 無し

#### レスポンス例
```json
{
  "セ・リーグ": [
    {
      "away": "中日",
      "date": "2025-04-06",
      "home": "ヤクルト",
      "id": 1,
      "league": "セ・リーグ",
      "stadium": "神宮",
      "starttime": "13:00:00"
    },
    {
      "away": "DeNA",
      "date": "2025-04-06",
      "home": "広島",
      "id": 2,
      "league": "セ・リーグ",
      "stadium": "マツダスタジアム",
      "starttime": "13:00:00"
    }
  ]
}
```

### 2. GET /scores/{$matchid}
- **説明**: 当日の試合進捗を取得
- **リクエストパラメータ**:
  - `matchid` (optional): フィルタリングするmatchid

#### レスポンス例
##### 試合中
```json
{
    {
      "id": 1,
      "homescore": "1",
      "awayscore": "1",
      "batter": "渡部 聖弥",
      "inning": "3回裏",
      "result": "左2塁打"
    }
}
```
##### 試合前
```json
{
    {
      "id": 1,
      "homescore": "",
      "awayscore": "",
      "batter": "",
      "inning": "試合前",
      "result": ""
    }
}
```
##### 試合終了
```json
{
    {
      "id": 1,
      "homescore": "5",
      "awayscore": "3",
      "batter": "",
      "inning": "試合終了",
      "result": ""
    }
}
```