CREATE TABLE IF NOT EXISTS `cheshi`.`departments` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `superior` INT(11) NULL DEFAULT NULL,
  `name` varchar(32) NOT NULL,
  UNIQUE KEY (`superior`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
