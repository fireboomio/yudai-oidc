# YuDai

取自古代身份验证令牌“鱼袋”

![鱼袋](https://static.aihanfu.net/uploadfile/2014/1128/20141128110743761.jpg)

## 功能介绍

我们初期参考了`Casdoor`的功能进行设计，目标是开发一款精简的符合`OIDC`协议的登录认证服务，将其作为`OpenAPI`数据源注册到[飞布 Fireboom](https://www.fireboom.cloud/zh)，支持手机短信登录和密码登录等。

> OIDC 是一种用于身份验证和授权的开放式协议。它建立在OAuth 2.0基础上，并为第三方应用程序提供了一种方便的方法来验证用户身份并获取用户信息，例如名称、邮件地址等。OIDC还支持单点登录（SSO），以便用户只需在一个地方登录，就可以访问多个应用程序。

> 在分布式架构中，我们的服务可能还要接入第三方服务商，比如使用微信、QQ登录，此时token认证商是三方的服务，这个时候就有可能需要我们提供token认证的方法给三方服务。cert通常代表着公私钥对中的私钥，用于对JWT进行签名，验证Token时使用公钥进行解密和验证。那么我们是需要把公钥暴露给三方服务的。在飞布控制台我们可以去按照一定的规范添加身份验证商。

## 如何使用

### 安装：

```sh
# Linux
curl -o ./yudai https://yudai-bin.fireboom.io/build-env/yudai-linux
# Mac
curl -o ./yudai https://yudai-bin.fireboom.io/build-env/yudai-mac
# Windows
curl -o ./yudai.exe https://yudai-bin.fireboom.io/build-env/yudai-windows.exe

chmod +x ./yudai
```

### 数据库初始化

请先准备好数据库，目前`YuDai`支持`mysql`和`postgres`，考虑到实际业务场景，我们暂时不打算支持`sqlite`数据库。

`YuDai`服务启动后会自动同步表信息到数据库，所以不需要初始化表结构。

### 配置

`YuDai`支持以环境变量的方式来配置，下面是所有支持的环境变量

- `YUDAI_PORT` 服务监听端口
- `YUDAI_DB_URL` 数据库连接字符串，`postgres`使用`postgres://user:password@host:port/dbname?sslmode=disable`，`mysql`使用`mysql://user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local`，更多参数请分别查看文档[postgres driver](https://pkg.go.dev/github.com/lib/pq)和[mysql driver](https://github.com/go-sql-driver/mysql)
- `YUDAI_DB_PREFIX` 数据库表前缀，可以用来区分多个服务或者和业务表区分开
- `YUDAI_PC_APP_ID` 微信pc扫码登录app_id
- `YUDAI_PC_APP_SECRET` 微信pc扫码登录app_secret
- `YUDAI_H5_APP_ID` 微信公众号h5登录app_id
- `YUDAI_H5_APP_SECRET` 微信公众号h5登录app_secret
- `YUDAI_APP_APP_ID` 微信app登录app_id
- `YUDAI_APP_APP_SECRET` 微信app登录app_secret
- `YUDAI_MINI_APP_ID` 微信小程序登录app_id
- `YUDAI_MINI_APP_SECRET` 微信小程序登录app_secret

### 短信配置

在provider表中配置短信的相关信息。

```sql
INSERT INTO "provider" ("owner", "name", "created_at", "type", "client_id", "client_secret", "sign_name", "template_code") VALUES ('your_name', 'provider_sms', '2023-01-17 01:22:33', 'Aliyun SMS', 'xxxxxxxxxxxxxxxx', 'xxxxxxxxxxxxxxxxxxxxxxxxxxx', 'app_name', 'temp_code');
```

### 密钥生成

`YuDai`默认会在`cert`目录下生成`token_jwt_key.key`和`token_jwt_key.pem`两个文件，用于生成JWT签名和验证，如果这2个文件不存在会自动生成，但每次重新生成会导致之前的jwt失效，因此在正式项目使用中可以通过 Volume 挂载的方式保持这2个文件不变。

### 启动

```sh
./yudai
```

## 如何和fireboom结合

[前往阅读](https://docs.fireboom.io/v/v1.0/ji-chu-ke-shi-hua-kai-fa/shen-fen-yan-zheng/yin-shi-mo-shi#oidc-pei-zhi)

### Fireboom如何配置（跟钩子、RBAC没关系）

将`YuDai`启动后，启动飞布控制台，在身份验证新建身份验证器。

参考官方文档：https://ansons-organization.gitbook.io/product-manual/kai-fa-wen-dang/yan-zheng-he-shou-quan/shen-fen-yan-zheng

根据`YuDai`提供的路由信息在飞布控制台添加并配置身份验证。

飞布根据 Issuer 生成服务发现地址，其中 [http://127.0.0.1:9825](http://127.0.0.1:10021/.well-known/openid-configuration)为`YuDai`服务地址，[http://127.0.0.1:9825/.well-known/openid-configuration](http://127.0.0.1:10021/.well-known/openid-configuration) 为服务发现地址，返回数据中定义了包括但不限于：

- 用户端点url：userinfo_endpoint
- jwks访问url：jwks_uri

飞布根据这些信息解析出`jwks_uri`后，就可以对请求进行`token`认证

![img](https://cos.ap-nanjing.myqcloud.com/test-1314985928/admin/image.png)

### Restful API

当需要在飞布中以`Restful API`的方式调用`YuDai`时，可以参考[swagger.json](docs/oidc.json)

### 补充：

> JWKS (JSON Web Key Set) 代表了用于身份验证和授权的一组JSON格式的Web密钥。它是OAuth 2.0和OpenID Connect等身份验证和授权协议中的一部分。JWKS通常用于在安全的方式下公开API的公钥和其他安全参数。
>
> JWKS包含一个包含一组密钥的JSON对象，每个密钥都有一个唯一的标识符（kid），并提供其算法（alg）、密钥类型（kty）和公钥（用于验证签名）或私钥（用于生成签名）等信息。通常，JWKS中的密钥是使用非对称加密算法生成的，如RSA或ECDSA。
>
> 在OAuth 2.0和OpenID Connect中，JWKS用于验证由身份提供者（如认证服务器）签名的访问令牌或ID令牌。客户端可以通过从指定的JWKS端点获取JWKS并验证签名，来确保令牌的身份验证和完整性。
>
> 总结来说，JWKS是一个存储公钥和其他安全参数的JSON对象，用于身份验证和授权协议中的安全验证和签名操作。

### 前端如何使用

前端请求登录接口，验证成功后会返回`access_token`、`refresh_token`和对应的过期时间。前端将`access_token`值、`access_token`过期时间、`refresh_token`值与`refresh_token`的过期时间存储到`localStorage`里。

当请求非白名单的接口时，从`localStorage`中取出`access_token`值与其过期时间.首先检验是否过期，若过期，则读取`refresh_token`与其对应的过期时间，若refreshToken没有过期，则携带`refresh_token`去请求响应的接口，获取新的`access_token`和新的`refresh_token`，并携带新的`access_token`去请求接口，若`refresh_token`也过期，则跳转到登录页面。

前端也可以不存储过期时间，全部由接口是否返回401状态码来判断`access_token`和`refresh_token`是否过期。

## Todo list

[ ] cookie模式支持
[x] pgsql支持
