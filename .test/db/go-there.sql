CREATE TABLE `go` (
    `path` text DEFAULT NULL,
    `target` text DEFAULT NULL,
    `user` text DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `users` (
    `username` varchar(255) DEFAULT NULL,
    `is_admin` tinyint(1) DEFAULT 0,
    `password_hash` varchar(255) DEFAULT NULL,
    `api_key_hash` varchar(255) DEFAULT NULL,
    `api_key_salt` varchar(30) DEFAULT NULL,
    UNIQUE KEY `users_api_key_salt_uindex` (`api_key_salt`),
    UNIQUE KEY `users_username_uindex` (`username`) USING HASH
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
