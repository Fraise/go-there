CREATE TABLE `users` (
    `id` int AUTO_INCREMENT PRIMARY KEY,
    `username` varchar(255) DEFAULT NULL,
    `is_admin` tinyint(1) DEFAULT 0,
    `password_hash` varchar(255) DEFAULT NULL,
    `api_key_hash` varchar(255) DEFAULT NULL,
    INDEX (username),
    INDEX (password_hash),
    UNIQUE (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE `go` (
    `path` varchar(255) DEFAULT NULL,
    `target` text DEFAULT NULL,
    `user_id` int,
    INDEX (path),
    FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
