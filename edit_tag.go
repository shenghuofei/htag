package main

import (
    "fmt"
)

// 向tag表中增加tag，多个tag用逗号(,)分割
func addTag(tags string) {
    tag_list := checkTagListArgv(tags)
    // 先检查tag格式是否符合要求
    checkTag(tags)
    tx, err := db.Begin()
    for _, tag := range tag_list {
        // 每次循环用的都是tx内部的连接，没有新建连接，效率高
        tx.Exec("INSERT IGNORE INTO tag(name) VALUES(?)", tag)
    }
    err = tx.Commit()
    checkErr(err, "hostAddTag commit transaction failed", SQLTXERR)
    fmt.Println("ADD Tag Done")
}


// 从tag表中删除tag，多个tag用逗号(,)分割
func delTag(tags string) {
    tag_list := checkTagListArgv(tags)
    tx, err := db.Begin()
    for _, tag := range tag_list {
        //先检查tag是否有机器关联
        rows, err := db.Query("SELECT ht.id FROM hosttag ht JOIN tag t ON ht.tid = t.id WHERE t.name=?", tag)
        checkErr(err, "delTag check tag is used or not failed", EXECSQLERR)
        if rows.Next() {
            /* 
            // 若有使用中的tag则退出的话如下:(由于事物未提交所以一个tag也不会删除)
            msg := fmt.Sprintf("%s is used now, can not delete", tag)
            printErr(msg, TAGUSEDERR)
            */

            // 跳过使用中的tag
            fmt.Println(tag, "is used now,can not delete")
            continue
        }
        // 每次循环用的都是tx内部的连接，没有新建连接，效率高
        tx.Exec("DELETE FROM tag WHERE name=?", tag)
    }
    err = tx.Commit()
    checkErr(err, "delTag commit transaction failed", SQLTXERR)
    fmt.Println("DELETE Tag Done")
}