-- MYSQL_USERに権限を付与
GRANT ALL PRIVILEGES ON *.* TO 'bbapi'@'%';
FLUSH PRIVILEGES;

DROP TABLE IF EXISTS matches;

CREATE TABLE matches (
    id INT AUTO_INCREMENT PRIMARY KEY,
    date DATE NOT NULL,
    home VARCHAR(50) NOT NULL,
    away VARCHAR(50) NOT NULL,
    league VARCHAR(50) NOT NULL,
    stadium VARCHAR(100) NOT NULL,
    starttime TIME NOT NULL,
    UNIQUE KEY link VARCHAR(255), -- 進捗ページのURL
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
);

CREATE TABLE scores (
    id INT AUTO_INCREMENT PRIMARY KEY,
    home_score VARCHAR(3),
    away_score VARCHAR(3),
    batter VARCHAR(30),
    inning VARCHAR(30) DEFAULT '試合前',
    result VARCHAR(100),
    match_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (match_id) REFERENCES matches(id)
);