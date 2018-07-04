CREATE TABLE IF NOT EXISTS `cheshi`.`features` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `auth` TINYINT(1) NULL DEFAULT NULL,
  `name` varchar(32) NULL DEFAULT NULL,
  `methods` varchar(32) NOT NULL DEFAULT '',
  `path` varchar(64) NOT NULL DEFAULT '',
  UNIQUE KEY (`methods`,`path`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `cheshi`.`features` WRITE;
INSERT INTO `cheshi`.`features`(`id`,`auth`,`name`,`methods`,`path`) VALUES
(1,0,"login","GET","/login"),
(2,0,"logout","GET","/logout"),
(3,0,"weixin auth","GET","/weixin/auth"),
(4,1,"users list","GET","/users"),
(5,1,"user create","POST","/user"),
(6,1,"user update","PUT","/user"),
(7,1,"weixin waitbind","GET","/weixin/waitbind"),
(8,1,"dep. list","GET","/departments"),
(9,1,"dep. create","POST","/department"),
(10,1,"dep. update","PUT","/department"),
(11,1,"dep. roles all list","GET","/departments/roles"),
(12,1,"dep. roles list","GET","/department/{depid:[0-9]+}/roles"),
(13,1,"dep. role create","POST","/department/{depid:[0-9]+}/role"),
(14,1,"dep. role update","PUT","/department/{depid:[0-9]+}/role"),
(15,1,"role access get","GET","/role/{role_id:[0-9]+}/access"),
(16,1,"role access update","PUT","/role/{role_id:[0-9]+}/access"),
(20,1,"home","GET","/");
UNLOCK TABLES;
