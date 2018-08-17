USE mysql;

CREATE DATABASE IF NOT EXISTS test
  DEFAULT CHARSET utf8mb4
  COLLATE utf8mb4_unicode_ci;

CREATE USER 'foo'@'localhost'
  IDENTIFIED BY 'test1234';
GRANT ALL PRIVILEGES ON test.* TO 'foo'@'localhost'
IDENTIFIED BY 'test1234';

SET GLOBAL TIME_ZONE = '+8:00';
SET TIME_ZONE = '+8:00';

FLUSH PRIVILEGES;


USE test;
-- ----------------------------
-- Table structure for tbl_user
-- ----------------------------
DROP TABLE IF EXISTS `tbl_user`;
CREATE TABLE `tbl_user`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `username` varchar(20) NOT NULL,
  `password` varchar(20) NOT NULL,
  `secret` char(32) NOT NULL,
  `cellphone` varchar(20) NOT NULL,
  `email` varchar(20) DEFAULT NULL,
  `created_time` datetime(0) NOT NULL,
  `updated_time` datetime(0) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `UQE_tbl_user_cellphone`(`cellphone`) USING BTREE,
  UNIQUE INDEX `UQE_tbl_user_username`(`username`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci;

-- ----------------------------
-- Records of tbl_user
-- ----------------------------
INSERT INTO `tbl_user` VALUES (2, 'test', '46931023339d', 'XOP3z5O1kGn3Y3JlOOBlmZIHkOdC8lxu', '156××××××××', '', '2018-08-15 17:32:54', '2018-08-15 17:32:54');
