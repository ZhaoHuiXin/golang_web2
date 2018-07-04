CREATE TABLE IF NOT EXISTS `cheshi`.`access` (
  `role_id` INT(11) NOT NULL DEFAULT 0,
  `user_id` INT(11) NOT NULL DEFAULT 0,
  `feat_id` INT(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`role_id`,`user_id`, `feat_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `cheshi`.`access` WRITE;
INSERT INTO `cheshi`.`access` (`role_id`,`feat_id`) VALUES
(1,1),
(1,2),
(1,3),
(1,4),
(1,5),
(1,6),
(1,7),
(1,8),
(1,9),
(1,10),
(1,11),
(1,12),
(1,13),
(1,14),
(1,15),
(1,16),
(1,20);
UNLOCK TABLES;
