# システムアーキテクチャ設計書

## 🎯 目的
システム全体のアーキテクチャと各コンポーネントの関係性

## 🏗️ アーキテクチャ構成
- **バックエンド**:
  - Go
  - MySQL
- **インフラ**:
  - AWS EC2
  - Docker
- **CI/CD**:
  - ・Github
  - ・AWS(codepipeline)
  - ・AWS(codedeploy)

## 🚀 技術選定理由
- Go: 高パフォーマンスなAPI開発が可能
- MySQL: データベースのスケーラビリティと安定性
- Docker: 開発・運用環境の統一
