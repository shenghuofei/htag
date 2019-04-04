## htag机器tag
### 功能说明
* 添加机器
* 删除机器
* 添加tag
* 删除tag
* 机器添加tag（一个或多个，tag格式为key=value,如：idc=dx）
* 机器删除tag（部分或全部）
* 机器修改tag
* 根据机器名查机器的tag
* 根据tag查机器列表
  1. 支持与(,)或(|)查询，逗号(,)表示与，管道(|)表示或(如:idc=dx,env=test|env=dev表示查询idc是dx且环境为test或者环境为dev的机器）
  2. 正常情况下与(,)或(|)的优先级一样，按从左向右查询；支持用小括号提高优先级（如:idc=lf|(idc=dx,env=prod)将先计算括号里的内容，然后作为整体与idc=lf做或运算）
  3. 实例(**注意后两个tag表达式只是优先级不一样**)：
     假如有h1,h2,h3,h4四台机器，其tag信息如下：
     

     | host | tag |
     | --- | --- |
     | h1 | idc=dx |
     | h1 | env=test |
     | h2 | idc=lf |
     | h2 | env=dev |
     | h3 | idc=dx |
     | h3 | evn=prod |
     | h4 | idc=lf |
     | h4 | env=dev |  



     * **idc=dx,env=test|env=dev** 查询结果为：**h1,h2,h4** ; 查询步骤：
       1. idc=dx的有h1,h3; **result=h1,h3**
       2. 并且所有env=test的有h1; h1,h3与h1取交集 **result=h1**
       3. 或者所有env=dev的有h2,h4; h1与h2,h4取并集 **result=h1,h2,h4**
     * **idc=lf|idc=dx,env=prod** 查询结果为：**h3** ; 查询步骤：
       1. idc=lf的有h2,h4; **result=h2,h4**
       2. 或者所有idc=dx的有h1,h3; h2,h4与h1,h3取并集 **result=h1,h2,h3,h4**
       3. 并且所有env=prod的有h3; h1,h2,h3,h4与h3取交集 **result=h3**
     * **idc=lf|(idc=dx,env=prod)** 查询结果为：**h2,h3,h4** ; 查询步骤：
       1. idc=lf的有h2,h4; **result=h2,h4**
       2. 有括号，所以先计算括号里的内容
          1. idc=dx的有h1,h3; **tmpresult=h1,h3**
          2. 并且所有env=prod的有h3; h1,h3与h3取交集 **tmpresult=h3**
       3. 最后result和tmpresult执行或运算；h2,h4与h3取并集 **result=h2,h3,h4**

### 最简化的sql表结构
* host表：
```
Table: host
Create Table: CREATE TABLE `host` (
`id` int(11) NOT NULL AUTO_INCREMENT,
`name` varchar(255) NOT NULL,
PRIMARY KEY (`id`),
UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8
```
* tag表：
```
Table: tag
Create Table: CREATE TABLE `tag` (
`id` int(11) NOT NULL AUTO_INCREMENT,
`name` varchar(255) NOT NULL,
PRIMARY KEY (`id`),
UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8
```
* hosttag表：
```
Table: hosttag
Create Table: CREATE TABLE `hosttag` (
`id` int(11) NOT NULL AUTO_INCREMENT,
`hid` int(11) DEFAULT NULL,
`tid` int(11) DEFAULT NULL,
PRIMARY KEY (`id`),
UNIQUE KEY `hosttag_uq` (`hid`,`tid`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8
```
### 主要sql语句
#### 机器添加tag
1. 机器添加tag(ignore 如果存在则忽略,不加的话会报Duplicate错),**添加一个可以用这个语句**:
`insert ignore into hosttag(hid,tid) select h.id,t.id from host h,tag t where h.name='h1' and t.name='idc=dx';`
2. 机器添加tag(ignore 如果存在则忽略,不加的话会报Duplicate错),**添加多个可以用这个语句**:
`insert ignore into hosttag(hid,tid) select h.id,t.id from host h,tag t where h.name='h1' and t.name in ('idc=dx','env=test');`
#### tag表中增加tag值
1. 添加tag值(ignore 如果存在则忽略,不加的话会报Duplicate错),**添加一个可以用这个语句**:
`insert ignore into tag(name) select 'idc=lf' from dual;`
2. 添加tag值(ignore 如果存在则忽略,不加的话会报Duplicate错),**添加多个可以用这个语句**:
`insert ignore into tag(name) values ('idc=rz'),('idc=gh');`
#### 根据tag查机器
1. 查询tagname='idc=dx'的所有机器
`select h.name from host h join hosttag ht on h.id=ht.hid join tag t on t.id=ht.tid where t.name='idc=dx';`
#### 根据机器查tag
1. 查询hostname='h1'的机器的所有tag
`select t.name from tag t join hosttag ht on t.id=ht.tid join host h on ht.hid=h.id where h.name='h1';`
#### 删除机器tag
1. 删除机器h1的某些tag:
`delete from hosttag where hid in (select id from host where name='h1') and tid in (select id from tag where name in ('idc=dx','env=test'));`
2. 删除机器h1的所有tag:
`delete from hosttag where hid = (select id from host where name='h1');`
#### 修改机器tag（分两步做，增加新的删除老的）:
1. 先添加新的tag:
`insert ignore into hosttag(hid,tid) select h.id,t.id from host h,tag t where h.name='h1' and t.name in ('idc=dx','env=test');`
2. 再删除老的tag:
`delete from hosttag where hid in (select id from host where name='h1') and tid in (select id from tag where name not in ('idc=dx','env=test'));`
### 编译
* 本地编译执行
  1. clone源码到`$GOPATH`下
  2. `cd $GOPATH/htag && go get ./...`
  3. `go build -o htag`
* 在mac上进行交叉编译(前2步跟本地编译一样):
  * mac上编译for linux:`GOOS=linux GOARCH=amd64 go build -o htag`
  * mac上编译for win:`GOOS=windows GOARCH=amd64 go build -o htag`
* 在linux上进行交叉编译(前2步跟本地编译一样):
  * linux上编译for mac:`GOOS=darwin GOARCH=amd64 go build -o htag`
  * linux上编译for win:`GOOS=windows GOARCH=amd64 go build -o htag`
### 使用说明(*如果参数里包含特殊字符一定要加引号,如'|'*)
* -action参数指定要进行操作的类型，包括如下actions:
  * addtag
  **添加tag:** 需要`-tlist`参数提供需要添加的tag列表,多个tag用逗号(,)分割, `htag -action addtag -tlist 'tag1,tag2'`
  * deltag
  **删除tag:** 需要`-tlist`参数提供需要删除的tag列表,多个tag用逗号(,)分割,`htag -action deltag -tlist 'tag1,tag2'`
  * addhost
  **添加机器:** 需要`-hlist`参数提供需要添加的机器列表,多个机器用逗号(,)分割,`htag -action addhost -hlist 'h1,h2'`
  * delhost
  **删除机器:** 需要`-hlist`参数提供需要删除的机器列表,多个机器用逗号(,)分割,`htag -action delhost -hlist 'h1,h2'`
  * gettaghost
  **根据tag查询机器列表:** 需要-t参数提供tag表达式,表达式说明见“功能说明”,***tag表达式一定要加引号***,`htag -action gettaghost -t 'idc=lf|(idc=dx,env=prod)'`
  * gethosttag
  **查询机器的tag:** 需要-h参数提供要查询的机器名,`htag -action gethosttag -h 'h1'`
  * addhosttag
  **给机器添加tag:** 需要-h参数提供要添加tag的机器名及-tlist参数提供机器要添加的tag列表(可以是一个,多个用逗号分隔),`htag -action addhosttag -h 'h1' -tlist 'tag1,tag2'`
  * delhostsometag
  **删除机器的部分tag:** 需要-h参数提供要删除tag的机器名及-tlist参数提供机器要删除的tag列表(可以是一个,多个用逗号分隔),`htag -action delhostsometag -h 'h1' -tlist 'tag1,tag2'`
  * delhostalltag
  **删除机器的所有tag:** 需要-h参数提供要删除tag的机器名,`htag -action delhostalltag -h 'h1'`
  * updatehosttag
  **修改机器的tag:** 需要-h参数提供要修改tag的机器名及-tlist参数提供机器的目标tag列表(可以是一个,多个用逗号分隔),`htag -action updatehosttag -h 'h1' -tlist 'tag1,tag2'`
* 退出状态码说明:
  *  _ = iota                   0是正常退出的状态码,跳过
  *  DBCONNERR                  db连接失败(1)
  *  ARGVERR                    参数错误(2)
  *  TAGERR                     tag格式错误(3)
  *  TAGNOTEXISTERR             机器添加tag时,要使用的tag不存在错误(4)
  *  CHKTAGERR                  检查tag格式时错误(5)
  *  EXECSQLERR                 sql执行失败(6)
  *  SQLSCANERR                 获取sql查询结果错误(7)
  *  SQLTXERR                   sql事物提交失败(8)
  *  TAGUSEDERR                 删除tag时,tag仍有机器关联错误,此错误目前不会退出,只跳过此tag的删除;此状态码暂时未用(9)
