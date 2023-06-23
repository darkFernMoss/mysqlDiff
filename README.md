# mysqlDiff

[English](./README.md) | [中文](./README-zh.md)



**This is a tool that compares the differences between two databases and generates updated SQL statement scripts. **



### mysqlDiff Introduction



Based on the two URLs entered to connect to the two databases targetDB and sourceDB respectively. mysqlDIff tool will analyze the two databases and give the sql script for updating from targetDB to sourceDB.



### Case classification

1. if targetDB contains tables that do not exist in sourceDB, the corresponding DROP TABLE statement is given. 2. if sourceDB contains tables that do not exist in sourceDB, the corresponding DROP TABLE statement is given.
2. if the sourceDB contains a table that does not exist in the targetDB, such as the STUDENT table, a table build statement for the STUDENT table will be given.
3. if sourceDB and targetDB contain several tables with the same name, such as tables, then sqlDiff tool will compare the differences between each pair of them with the same name and give the table structure update statement from targetDB->sourceDB, even the index update sql statement.





### How to use

1. Download the corresponding release according to your operating system (the following is an example of linux environment).
2. Switch the directory in the command line to the directory corresponding to the release bution you just downloaded and use sqlDiff_mysqlXXX -h to see how to use it (make sure the version of mysql you want to link to is the same as the release description).
3. Execute the binary ask file and the SQL script will be output in the command line below.





### example



#### input

./sqlDiff_mysql8.0.32 -sourceDB "root:root@tcp(localhost:3306)/db_cs_account1.2.0" -targetDB "root:root@tcp(localhost:3306)/db_cs_account"

#### output

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

