package main

import (
	"fmt"
	db "github.com/liuzhiyi/go-db"
)

func createData() {
	transaction := db.F.GetConnect("write").BeginTransaction()
	defer func() {
		err := recover()
		fmt.Println(err)
		transaction.Rollback()
	}()
	transaction.MustExec("DROP TABLE IF EXISTS `core_website`;")
	transaction.MustExec(" CREATE TABLE `core_website` ( " +
		"`website_id` smallint(5) unsigned NOT NULL auto_increment, " +
		"`code` varchar(32) NOT NULL default '', " +
		"`name` varchar(64) NOT NULL default '', " +
		"`sort_order` smallint(5) unsigned NOT NULL default '0', " +
		"`is_active` tinyint(1) unsigned NOT NULL default '0', " +
		"PRIMARY KEY  (`website_id`), " +
		"UNIQUE KEY `code` (`code`), " +
		"KEY `is_active` (`is_active`,`sort_order`) " +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Websites';")
	transaction.MustExec("insert  into `core_website`(`website_id`,`code`,`name`,`sort_order`,`is_active`) values (1,'default','Default',0,1),(2,'base','Main Website',0,1);")

	transaction.MustExec("DROP TABLE IF EXISTS `core_api`;")
	transaction.MustExec(" CREATE TABLE `core_api` ( " +
		"`api_id` smallint(5) unsigned NOT NULL auto_increment, " +
		"`website_id` smallint(5) unsigned NOT NULL, " +
		"`api_name` varchar(255) NOT NULL default '', " +
		"PRIMARY KEY  (`api_id`), " +
		"UNIQUE KEY (`api_id`) " +
		// "foreign key (`website_id`) references `core_website` (`website_id`) on delete cascade on update cascade" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Websites';")
	transaction.MustExec("insert  into `core_api`(`api_id`, `website_id`, `api_name`) values (1,1, 'li ming'),(2, 2, 'liu hua');")
	transaction.Commit()
}
