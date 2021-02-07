/*
Navicat MariaDB Data Transfer

Source Server         : localhost
Source Server Version : 100321
Source Host           : localhost:3306
Source Database       : dgo

Target Server Type    : MariaDB
Target Server Version : 100321
File Encoding         : 65001

Date: 2020-04-10 12:49:04
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for cms_admin_option_log
-- ----------------------------
DROP TABLE IF EXISTS `cms_admin_option_log`;
CREATE TABLE `cms_admin_option_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_name` varchar(50) NOT NULL DEFAULT '' COMMENT '用户名',
  `user_id` int(11) NOT NULL COMMENT '用户ID',
  `path` varchar(128) NOT NULL DEFAULT '' COMMENT '路径',
  `method` varchar(50) NOT NULL DEFAULT '' COMMENT '请求方法',
  `option` varchar(50) NOT NULL DEFAULT '' COMMENT '操作',
  `created_at` datetime NOT NULL DEFAULT current_timestamp() COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='后台操作日志表';

-- ----------------------------
-- Records of cms_admin_option_log
-- ----------------------------

-- ----------------------------
-- Table structure for cms_admin_permissions
-- ----------------------------
DROP TABLE IF EXISTS `cms_admin_permissions`;
CREATE TABLE `cms_admin_permissions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '菜单名称',
  `pid` smallint(6) NOT NULL COMMENT '父级菜单ID',
  `path` varchar(128) NOT NULL DEFAULT '' COMMENT '路径',
  `is_show` tinyint(4) NOT NULL DEFAULT 1 COMMENT '导航展示 (0 否 1 是)',
  `is_record` tinyint(1) NOT NULL DEFAULT 1 COMMENT '记录日志 (0 否 1 是)',
  `is_modify` tinyint(1) NOT NULL DEFAULT 1 COMMENT '权限传递 (0 否 1 是)',
  `created_at` datetime NOT NULL DEFAULT current_timestamp() COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `name` (`name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=30 DEFAULT CHARSET=utf8mb4 COMMENT='权限菜单';

-- ----------------------------
-- Records of cms_admin_permissions
-- ----------------------------
INSERT INTO `cms_admin_permissions` VALUES ('1', '权限管理', '0', '', '1', '1', '1', '2020-01-10 02:16:36', '2020-01-10 02:16:38');
INSERT INTO `cms_admin_permissions` VALUES ('2', '用户管理', '1', '/user', '1', '1', '1', '2020-01-10 02:17:10', '2020-01-10 02:17:12');
INSERT INTO `cms_admin_permissions` VALUES ('3', '角色管理', '1', '/role', '1', '1', '1', '2020-01-10 02:17:34', '2020-01-10 02:17:36');
INSERT INTO `cms_admin_permissions` VALUES ('4', '添加用户页面', '2', '/useraddpage', '0', '1', '1', '2020-01-10 02:19:11', '2020-04-03 11:35:52');
INSERT INTO `cms_admin_permissions` VALUES ('5', '添加用户', '2', '/useradd', '0', '1', '1', '2020-01-10 02:19:38', '2020-01-10 02:19:41');
INSERT INTO `cms_admin_permissions` VALUES ('6', '编辑用户页面', '2', '/usereditpage', '0', '1', '1', '2020-01-10 22:30:50', '2020-01-10 22:30:52');
INSERT INTO `cms_admin_permissions` VALUES ('7', '编辑用户', '2', '/useredit', '0', '1', '1', '2020-01-10 22:31:44', '2020-01-10 22:31:46');
INSERT INTO `cms_admin_permissions` VALUES ('8', '删除用户', '2', '/userdel', '0', '1', '1', '2020-01-10 22:32:20', '2020-01-10 22:32:22');
INSERT INTO `cms_admin_permissions` VALUES ('9', '添加角色页面', '3', '/roleaddpage', '0', '1', '1', '2020-01-10 22:32:59', '2020-01-10 22:33:02');
INSERT INTO `cms_admin_permissions` VALUES ('10', '添加角色', '3', '/roleadd', '0', '1', '1', '2020-01-10 22:33:26', '2020-01-10 22:33:28');
INSERT INTO `cms_admin_permissions` VALUES ('11', '编辑角色页面', '3', '/roleeditpage', '0', '1', '1', '2020-01-10 22:33:53', '2020-01-10 22:33:55');
INSERT INTO `cms_admin_permissions` VALUES ('12', '编辑角色', '3', '/roleedit', '0', '1', '1', '2020-01-10 22:34:17', '2020-01-10 22:34:19');
INSERT INTO `cms_admin_permissions` VALUES ('13', '删除角色', '3', '/roledel', '0', '1', '1', '2020-01-10 22:36:47', '2020-01-10 22:36:49');
INSERT INTO `cms_admin_permissions` VALUES ('14', '获取角色权限', '3', '/getpermissionsofrole', '0', '1', '1', '2020-01-10 22:38:19', '2020-01-10 22:38:23');
INSERT INTO `cms_admin_permissions` VALUES ('15', '分配权限页面', '3', '/permissionsofrole', '0', '1', '1', '2020-01-10 22:39:22', '2020-01-10 22:39:24');
INSERT INTO `cms_admin_permissions` VALUES ('16', '权限保存', '3', '/permissionsofrolesave', '0', '1', '1', '2020-01-10 22:41:40', '2020-01-10 22:41:43');
INSERT INTO `cms_admin_permissions` VALUES ('17', '菜单管理', '1', '/permission', '1', '1', '0', '2020-01-10 22:42:40', '2020-01-10 22:42:43');
INSERT INTO `cms_admin_permissions` VALUES ('18', '添加菜单页面', '17', '/permissionaddpage', '0', '1', '0', '2020-01-10 22:43:14', '2020-01-10 22:43:16');
INSERT INTO `cms_admin_permissions` VALUES ('19', '添加菜单', '17', '/permissionadd', '0', '1', '0', '2020-01-10 22:33:43', '2020-01-10 22:43:44');
INSERT INTO `cms_admin_permissions` VALUES ('20', '编辑菜单页面', '17', '/permissioneditpage', '0', '1', '0', '2020-01-10 22:34:12', '2020-01-10 22:44:13');
INSERT INTO `cms_admin_permissions` VALUES ('21', '编辑菜单', '17', '/permissionedit', '0', '1', '0', '2020-01-10 22:34:33', '2020-01-10 22:44:34');
INSERT INTO `cms_admin_permissions` VALUES ('22', '删除菜单', '17', '/permissiondel', '0', '1', '0', '2020-01-10 22:34:57', '2020-01-10 22:44:58');
INSERT INTO `cms_admin_permissions` VALUES ('23', '日志管理', '1', '/logrecord', '1', '1', '0', '2020-01-11 02:52:24', '2020-01-11 02:52:24');
INSERT INTO `cms_admin_permissions` VALUES ('24', '删除日志', '23', '/logrecorddel', '0', '1', '0', '2020-01-11 04:10:18', '2020-01-11 04:10:18');
INSERT INTO `cms_admin_permissions` VALUES ('25', '数据库管理', '1', '/tables', '1', '0', '0', '2020-04-08 19:44:01', '2020-04-10 12:44:20');
INSERT INTO `cms_admin_permissions` VALUES ('26', '数据表优化', '25', '/tables/optimize', '0', '0', '0', '2020-04-08 19:44:29', '2020-04-10 12:44:25');
INSERT INTO `cms_admin_permissions` VALUES ('27', 'model生成器', '25', '/generatepage', '0', '0', '0', '2020-04-08 19:44:54', '2020-04-10 12:44:31');
INSERT INTO `cms_admin_permissions` VALUES ('28', 'model代码生成', '25', '/tables/generate', '0', '0', '0', '2020-04-08 19:45:16', '2020-04-10 12:44:36');
INSERT INTO `cms_admin_permissions` VALUES ('29', '更改用户状态', '2', '/userban', '0', '1', '1', '2020-04-20 14:46:31', '2020-04-20 14:46:31');
INSERT INTO `cms_admin_permissions`  VALUES ('30', '配置管理', '0', '', '1', '1', '1', '2020-08-21 14:25:27', '2020-08-21 14:25:33');
INSERT INTO `cms_admin_permissions`  VALUES ('31', '数据配置', '30', '/keys', '1', '0', '1', '2020-08-21 17:07:51', '2020-08-21 17:07:51');
INSERT INTO `cms_admin_permissions`  VALUES ('32', '数据配置排序', '31', '/keys/sort', '0', '1', '1', '2020-08-21 17:08:51', '2020-08-21 17:08:51');
INSERT INTO `cms_admin_permissions`  VALUES ('33', '数据配置删除', '31', '/keys/delete', '0', '1', '1', '2020-08-21 17:09:11', '2020-08-21 17:09:11');
INSERT INTO `cms_admin_permissions`  VALUES ('34', '新增数据配置页面', '31', '/keys/addpage', '0', '0', '1', '2020-08-21 17:09:40', '2020-08-21 17:09:40');
INSERT INTO `cms_admin_permissions`  VALUES ('35', '新增数据配置', '31', '/keys/add', '0', '1', '1', '2020-08-21 17:10:02', '2020-08-21 17:10:02');
INSERT INTO `cms_admin_permissions`  VALUES ('36', '编辑数据配置页面', '31', '/keys/editpage', '0', '0', '1', '2020-08-21 17:10:32', '2020-08-21 17:10:32');
INSERT INTO `cms_admin_permissions`  VALUES ('37', '编辑数据配置', '31', '/keys/edit', '0', '1', '1', '2020-08-21 17:10:51', '2020-08-21 17:10:51');




-- ----------------------------
-- Table structure for cms_admin_role_has_permissions
-- ----------------------------
DROP TABLE IF EXISTS `cms_admin_role_has_permissions`;
CREATE TABLE `cms_admin_role_has_permissions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `role_id` int(11) NOT NULL COMMENT '角色ID',
  `permission_id` int(11) NOT NULL COMMENT '权限ID',
  `created_at` datetime NOT NULL DEFAULT current_timestamp() COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `role_id` (`role_id`,`permission_id`)
) ENGINE=InnoDB AUTO_INCREMENT=232 DEFAULT CHARSET=utf8mb4 COMMENT='角色权限关联';

-- ----------------------------
-- Records of cms_admin_role_has_permissions
-- ----------------------------
INSERT INTO `cms_admin_role_has_permissions` VALUES ('187', '1', '1', '2020-01-12 16:06:47');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('188', '1', '2', '2020-01-12 16:06:47');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('189', '1', '4', '2020-01-12 16:06:47');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('190', '1', '5', '2020-01-12 16:06:47');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('191', '1', '6', '2020-01-12 16:06:47');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('192', '1', '7', '2020-01-12 16:06:47');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('193', '1', '8', '2020-01-12 16:06:47');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('194', '1', '3', '2020-01-12 16:06:47');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('195', '1', '9', '2020-01-12 16:06:47');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('196', '1', '10', '2020-01-12 16:06:47');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('197', '1', '11', '2020-01-12 16:06:47');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('198', '1', '12', '2020-01-12 16:06:47');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('199', '1', '13', '2020-01-12 16:06:47');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('200', '1', '14', '2020-01-12 16:06:47');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('201', '1', '15', '2020-01-12 16:06:47');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('202', '1', '16', '2020-01-12 16:06:47');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('203', '1', '17', '2020-01-12 16:06:47');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('204', '1', '18', '2020-01-12 16:06:47');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('205', '1', '19', '2020-01-12 16:06:47');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('206', '1', '20', '2020-01-12 16:06:47');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('207', '1', '21', '2020-01-12 16:06:47');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('208', '1', '22', '2020-01-12 16:06:47');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('209', '1', '23', '2020-01-12 16:06:47');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('210', '1', '24', '2020-01-12 16:06:47');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('227', '1', '25', '2020-04-09 16:09:57');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('228', '1', '26', '2020-04-09 16:09:57');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('229', '1', '27', '2020-04-09 16:09:57');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('230', '1', '28', '2020-04-09 16:09:57');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('231', '1', '29', '2020-04-20 14:46:43');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('232', '1', '30', '2020-08-21 17:17:18');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('233', '1', '31', '2020-08-21 17:11:14');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('234', '1', '32', '2020-08-21 17:11:14');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('235', '1', '33', '2020-08-21 17:11:14');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('236', '1', '34', '2020-08-21 17:11:14');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('237', '1', '35', '2020-08-21 17:11:14');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('238', '1', '36', '2020-08-21 17:11:14');
INSERT INTO `cms_admin_role_has_permissions` VALUES ('239', '1', '37', '2020-08-21 17:11:14');

-- ----------------------------
-- Table structure for cms_admin_roles
-- ----------------------------
DROP TABLE IF EXISTS `cms_admin_roles`;
CREATE TABLE `cms_admin_roles` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '角色名称',
  `created_at` datetime NOT NULL DEFAULT current_timestamp() COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

-- ----------------------------
-- Records of cms_admin_roles
-- ----------------------------
INSERT INTO `cms_admin_roles` VALUES ('1', '超级管理员', '2020-01-09 16:48:33', '2020-01-09 16:48:36');

-- ----------------------------
-- Table structure for cms_admin_user_has_roles
-- ----------------------------
DROP TABLE IF EXISTS `cms_admin_user_has_roles`;
CREATE TABLE `cms_admin_user_has_roles` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `admin_id` int(11) NOT NULL COMMENT '用户ID',
  `role_id` int(11) NOT NULL COMMENT '角色ID',
  `created_at` datetime NOT NULL DEFAULT current_timestamp() COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `admin_id` (`admin_id`,`role_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联';

-- ----------------------------
-- Records of cms_admin_user_has_roles
-- ----------------------------
INSERT INTO `cms_admin_user_has_roles` VALUES ('1', '1', '1', '2020-01-10 22:39:19');

-- ----------------------------
-- Table structure for cms_admin_users
-- ----------------------------
DROP TABLE IF EXISTS `cms_admin_users`;
CREATE TABLE `cms_admin_users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '用户名',
  `password` varchar(255) NOT NULL DEFAULT '' COMMENT '密码',
  `avatar` varchar(255) NOT NULL DEFAULT '' COMMENT '头像路径',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '0 禁用 1启用',
  `login_ip` varchar(50) NOT NULL DEFAULT '' COMMENT '登录IP',
  `created_at` datetime NOT NULL DEFAULT current_timestamp() COMMENT '创建时间',
  `login_at` datetime NOT NULL COMMENT '最近登录时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='后台管理员表';

-- ----------------------------
-- Records of cms_admin_users
-- ----------------------------
INSERT INTO `cms_admin_users` VALUES ('1', 'admin', '$2a$10$IoOAzYgXD.R8NUZxlFMSke/Q3ByUyfq7XQpAwoVDLXo4rJkxZQtmy', '/upload/default.jpg', '1', '::1', '2020-01-08 18:11:32', '2020-04-10 12:35:05');

-- ----------------------------
-- Table structure for cms_burst_record
-- ----------------------------
DROP TABLE IF EXISTS `cms_burst_record`;
CREATE TABLE `cms_burst_record` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL COMMENT '用户id',
  `temp_folder_name` varchar(255) NOT NULL COMMENT '分片上传的临时缓存目录',
  `file_name` varchar(255) NOT NULL COMMENT '上传的文件名',
  `file_total_size` varchar(20) NOT NULL COMMENT '文件大小（字节）',
  `burst_count` int(11) NOT NULL DEFAULT 0 COMMENT '当前上传完成的分片块数',
  `burst_total` int(11) NOT NULL COMMENT '上传文件的总分片块数',
  `created_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='分片上传记录';

-- ----------------------------
-- Records of cms_burst_record
-- ----------------------------


-- ----------------------------
-- Table structure for cms_keys
-- ----------------------------
DROP TABLE IF EXISTS `cms_keys`;
CREATE TABLE `cms_keys` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL DEFAULT '' COMMENT '描述',
  `keyx1` varchar(64) NOT NULL DEFAULT '' COMMENT 'key1',
  `keyx2` varchar(64) NOT NULL DEFAULT '' COMMENT 'key2',
  `valuex` varchar(255) NOT NULL DEFAULT '' COMMENT 'value',
  `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '状态（0：草稿，1：正常）',
  `sort_num` int(11) NOT NULL DEFAULT 0 COMMENT '排序 （desc）',
  `created_at` datetime NOT NULL DEFAULT current_timestamp() COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `keyx` (`keyx1`,`keyx2`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='配置数据';

-- ----------------------------
-- Records of cms_keys
-- ----------------------------
