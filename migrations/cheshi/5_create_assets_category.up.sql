CREATE TABLE IF NOT EXISTS `cheshi`.`assets_category` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `name` varchar(32) NOT NULL,
  `parent` INT(11) UNSIGNED NOT NULL DEFAULT '0',
  unique key pname (`parent`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `cheshi`.`assets_category` WRITE;
INSERT INTO `cheshi`.`assets_category`(id,name,parent) VALUES (1,"固定资产分类",0);
INSERT INTO `cheshi`.`assets_category`(id,name,parent) VALUES (2,"电子资产分类",0);
UNLOCK TABLES;
