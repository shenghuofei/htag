package main

import (
	"fmt"
	"strings"
)

// 查询机器的tag
func getHostTag(host string) {
	checkHostArgv(host)
	sql := fmt.Sprintf("SELECT t.name FROM %s t JOIN %s ht ON t.id=ht.tid JOIN %s h ON ht.hid=h.id WHERE h.name=?", TagTable, HostTagTable, HostTable)
	// rows, err := db.Query("SELECT t.name FROM tag t JOIN hosttag ht ON t.id=ht.tid JOIN host h ON ht.hid=h.id WHERE h.name=?", host)
	rows, err := db.Query(sql, host)
	defer rows.Close()
	checkErr(err, "selectHostTag exec sql failed", EXECSQLERR)
	tags := []string{}
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		checkErr(err, "selectHostTag scan data failed", SQLSCANERR)
		tags = append(tags, name)
	}
	fmt.Println(host, " : ", strings.Join(tags, ","))
}

// 给机器添加tag，即添加tag和host的关联关系
func hostAddTag(host, tags string) {
	tag_list := checkTagListArgv(tags)
	checkHostArgv(host)
	/*
	   // 使用原始sql命令操作
	   forSqlTagstr :=  ""
	   for _, tag := range tag_list {
	       if forSqlTagstr == "" {
	           forSqlTagstr = fmt.Sprintf("'%s'", tag)
	       } else {
	           forSqlTagstr = fmt.Sprintf("%s,'%s'", forSqlTagstr, tag)
	       }
	   }
	   sqlStr := fmt.Sprintf("INSERT IGNORE INTO %s(hid,tid) SELECT h.id,t.id FROM %s h,%s t WHERE h.name='%s' and t.name in (%s)", HostTagTable, HostTable, TagTable, host, forSqlTagstr)
	   fmt.Println(sqlStr)
	   _, err := db.Exec(sqlStr)
	   checkErr(err, "hostAddTag tag failed", EXECSQLERR)
	*/
	// 使用golang sql包的事物
	tx, err := db.Begin()
	sql := fmt.Sprintf("SELECT id FROM %s WHERE name=?", TagTable)
	txsql := fmt.Sprintf("INSERT IGNORE INTO %s(hid,tid) SELECT h.id,t.id FROM %s h,%s t WHERE h.name=? AND t.name=?", HostTagTable, HostTable, TagTable)
	for _, tag := range tag_list {
		// 先检查tag是否存在，要在插入事物前检查，所以不能用同一个事物
		// rows, err := db.Query("SELECT id FROM tag WHERE name=?", tag)
		rows, err := db.Query(sql, tag)
		defer rows.Close()
		checkErr(err, "hostAddTag check tag exist or not failed", EXECSQLERR)
		if !rows.Next() {
			printErr(tag+" tag is not exist, none tag has been added", TAGNOTEXISTERR)
		}
		// 每次循环用的都是tx内部的连接，没有新建连接，效率高
		// tx.Exec("INSERT IGNORE INTO hosttag(hid,tid) SELECT h.id,t.id FROM host h,tag t WHERE h.name=? AND t.name=?", host, tag)
		tx.Exec(txsql, host, tag)
	}
	err = tx.Commit()
	checkErr(err, "hostAddTag commit transaction failed", SQLTXERR)
	fmt.Println(host, "Add Tag Done")
}

// 删除机器的部分tag
func delHostSomeTag(host, tags string) {
	tag_list := checkTagListArgv(tags)
	checkHostArgv(host)
	tx, err := db.Begin()
	txsql := fmt.Sprintf("DELETE FROM %s WHERE hid = (SELECT id FROM %s WHERE name=?) AND tid = (SELECT id FROM %s WHERE name=?)", HostTagTable, HostTable, TagTable)
	for _, tag := range tag_list {
		// 每次循环用的都是tx内部的连接，没有新建连接，效率高
		// tx.Exec("DELETE FROM hosttag WHERE hid = (SELECT id FROM host WHERE name=?) AND tid = (SELECT id FROM tag WHERE name=?)", host, tag)
		tx.Exec(txsql, host, tag)
	}
	err = tx.Commit()
	checkErr(err, "delHostSomeTag commit transaction failed", SQLTXERR)
	fmt.Println(host, "Delete Tag Done")
}

// 删除机器的所有tag
func delHostTag(host string) {
	checkHostArgv(host)
	/*
	   // use db.Exec()
	   _, err := db.Exec("DELETE FROM hosttag WHERE hid = (SELECT id FROM host WHERE name=?)", host)
	   checkErr(err, "delHostTag execute sql failed", EXECSQLERR)
	*/

	// use db.Prepare()
	sql := fmt.Sprintf("DELETE FROM %s WHERE hid = (SELECT id FROM %s WHERE name=?)", HostTagTable, HostTable)
	// stm, _ := db.Prepare("DELETE FROM hosttag WHERE hid = (SELECT id FROM host WHERE name=?)")
	stm, _ := db.Prepare(sql)
	defer stm.Close()
	_, err := stm.Exec(host)
	checkErr(err, "delHostTag execute sql failed", EXECSQLERR)
	fmt.Println(host, "Delete All Tag Done")
}

// 更新机器的tag
func updateHostTag(host, tags string) {
	tag_list := checkTagListArgv(tags)
	checkHostArgv(host)
	// 在一个事物中，先删除老tag再添加新tag，如果不是事物最好先添加再删除
	tx, err := db.Begin()
	// 删除所有老tag
	txsql := fmt.Sprintf("DELETE FROM %s WHERE hid = (SELECT id FROM %s WHERE name=?)", HostTagTable, HostTable)
	// tx.Exec("DELETE FROM hosttag WHERE hid = (SELECT id FROM host WHERE name=?)", host)
	tx.Exec(txsql, host)
	// 添加新tag
	txsql = fmt.Sprintf("INSERT IGNORE INTO %s(hid,tid) SELECT h.id,t.id FROM %s h,%s t WHERE h.name=? AND t.name=?", HostTagTable, HostTable, TagTable)
	for _, tag := range tag_list {
		// tx.Exec("INSERT IGNORE INTO hosttag(hid,tid) SELECT h.id,t.id FROM host h,tag t WHERE h.name=? AND t.name=?", host, tag)
		tx.Exec(txsql, host, tag)
	}

	err = tx.Commit()
	checkErr(err, "updateHostTag commit transaction failed", SQLTXERR)
	fmt.Println(host, "Update Tag Done")
}
