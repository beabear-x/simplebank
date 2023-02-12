CREATE TABLE `users` (
  `username` varchar(255) PRIMARY KEY,
  `hashed_password` varchar(255) NOT NULL,
  `full_name` varchar(255) NOT NULL,
  `email` varchar(255) UNIQUE NOT NULL,
  `password_changed_at` timestamp NOT NULL DEFAULT "1970-01-01 08:00:00",
  `created_at` timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

ALTER TABLE `accounts` ADD FOREIGN KEY (`owner`) REFERENCES `users` (`username`);

-- CREATE UNIQUE INDEX `accounts_index_1` ON `accounts` (`owner`, `currency`);
ALTER TABLE `accounts` ADD CONSTRAINT `accounts_index_1` UNIQUE (`owner`, `currency`);