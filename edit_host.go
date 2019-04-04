package main

import (
    "fmt"
)

// 向host表中增加host，多个host用逗号(,)分割
func addHost(hosts string) {
    host_list := checkHostListArgv(hosts)
    tx, err := db.Begin()
    for _, host := range host_list {
        // 每次循环用的都是tx内部的连接，没有新建连接，效率高
        tx.Exec("INSERT IGNORE INTO host(name) VALUES(?)", host)
    }
    err = tx.Commit()
    checkErr(err, "addHost commit transaction failed", SQLTXERR)
    fmt.Println("ADD Host Done")
}

// 从host表中删除host，多个host用逗号(,)分割
func delHost(hosts string) {
    host_list := checkHostListArgv(hosts)
    tx, err := db.Begin()
    for _, host := range host_list {
        // 先删除host关联的tag
        tx.Exec("DELETE FROM hosttag WHERE hid = (SELECT id FROM host WHERE name=?)", host)
        // 再删除host
        tx.Exec("DELETE FROM host WHERE name=?", host)
    }
    err = tx.Commit()
    checkErr(err, "delHost commit transaction failed", SQLTXERR)
    fmt.Println("DELETE Host Done")
}
