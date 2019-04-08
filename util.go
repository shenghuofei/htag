package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// 检查错误，任何报错均退出
func checkErr(err error, msg string, code int) {
	if err != nil {
		fmt.Println("ERROR:", err, msg)
		os.Exit(code)
	}
}

// 检查tag的格式是否符合要求
func checkTag(tag string) {
	pat := `^([\w-]+=[\w-]+,)*([\w-]+=[\w-]+)$`
	// pat := `^([\w-]+=[\w-]+)$`
	match, err := regexp.MatchString(pat, tag)
	checkErr(err, "check tag pattern fail", CHKTAGERR)
	if !match {
		printErr("tag format must be key=value,key value can include [a-zA-Z0-9_-] only", TAGERR)
	}
}

// 输出错误信息并退出
func printErr(msg string, code int) {
	fmt.Println("ERROR:", msg)
	os.Exit(code)
}

// 检查命令行host参数
func checkHostArgv(host string) {
	if host == "" {
		printErr("host argv is null, please add -h argv", ARGVERR)
	}
}

// 检查命令行tag参数
func checkTagArgv(tag string) {
	if tag == "" {
		printErr("tag argv is null, please add -t argv", ARGVERR)
	}
}

// 检查命令行tlist参数
func checkTagListArgv(tags string) []string {
	// 注意: 空字符串Split后返回的是长度为1的Slice,Slice中的元素是空字符串
	if tags == "" {
		printErr("tag list argv is null, please add -tlist argv", ARGVERR)
	}
	tag_list := strings.Split(tags, ",")
	return tag_list
}

// 检查命令行hlist参数
func checkHostListArgv(hosts string) []string {
	// 注意: 空字符串Split后返回的是长度为1的Slice,Slice中的元素是空字符串
	if hosts == "" {
		printErr("host list argv is null, please add -hlist argv", ARGVERR)
	}
	host_list := strings.Split(hosts, ",")
	return host_list
}

// 机器求并集
func hostOr(m1, m2 map[string]bool) map[string]bool {
	// 为提高性能，始终只遍历小的
	res := m1
	// fmt.Println("or",m1,m2)
	if len(m1) < len(m2) {
		// res始终初始化为大的集合，然后判断小集合，把小集合中有而大集合中没有的元素加入res中
		res = m2
		for key, value := range m1 {
			// fmt.Println("m1",key)
			if ok, _ := m2[key]; !ok {
				res[key] = value
			}
		}
	} else {
		for key, value := range m2 {
			// fmt.Println("m2",key)
			if ok, _ := m1[key]; !ok {
				res[key] = value
			}
		}
	}
	return res
}

// 机器求交集
func hostAnd(m1, m2 map[string]bool) map[string]bool {
	// 为提高性能，始终只遍历小的
	// fmt.Println("and",m1,m2)
	res := map[string]bool{}
	if len(m1) < len(m2) {
		for key, value := range m1 {
			// fmt.Println("m1",key)
			if ok, _ := m2[key]; ok {
				res[key] = value
			}
		}
	} else {
		for key, value := range m2 {
			// fmt.Println("m2",key)
			if ok, _ := m1[key]; ok {
				res[key] = value
			}
		}
	}
	return res
}
