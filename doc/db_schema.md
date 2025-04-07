
#### 2. **`db_schema.md`**  
データベースのテーブル設計やER図を記載。
# DBスキーマ設計書

## 🎯 目的
試合情報を格納するためのデータベース設計

### テーブル：matches

| カラム名     | 型           | 説明                    |
|--------------|--------------|-------------------------|
| id           | INT          | 主キー、自動インクリメント |
| date         | DATE         | 試合の日付              |
| home         | VARCHAR(50)  | ホームチーム名         |
| away         | VARCHAR(50)  | アウェイチーム名       |
| league       | VARCHAR(50)  | リーグ名                |
| stadium      | VARCHAR(100) | スタジアム名            |
| starttime    | TIME         | 試合開始時刻            |
| link         | VARCHAR(255) | 試合進捗のURL           |
| created_at   | TIMESTAMP    | 作成日時（自動）        |

---

### テーブル：scores

| カラム名      | 型           | 説明                        |
|---------------|--------------|-----------------------------|
| id            | INT          | 主キー、自動インクリメント     |
| match_id      | INT          | `matches.id` への外部キー      |
| home_score    | VARCHAR(50)  | ホームチームスコア           |
| away_score    | VARCHAR(50)  | アウェイチームスコア         |
| batter        | VARCHAR(50)  | 打席の選手名                  |
| inning        | VARCHAR(20)  | イニング                     |
| result        | VARCHAR(100) | 投打の結果                   |
| created_at    | TIMESTAMP    | 作成日時（自動）              |

---