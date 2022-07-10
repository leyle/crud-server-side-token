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

# 创建了此 key 后，可以保存起来，后续命令也可以使用，而其他程序集成此库时，如果已经用 cli 签发了 token，那么也要使用相同的 key 值。
```

**key.yaml** 例子

```yaml
key: "v}#0CEYuG%M8#c77HkUp"
```



---

### create token

```shell
# 创建 token 需要两个参数
# 1. 给谁创建 token，无论是其他 service，或者某个具体的 user，此处统一命名为 userId
# 2. 签发 token 时使用的 aes 加密的 key 值

# 因为签发需要使用到 key，为了避免在 shell history 中看到输入的 key 值，使用从文件中读取 key
# key.yaml 见上面

# 例子
./sstcli -secretFile /path/to/key.yaml -createToken cdiservice

# 需要特别注意的是，整个程序的生命周期中，一旦开始签发 token 后，aes key 应该妥善保存，并且不要变化（除非泄漏了），否则会造成已签发的 token 全部失效

# 如果整个程序没有问题，那么就会有大概如下的输出
# 其中 `token=` 后面的值以 SST- 开头的即为签发的 token 值
# SST-aTUKdDiJczOu/8vPVLoeA9rN5m7aEpWIhL1Ue4gIb28aX7EWhJrNTFr8P9U5tCMt/A==
./sstcli -secretFile /tmp/key.yaml -createToken cdiservice
2022-07-10T16:38:39+08:00 DBG ../../sstapp/db.go:35 > get sqlite db ok, db name[/Users/hmac/.config/sst/sst.db]
2022-07-10T16:38:39+08:00 DBG ../../sstapp/db.go:48 > create sqlite3 table, affected rows[0]
2022-07-10T16:38:39+08:00 INF ../../sstapp/db.go:117 > load revocation list from sqlite succeed, total num[7]
2022-07-10T16:38:39+08:00 TRC ../../sstapp/api.go:18 > GenerateToken succeed token=SST-1a2e5e90b045cd149017d6ea1ede7775bc6e1a578c671d7b7858c8970f7785d02eba574caa05f7cece692ace2a52014092 userId=cdiservice
2022-07-10T16:38:39+08:00 INF ../../sstapp/api.go:19 > GenerateToken succeed userId=cdiservice
2022-07-10T16:38:39+08:00 INF main.go:92 > token=SST-1a2e5e90b045cd149017d6ea1ede7775bc6e1a578c671d7b7858c8970f7785d02eba574caa05f7cece692ace2a52014092
```



---

### verify token

```shell
# 这个命令在终端中使用的很少，一般是集成此 sst library 的其他程序在自己的程序调用相关方法
# 此命令使用与 create token 类似
# 此处展示一下执行过程，不做详细说明

# valid token
./sstcli -secretFile /tmp/key.yaml -verifyToken SST-1a2e5e90b045cd149017d6ea1ede7775bc6e1a578c671d7b7858c8970f7785d02eba574caa05f7cece692ace2a52014092
2022-07-10T16:39:30+08:00 DBG ../../sstapp/db.go:35 > get sqlite db ok, db name[/Users/hmac/.config/sst/sst.db]
2022-07-10T16:39:30+08:00 DBG ../../sstapp/db.go:48 > create sqlite3 table, affected rows[0]
2022-07-10T16:39:30+08:00 INF ../../sstapp/db.go:117 > load revocation list from sqlite succeed, total num[7]
2022-07-10T16:39:30+08:00 DBG ../../sstapp/api.go:58 > verify token, decode succeed t=1657442319 userId=cdiservice
2022-07-10T16:39:30+08:00 INF main.go:110 > token is valid token=SST-1a2e5e90b045cd149017d6ea1ede7775bc6e1a578c671d7b7858c8970f7785d02eba574caa05f7cece692ace2a52014092 userId=cdiservice

