# mysqlDiff

[English](./README.md) | [中文](./README-zh.md)



**这是一个比较两个数据库之间的差异并生成更新的SQL语句脚本的工具。**



### mysqlDiff介绍



根据输入的两个url分别连接到两个数据库targetDB和sourceDB。mysqlDIff工具会分析两个数据库并给出由targetDB更新至sourceDB的sql脚本。



### 情况分类

1. 若targetDB含有sourceDB中不存在的表，则会给出对应的DROP TABLE语句。
2. 若sourceDB含有targetDB中不存在的表，如表student，则会给出表student的建表语句。
3. 若sourceDB与targetDB含有若干个名字相同的表，如tables，那么sqlDiff工具将会比较它们每一对相同名称的表之间的差异，并给出由targetDB->sourceDB的表结构更新语句，甚至是索引更新的sql语句。



### 如何使用

1. 根据您的操作系统下载对应的发行版（以下以linux环境为举例）。
2. 在命令行中切换目录至刚刚下载的发行版对应的目录，使用 sqlDiff_mysqlXXX -h 来查看如何使用（请确保您要链接的的mysql版本与发行版描述的一致）。
3. 执行二进制问文件，SQL脚本将会输出在下方的命令行中。





### 示例



#### 输入

./sqlDiff_mysql8.0.32 -sourceDB "root:root@tcp(localhost:3306)/db_cs_account1.2.0" -targetDB "root:root@tcp(localhost:3306)/db_cs_account"

#### 输出

#Table changes from database db_cs_account to database db_cs_account1.2.0:
CREATE TABLE `tb_test` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(32) DEFAULT 'zhangsan',
  `age` int NOT NULL,
  `hobbies` json DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


#The sql statement corresponding to the change of table tb_account ：
ALTER TABLE tb_account
ADD COLUMN delete_tm datetime COMMENT '账户删除时间',
ADD COLUMN testColumn varchar(234) NOT NULL,
DROP COLUMN account_name,
MODIFY COLUMN user_id int NOT NULL DEFAULT 34535 COMMENT '宵明用户ID',
MODIFY COLUMN platform_id int NOT NULL DEFAULT 5555 COMMENT '云平台id',
MODIFY COLUMN status json NOT NULL COMMENT '状态',
MODIFY COLUMN insert_tm datetime DEFAULT CURRENT_TIMESTAMP;


#Index update statement for table tb_account：
CREATE UNIQUE INDEX id ON tb_account(id);
CREATE UNIQUE INDEX yhi ON tb_account(platform_id,account_uid);
CREATE UNIQUE INDEX yhido ON tb_account(user_id,platform_id,insert_tm);
CREATE INDEX aka ON tb_account(user_id,id);















