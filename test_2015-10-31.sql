# ************************************************************
# Sequel Pro SQL dump
# Version 4135
#
# http://www.sequelpro.com/
# http://code.google.com/p/sequel-pro/
#
# Host: 127.0.0.1 (MySQL 5.6.26)
# Database: test
# Generation Time: 2015-10-30 23:02:12 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table task
# ------------------------------------------------------------

CREATE TABLE `task` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `num` int(11) DEFAULT '0' COMMENT '总链接数',
  `success` int(11) DEFAULT '0' COMMENT '总完成数',
  `create_time` int(11) DEFAULT '0' COMMENT '任务开始时间',
  `title` varchar(255) DEFAULT NULL COMMENT '名称',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# Dump of table url
# ------------------------------------------------------------

CREATE TABLE `url` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `duration` int(11) DEFAULT '0' COMMENT '耗时',
  `module` varchar(255) DEFAULT NULL COMMENT '模块',
  `controller` varchar(255) DEFAULT NULL COMMENT '控制器',
  `action` varchar(255) DEFAULT NULL COMMENT '操作',
  `url` varchar(255) DEFAULT NULL COMMENT '链接地址',
  `last_url` varchar(255) DEFAULT NULL COMMENT '最终转到的地址',
  `status` int(11) DEFAULT '0' COMMENT 'http状态码',
  `create_time` int(11) DEFAULT '0' COMMENT '记录添加时间',
  `server` varchar(255) DEFAULT NULL COMMENT 'server',
  `same` tinyint(4) NOT NULL DEFAULT '1' COMMENT '链接地址与最终地址是否一样 1是 0否',
  `err` varchar(255) DEFAULT NULL COMMENT '最近一次错误信息',
  `is_err` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否出错',
  `task_id` int(11) DEFAULT '0' COMMENT '任务id',
  `times` int(11) DEFAULT '0' COMMENT '访问次数',
  PRIMARY KEY (`id`),
  KEY `task_id` (`task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
