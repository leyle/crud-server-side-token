# server side token

[TOC]

server side token(SST) 是一个给微服务内部互相之间进行权限验证的库。其目的是给 client 签发 token，在 server 内部对 token 进行验证，不依赖外部的验证方法或 api。

主要的功能有：

- create token
- verify token
- revoke token
- cli commands to process token



---

## 用法

sst 的用法分为两部分：

- 其他程序集成这个 sst library，在这些程序中使用 verifyToken 等相关函数进行 token 的验证
- 使用 sstcli 这个命令，进行 token 的生成与吊销等操作



## 其他程序集成这个 sst library

集成方法的步骤大概是：

1. 提供必要的参数，初始化 sst option
2. 调用 sst 暴露出来的方法



下面是一个 golang 例子

```go
// 首先初始化配置 sst
// 提供两个参数
// 分别是 aes key、zerolog 的 logger

// 初始化
	aesKey := "somekeyvalue"
	logger := logmiddleware.GetLogger(logmiddleware.LogTargetConsole)
	sst, err := sstapp.NewSSTokenOption(aesKey, &logger)
	if err != nil {
		fmt.Println("create server side token object failed")
		fmt.Println(err)
		os.Exit(1)
	}

// 使用

	token, err := sst.GenerateToken("useridvalue")
```



更多的完整例子，可以看代码中的 `samples/`目录。



---

## 使用 cli command 来操作 token

sst 提供了一个 cli 程序来处理 token 的签发、验证、revoke 等操作。

main 程序代码在 `cmd/sstcli/`目录下，如果临时使用，直接 `go build` 即可生成一个 `sstcli`的可执行程序。

使用 `./sstcli -h` 即可查看支持的命令及用法。

sst 使用了 `sqlite3`用来存储被 revoke 的 token 数据，这个文件的目录是 `$HOME/.config/sst/sst.db`。



### create random aes key

```shell
# 为了方便用户使用这套库，sstcli 程序同时提供了一个生成 aes key 的功能
# 输入的参数是 aes key 的长度，一般至少为 20 个值
# 内部实现中，拿到了 aes key，还会进行 md5 编码，转换为 32 个字符的长度使用。

./sstcli -createAesKey 20

# 可能的输出如下所示
# 在 `aesKey=` 后面的内容即为 random key 值
# 此例中为 ;d$d,ty8UZ<a:TA6$uVN
2022-07-07T11:25:42+08:00 INF main.go:63 > Generate aesKey succeeded aesKey=;d$d,ty8UZ<a:TA6$uVN

```



---

### create token

```shell
# 创建 token 需要两个参数
# 1. 给谁创建 token，无论是其他 service，或者某个具体的 user，此处统一命名为 userId
# 2. 签发 token 时使用的 aes 加密的 key 值

# 在输入中，为了安全，aes key 不在 commond line 中输入，而是使用 stdin 输入的方式

# 例子
./sstcli -createToken cdiservice

# 执行上面的命令后，会有个提示输入 aes key 的提示，程序会等待用户输入
# 此时，可以输入一个 aes key 值，比如从上面的 create random aes key 处得到的一个 key 值
# 需要特别注意的是，整个程序的生命周期中，一旦开始签发 token 后，aes key 应该妥善保存，并且不要变化（除非泄漏了），否则会造成已签发的 token 全部失效

# 假设此时输入 ;d$d,ty8UZ<a:TA6$uVN
# 如果整个程序没有问题，那么就会有大概如下的输出
# 其中 `token=` 后面的值即为签发的 token 值
./sstcli -createToken cdiservice
2022-07-07T11:29:30+08:00 INF main.go:140 > input aes key:
;d$d,ty8UZ<a:TA6$uVN
2022-07-07T11:31:55+08:00 DBG ../../sstapp/db.go:35 > get sqlite db ok, db name[/tmp/sstcli.tmp.db]
2022-07-07T11:31:55+08:00 DBG ../../sstapp/db.go:48 > create sqlite3 table, affected rows[0]
2022-07-07T11:31:55+08:00 INF ../../sstapp/db.go:117 > load revoke list from sqlite succeed, total num[1]
2022-07-07T11:31:55+08:00 TRC ../../sstapp/service.go:16 > GenerateToken succeed token=2_w9YiybBgjK2EeFo44b4seqKa2CHvbcr8hFTN2CXFbvcLwVOO5R7qePNEU3VY4t userId=cdiservice
2022-07-07T11:31:55+08:00 INF ../../sstapp/service.go:17 > GenerateToken succeed userId=cdiservice
2022-07-07T11:31:55+08:00 INF main.go:79 > token=2_w9YiybBgjK2EeFo44b4seqKa2CHvbcr8hFTN2CXFbvcLwVOO5R7qePNEU3VY4t
```



---

### verify token

