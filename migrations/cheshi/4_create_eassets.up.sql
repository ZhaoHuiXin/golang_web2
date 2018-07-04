CREATE TABLE IF NOT EXISTS `cheshi`.`eassets` (
  `id` int NOT NULL PRIMARY KEY AUTO_INCREMENT,
  `kind` int NULL,
  `op` int NULL,
  `label` varchar(32) NOT NULL,
  `model` varchar(32) NULL,
  `serial` varchar(32) NULL,
  `buy_at` DATE NULL,
  `cost` DECIMAL(8,2) NULL,
  `comment` varchar(255) NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
