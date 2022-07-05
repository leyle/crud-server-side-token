# server side token

server side token(SST) 是一个给微服务内部互相之间进行权限验证的库。其目的是给 client 签发 token，在 server 内部对 token 进行验证，不依赖外部的验证方法或 api。

主要的功能有：

- create token
- verify token
- revoke token
- query token
- restful api



---

## api list

### create token

**request**

```json
{
    "client": "cdi-service"
}
```

**response**

```json
{
    "token": "f45de51cdef30991551e41e882dd7b5404799648a0a00753f44fc966e6153fc1"
}
```



---

### verify token

**request**

```shell
// set token in headers as X-Server-Side-Token:token
```

**response**

`valid`

```json
{
    "status": "VALID",
    "client": "cdi-service"
}
```

`invalid`

```json
{
    "status": "INVALID"
}
```



---

### revoke token