```shell
# 这个命令在终端中使用的很少，一般是集成此 sst library 的其他程序在自己的程序调用相关方法
# 此命令使用与 create token 类似
# 此处展示一下执行过程，不做详细说明

# valid token
./sstcli -verifyToken 2_w9YiybBgjK2EeFo44b4seqKa2CHvbcr8hFTN2CXFbvcLwVOO5R7qePNEU3VY4t
2022-07-07T11:36:19+08:00 INF main.go:140 > input aes key:
;d$d,ty8UZ<a:TA6$uVN
2022-07-07T11:36:29+08:00 DBG ../../sstapp/db.go:35 > get sqlite db ok, db name[/Users/hmac/.config/sstcli/sstcli.db]
2022-07-07T11:36:29+08:00 DBG ../../sstapp/db.go:48 > create sqlite3 table, affected rows[0]
2022-07-07T11:36:29+08:00 INF ../../sstapp/db.go:117 > load revoke list from sqlite succeed, total num[1]
2022-07-07T11:36:29+08:00 DBG ../../sstapp/service.go:46 > verify token, decode succeed t=1657164715 userId=cdiservice
2022-07-07T11:36:29+08:00 INF main.go:98 > token is valid token=2_w9YiybBgjK2EeFo44b4seqKa2CHvbcr8hFTN2CXFbvcLwVOO5R7qePNEU3VY4t userId=cdiservice

# invalid token(aes key invalid)
./sstcli -verifyToken 2_w9YiybBgjK2EeFo44b4seqKa2CHvbcr8hFTN2CXFbvcLwVOO5R7qePNEU3VY4t
2022-07-07T11:36:57+08:00 INF main.go:140 > input aes key:
abc
2022-07-07T11:36:59+08:00 DBG ../../sstapp/db.go:35 > get sqlite db ok, db name[/Users/hmac/.config/sstcli/sstcli.db]
2022-07-07T11:36:59+08:00 DBG ../../sstapp/db.go:48 > create sqlite3 table, affected rows[0]
2022-07-07T11:36:59+08:00 INF ../../sstapp/db.go:117 > load revoke list from sqlite succeed, total num[1]
2022-07-07T11:36:59+08:00 WRN ../../sstapp/service.go:36 > decrypt aes token failed error="unpad error. This could happen when incorrect encryption key is used" token=2_w9YiybBgjK2EeFo44b4seqKa2CHvbcr8hFTN2CXFbvcLwVOO5R7qePNEU3VY4t
2022-07-07T11:36:59+08:00 WRN main.go:94 > invalid token, invalid token format, maybe wrong token, or maybe wrong aes key token=2_w9YiybBgjK2EeFo44b4seqKa2CHvbcr8hFTN2CXFbvcLwVOO5R7qePNEU3VY4t

# invalid token(token itself isn't valid)
./sstcli -verifyToken invalidtoken
2022-07-07T11:37:33+08:00 INF main.go:140 > input aes key:
;d$d,ty8UZ<a:TA6$uVN
2022-07-07T11:37:38+08:00 DBG ../../sstapp/db.go:35 > get sqlite db ok, db name[/Users/hmac/.config/sstcli/sstcli.db]
2022-07-07T11:37:38+08:00 DBG ../../sstapp/db.go:48 > create sqlite3 table, affected rows[0]
2022-07-07T11:37:38+08:00 INF ../../sstapp/db.go:117 > load revoke list from sqlite succeed, total num[1]
2022-07-07T11:37:38+08:00 WRN ../../sstapp/service.go:36 > decrypt aes token failed error="blocksize must be multipe of decoded message length" token=invalidtoken
2022-07-07T11:37:38+08:00 WRN main.go:94 > invalid token, invalid token format, maybe wrong token, or maybe wrong aes key token=invalidtoken
```



---

### revoke token

```shell
# 与 create token 等类似
# 需要注意的是，当 revoke 一个 valid token 时，会把此 token 值写入 sqlite3 数据库中
# 所以，需要特别注意此数据库的路径
# 默认的数据路径是 /tmp/
./sstcli -revokeToken 2_w9YiybBgjK2EeFo44b4seqKa2CHvbcr8hFTN2CXFbvcLwVOO5R7qePNEU3VY4t
2022-07-07T11:38:55+08:00 INF main.go:140 > input aes key:
;d$d,ty8UZ<a:TA6$uVN
2022-07-07T11:39:00+08:00 DBG ../../sstapp/db.go:35 > get sqlite db ok, db name[/Users/hmac/.config/sstcli/sstcli.db]
2022-07-07T11:39:00+08:00 DBG ../../sstapp/db.go:48 > create sqlite3 table, affected rows[0]
2022-07-07T11:39:00+08:00 INF ../../sstapp/db.go:117 > load revoke list from sqlite succeed, total num[1]
2022-07-07T11:39:00+08:00 DBG ../../sstapp/service.go:46 > verify token, decode succeed t=1657164715 userId=cdiservice
2022-07-07T11:39:00+08:00 INF ../../sstapp/db.go:88 > save token into revoke list succeed token=2_w9YiybBgjK2EeFo44b4seqKa2CHvbcr8hFTN2CXFbvcLwVOO5R7qePNEU3VY4t
2022-07-07T11:39:00+08:00 INF main.go:120 > revoke token succeeded token=2_w9YiybBgjK2EeFo44b4seqKa2CHvbcr8hFTN2CXFbvcLwVOO5R7qePNEU3VY4t userId=

```

