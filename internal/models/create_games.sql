CREATE TABLE matches (
    id INT AUTO_INCREMENT PRIMARY KEY,
    date DATE NOT NULL,
    home VARCHAR(50) NOT NULL,
    away VARCHAR(50) NOT NULL,
    league VARCHAR(50) NOT NULL,
    stadium VARCHAR(100) NOT NULL,
    starttime TIME NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT '試合前', -- "予定", "試合中", "終了" などを想定
    link VARCHAR(255), -- 進捗ページのURL
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);