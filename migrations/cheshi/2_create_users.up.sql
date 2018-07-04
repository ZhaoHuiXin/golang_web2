CREATE TABLE IF NOT EXISTS `cheshi`.`users` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `email` varchar(32) UNIQUE,
  `serial` varchar(32) UNIQUE,
  `token` bigint,
  `dep_id` INT(11) UNSIGNED NULL DEFAULT NULL,
  `role_id` INT(11) UNSIGNED NULL DEFAULT NULL,
  `sex` tinyint(1),
  `name` varchar(32),
  `nickname` varchar(32),
  `unionid` varchar(32) NOT NULL UNIQUE,
  `headimgurl` varchar(255),
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


LOCK TABLES `cheshi`.`users` WRITE;
INSERT INTO `cheshi`.`users`(id,email,serial,name,unionid,role_id)
VALUES(1,"xdj@autoforce.net","13","xudejian","oooUVxCXgXZCIPdeBG-rgSpPGRIg",1);
UNLOCK TABLES;