# invalid token(aes key invalid)
./sstcli -secretFile /tmp/key.yaml -verifyToken SST-1a2e5e90b045cd149017d6ea1ede7775bc6e1a578c671d7b7858c8970f7785d02eba574caa05f7cece692ace2a52014092
2022-07-10T16:39:49+08:00 DBG ../../sstapp/db.go:35 > get sqlite db ok, db name[/Users/hmac/.config/sst/sst.db]
2022-07-10T16:39:49+08:00 DBG ../../sstapp/db.go:48 > create sqlite3 table, affected rows[0]
2022-07-10T16:39:49+08:00 INF ../../sstapp/db.go:117 > load revocation list from sqlite succeed, total num[7]
2022-07-10T16:39:49+08:00 WRN ../../sstapp/model.go:137 > decrpyt cipher text failed error="cipher: message authentication failed"
2022-07-10T16:39:49+08:00 WRN ../../sstapp/api.go:48 > decrypt aes token failed error="cipher: message authentication failed" token=SST-1a2e5e90b045cd149017d6ea1ede7775bc6e1a578c671d7b7858c8970f7785d02eba574caa05f7cece692ace2a52014092
2022-07-10T16:39:49+08:00 WRN main.go:106 > invalid token, invalid token format, maybe wrong token, or maybe wrong aes key token=SST-1a2e5e90b045cd149017d6ea1ede7775bc6e1a578c671d7b7858c8970f7785d02eba574caa05f7cece692ace2a52014092

# invalid token(token itself isn't valid)
./sstcli -secretFile /tmp/key.yaml -verifyToken SST-1a2e5e90b045cd149017d6ea1ede7775bc6e1a578c671d7b7858c8970f7785d02eba574caa05f7cece692ace2a52014093
2022-07-10T16:40:09+08:00 DBG ../../sstapp/db.go:35 > get sqlite db ok, db name[/Users/hmac/.config/sst/sst.db]
2022-07-10T16:40:09+08:00 DBG ../../sstapp/db.go:48 > create sqlite3 table, affected rows[0]
2022-07-10T16:40:09+08:00 INF ../../sstapp/db.go:117 > load revocation list from sqlite succeed, total num[7]
2022-07-10T16:40:09+08:00 WRN ../../sstapp/model.go:137 > decrpyt cipher text failed error="cipher: message authentication failed"
2022-07-10T16:40:09+08:00 WRN ../../sstapp/api.go:48 > decrypt aes token failed error="cipher: message authentication failed" token=SST-1a2e5e90b045cd149017d6ea1ede7775bc6e1a578c671d7b7858c8970f7785d02eba574caa05f7cece692ace2a52014093
2022-07-10T16:40:09+08:00 WRN main.go:106 > invalid token, invalid token format, maybe wrong token, or maybe wrong aes key token=SST-1a2e5e90b045cd149017d6ea1ede7775bc6e1a578c671d7b7858c8970f7785d02eba574caa05f7cece692ace2a52014093
```



---

### revoke token

```shell
# 与 create token 等类似
# 需要注意的是，当 revoke 一个 valid token 时，会把此 token 值写入 sqlite3 数据库中
# 所以，需要特别注意此数据库的路径
# 默认的数据路径是 $HOME/.config/sst/sst.db
./sstcli -secretFile /tmp/key.yaml -revokeToken SST-1a2e5e90b045cd149017d6ea1ede7775bc6e1a578c671d7b7858c8970f7785d02eba574caa05f7cece692ace2a52014092
2022-07-10T16:40:43+08:00 DBG ../../sstapp/db.go:35 > get sqlite db ok, db name[/Users/hmac/.config/sst/sst.db]
2022-07-10T16:40:43+08:00 DBG ../../sstapp/db.go:48 > create sqlite3 table, affected rows[0]
2022-07-10T16:40:43+08:00 INF ../../sstapp/db.go:117 > load revocation list from sqlite succeed, total num[7]
2022-07-10T16:40:43+08:00 DBG ../../sstapp/api.go:58 > verify token, decode succeed t=1657442319 userId=cdiservice
2022-07-10T16:40:43+08:00 INF ../../sstapp/db.go:88 > save token into revocation list succeed token=SST-1a2e5e90b045cd149017d6ea1ede7775bc6e1a578c671d7b7858c8970f7785d02eba574caa05f7cece692ace2a52014092
2022-07-10T16:40:43+08:00 INF main.go:131 > revoke token succeeded token=SST-1a2e5e90b045cd149017d6ea1ede7775bc6e1a578c671d7b7858c8970f7785d02eba574caa05f7cece692ace2a52014092 userId=
```

