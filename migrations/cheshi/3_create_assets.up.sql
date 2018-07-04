CREATE TABLE IF NOT EXISTS `cheshi`.`assets` (
  `id` int NOT NULL PRIMARY KEY AUTO_INCREMENT,
  `code` varchar(32) NULL UNIQUE,
  `label` varchar(32) NOT NULL,
  `kind` int NULL,
  `model` varchar(32) NULL,
  `serial` varchar(32) NULL,
  `buy_at` DATE NULL,
  `buy_val` DECIMAL(8,2) NULL,
  `comment` varchar(255) NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
