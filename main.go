package main

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	db, _ = sql.Open("mysql", DBdsn)
	defer db.Close()
	// for connect pool
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	err := db.Ping()
	checkErr(err, "connect to db failed", DBCONNERR)

	/*
	   actions := map[string]bool {
	       "addtag": true,
	       "deltag": true,
	       "addhost": true,
	       "delhost": true,
	       "gettaghost": true,
	       "gethosttag": true,
	       "addhosttag": true,
	       "delhostsometag": true,
	       "delhostalltag": true,
	       "updatehosttag": true
	   }
	*/

	action := flag.String("action", "", "configuration action, action list: [addtag,deltag,addhost,delhost,gettaghost,gethosttag,addhosttag,delhostsometag,delhostalltag,updatehosttag]")
	tag := flag.String("t", "", "configuration tag expression for actions:[gettaghost], such as 'idc=dx|idc=lf'")
	host := flag.String("h", "", "configuration host name for actions:[gethosttag,addhosttag,delhostsometag,delhostalltag,updatehosttag], such as 'h1'")
	hostlist := flag.String("hlist", "", "configuration host list for actions:[addhost,delhost], split with ',', such as 'h1,h2'")
	taglist := flag.String("tlist", "", "configuration tag list for actions:[addtag,deltag,addhosttag,delhostsometag,updatehosttag], split with ',', such as 'idc=dx,idc=lf'")
	flag.Parse()
	if *action == "" {
		printErr("no action support", ARGVERR)
	}
	/*
	   if _, ok := actions[*action]; !ok {
	       fmt.Println("ERROR: do not support this action")
	       os.Exit(2)
	   }
	*/

	switch *action {
	case "addhost":
		addHost(*hostlist)
	case "delhost":
		delHost(*hostlist)
	case "addtag":
		addTag(*taglist)
	case "deltag":
		delTag(*taglist)
	case "gettaghost":
		getTagHost(*tag)
	case "gethosttag":
		getHostTag(*host)
	case "addhosttag":
		hostAddTag(*host, *taglist)
	case "delhostsometag":
		delHostSomeTag(*host, *taglist)
	case "delhostalltag":
		delHostTag(*host)
	case "updatehosttag":
		updateHostTag(*host, *taglist)
	default:
		printErr("do not support this action", ARGVERR)
	}
}
