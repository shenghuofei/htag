package main

import (
	"fmt"
	// "os"
	"strings"
)

// 获取tag的优先级，返回去掉()后的tag列表和分割符列表，括号内的保持原样，()不支持括号嵌套
func getTagPrio(tags string) (tag_list, sep []string) {
	tags = strings.Trim(tags, " ")
	tags = strings.Trim(tags, "|")
	tags = strings.Trim(tags, ",")
	taglen := len(tags)
	tag := ""
	var i, k int
	for i = 0; i < taglen; i++ {
		if string(tags[i]) == "(" {
			if tag != "" {
				tag_list = append(tag_list, tag)
				tag = ""
			}
			for k = i + 1; k < taglen; k++ {
				if string(tags[k]) == ")" {
					if tag != "" {
						tag_list = append(tag_list, tag)
						tag = ""
					}
					i = k
					break
				} else if string(tags[k]) == " " {
					continue
				} else if string(tags[k]) == "(" {
					printErr("request tag format error,nesting is not support", TAGERR)
					// fmt.Println("ERROR: request args format error,nesting is not support")
					// os.Exit(1)
				} else {
					tag = tag + string(tags[k])
				}
			}
			if k == taglen {
				printErr("request tag format error,( and ) not pair", TAGERR)
				// fmt.Println("ERROR: request args format error,( and ) not pair")
				// os.Exit(1)
			}
		} else if string(tags[i]) == " " {
			continue
		} else if string(tags[i]) == "," || string(tags[i]) == "|" {
			sep = append(sep, string(tags[i]))
			if tag != "" {
				tag_list = append(tag_list, tag)
				tag = ""
			}
		} else {
			tag = tag + string(tags[i])
		}
	}
	if tag != "" {
		tag_list = append(tag_list, tag)
	}
	if len(tag_list) == 0 {
		printErr("tag format must be key=value and split with , or |", TAGERR)
		// fmt.Println("ERROR: tag format must be key=value and split with , or |")
		// os.Exit(1)
	}
	return
}

// 将复合tag拆分为单个k=v的tag,返回tag列表及分隔符列表
func getTagList(tags string) (tag_list, sep []string) {
	tags = strings.Trim(tags, " ")
	tags = strings.Trim(tags, "|")
	tags = strings.Trim(tags, ",")
	tag := ""
	for _, c := range tags {
		if string(c) == "," || string(c) == "|" {
			sep = append(sep, string(c))
			tag_list = append(tag_list, tag)
			tag = ""
		} else if string(c) == " " {
			continue
		} else {
			tag = tag + string(c)
		}
	}
	if tag != "" {
		tag_list = append(tag_list, tag)
	}
	if len(tag_list) == 0 {
		printErr("tag format must be key=value and split with , or |", TAGERR)
		// fmt.Println("ERROR: tag format must be key=value and split with , or |")
		// os.Exit(1)
	}
	for _, tag = range tag_list {
		checkTag(tag)
	}
	return
}

// 查找单个tag的所有实例
func getOneTagHost(tag string) map[string]bool {
	// fmt.Println("tag:", tag)
	// sql := fmt.Sprintf("select name from instances where id in (select ins_id from tags where tag='%s')", tag)
	rows, err := db.Query("SELECT h.name from host h JOIN hosttag ht ON h.id=ht.hid JOIN tag t ON t.id=ht.tid WHERE t.name=?", tag)
	checkErr(err, "getOneTagHost execute sql fail", EXECSQLERR)
	name := ""
	res := map[string]bool{}
	for rows.Next() {
		name = ""
		rows.Scan(&name)
		res[name] = true
	}
	// fmt.Println("list:",res)
	return res
}

// 查询优先级tag列表中一个tag(这个tag可能是复合tag)的所有实例，并按逻辑进行取舍
func getOnePrioTagHost(tags string) map[string]bool {
	// fmt.Println("getOnePrioTagHost", tags)
	tag_list, sep := getTagList(tags)
	tmp_res := getOneTagHost(tag_list[0])
	taglen := len(tag_list)
	for i := 1; i < taglen; i++ {
		res := getOneTagHost(tag_list[i])
		// fmt.Println(sep[i-1])
		if sep[i-1] == "," {
			tmp_res = hostAnd(tmp_res, res)
		} else if sep[i-1] == "|" {
			tmp_res = hostOr(tmp_res, res)
		} else {
			printErr("tag split just for , and |", TAGERR)
			// fmt.Println("ERROR: tag split just for , and |")
			// os.Exit(1)
		}
	}
	// fmt.Println("getOnePrioTagHost", tmp_res)
	return tmp_res
}

// 查询优先级tag列表中所有tag的所有实例，并按逻辑进行取舍
func getAllPrioTagHost(tags string) map[string]bool {
	prio_tag_list, prio_sep := getTagPrio(tags)
	// fmt.Println(prio_tag_list, prio_sep)
	tmp_res := getOnePrioTagHost(prio_tag_list[0])
	taglen := len(prio_tag_list)
	for i := 1; i < taglen; i++ {
		// fmt.Println("getAllPrioTagHost", prio_tag_list[i])
		res := getOnePrioTagHost(prio_tag_list[i])
		// fmt.Println("getAllPrioTagHost", prio_sep[i-1])
		if prio_sep[i-1] == "," {
			tmp_res = hostAnd(tmp_res, res)
		} else if prio_sep[i-1] == "|" {
			tmp_res = hostOr(tmp_res, res)
		} else {
			printErr("tag split just for , and |", TAGERR)
			// fmt.Println("ERROR: tag split just for , and |")
			// os.Exit(1)
		}
		// fmt.Println("getAllPrioTagHost", tmp_res)
	}
	return tmp_res
}

// 查询tag的所有实例并按逻辑取舍，然后输出结果
func getTagHost(tags string) {
	checkTagArgv(tags)
	res := getAllPrioTagHost(tags)
	for host, _ := range res {
		fmt.Println(host)
	}
	// fmt.Println(len(res), "results")
}
