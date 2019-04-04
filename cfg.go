package main

var (
	DBdsn string = "dbuser:dbpasswd@tcp(dbhost:dbport)/databasename?charset=utf8" // DB连接配置
)

// 错误退出码
const (
	_              = iota // 0是正常退出的状态码,跳过
	DBCONNERR             // db连接失败
	ARGVERR               // 参数错误
	TAGERR                // tag格式错误
	TAGNOTEXISTERR        // 机器添加tag时,要使用的tag不存在错误
	CHKTAGERR             // 检查tag格式时错误
	EXECSQLERR            // sql执行失败
	SQLSCANERR            // 获取sql查询结果错误
	SQLTXERR              // sql事物提交失败
	TAGUSEDERR            // 删除tag时,tag仍有机器关联错误,此错误目前不会退出,只跳过此tag的删除;此状态码暂时未用
)
