# YuDai

取自古代身份验证令牌“鱼袋”

![鱼袋](https://static.aihanfu.net/uploadfile/2014/1128/20141128110743761.jpg)

# 1. 功能介绍

我们的演示服务基于 Casdoor 进行了精简，自研部署了一套符合 OIDC 协议的登录认证服务，将其作为 OpenAPI数据源注册到飞步，支持手机短信登录和密码登录

OIDC 是一种用于身份验证和授权的开放式协议。它建立在OAuth 2.0基础上，并为第三方应用程序提供了一种方便的方法来验证用户身份并获取用户信息，例如名称、邮件地址等。OIDC还支持单点登录（SSO），以便用户只需在一个地方登录，就可以访问多个应用程序。

在分布式架构中，我们的服务可能还要接入第三方服务商，比如使用微信、QQ登录，此时token认证商是三方的服务，这个时候就有可能需要我们提供token认证的方法给三方服务。cert通常代表着公私钥对中的私钥，用于对JWT进行签名，验证Token时使用公钥进行解密和验证。那么我们是需要把公钥暴露给三方服务的。在飞布控制台我们可以去按照一定的规范添加身份验证商。

# 2. 如何启动

## 2.1 安装：
curl -o ./yudai https://yudai-bin.fireboom.io/build/yudai-linux

## 2.2 数据库初始化

数据库脚本位于[oidc.sql](docs/oidc.sql), 请在Mysql数据库中执行。

## 2.3 配置文件

在conf/config.yaml中配置自己的数据库端口，用户名，密码及微信开放平台等信息[config.yaml](conf/config.yaml)

```Yaml
mysql:
  host: "localhost" // 数据库地址
  port: 3306 // 数据库端口
  user: "root" // 数据库用户名
  password: "123456" // 数据库密码
  dbname: "oidc" // 数据库名
  max_open_conns: 200 // 最大连接数
  max_idle_conns: 50 // 最大空闲连接数
wxlogin:
  h5:
    appid: "appid" // 微信公众号appid
    secret: "secret" // 微信公众号secret
  pc:
    appid: "appid" // 微信pc扫码登录appid
    secret: "secret" // 微信pc扫码登录secret
  app:
    appid: "appid" // 微信app登录appid
    secret: "secret" // 微信app登录secret
  mini:
    appid: "appid" // 微信小程序appid
    secret: "secret" // 微信小程序secret
```

## 2.4 短信配置

执行完sql之后，在provider表中配置短信的相关信息。

```Go
package objcet

type Provider struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`       // 创建者
	Name        string `xorm:"varchar(100) notnull pk unique" json:"name"` // 名称
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`            // 创建时间

	Type         string `xorm:"varchar(100)" json:"type"`          // 类型
	ClientId     string `xorm:"varchar(100)" json:"clientId"`      // 客户端ID
	ClientSecret string `xorm:"varchar(2000)" json:"clientSecret"` // 客户端秘钥
	SignName     string `xorm:"varchar(100)" json:"signName"`      // 签名名称
	TemplateCode string `xorm:"varchar(100)" json:"templateCode"`  // 模板code
}
```

在object/sms.go中会获取sms短信提供方的相关信息用于短信发送。

在api/verification.go会根据获取到的sms信息进行短信的发送。

## 2.5 架构图 https://docs.fireboom.io/v/v1.0/ji-chu-ke-shi-hua-kai-fa/shen-fen-yan-zheng/yin-shi-mo-shi#gong-zuo-yuan-li

# 3. 登录功能 调用方式参考swagger文档 [oidc.json](docs/oidc.json)

## 3.1 账户密码登录
```
POST /api/login
{
  "username": "admin",
  "password": "123456",
  "loginType": "password"
}
```

## 3.2 手机号验证码登录
```
POST /api/login
{
  "phone": "13800138000",
  "code": "123456",
  "loginType": "phone"
}
```

## 3.3 微信公众号h5登录
```
POST /api/login
{
  "code": "code",
  "loginType": "h5"
}
```

## 3.4 微信pc扫码登录
```
POST /api/login
{
  "code": "code",
  "loginType": "pc"
}
```

## 3.5 微信小程序登录
```
POST /api/login
{
  "code": "code",
  "loginType": "mini"
}
```

## 3.6 微信app登录
```
POST /api/login
{
  "code": "code",
  "loginType": "app"
}
```

## 3.7 密钥生成

在 `object.jwks.go` 中需要初始化 cert, 启动项目后在 `object.jwks.go` 程序初始化时会查找项目根目录下的两个文件（用于OIDC服务发现的private key 和 certificate），如果没有会在初始化时生成这两个文件。

# 4. 如何和fireboom结合 https://docs.fireboom.io/v/v1.0/ji-chu-ke-shi-hua-kai-fa/shen-fen-yan-zheng/yin-shi-mo-shi#oidc-pei-zhi

## 4.1 Fireboom如何配置（跟钩子、RBAC没关系）

将YuDai启动之后，启动飞布控制台，在身份验证新建身份验证器。

参考官方文档：https://ansons-organization.gitbook.io/product-manual/kai-fa-wen-dang/yan-zheng-he-shou-quan/shen-fen-yan-zheng

根据 YuDai 提供的路由信息在飞布控制台添加并配置身份验证。

飞布根据 Issuer 生成服务发现地址，其中 [http://127.0.0.1:9825](http://127.0.0.1:10021/.well-known/openid-configuration)为 YuDai 服务地址，[http://127.0.0.1:9825/.well-known/openid-configuration](http://127.0.0.1:10021/.well-known/openid-configuration) 为服务发现地址，返回数据中定义了包括但不限于：

- 用户端点url：userinfo_endpoint
- jwks访问url：jwks_uri

飞布根据这些信息解析出`jwks_uri`后，就可以对请求进行`token`认证

![img](https://cos.ap-nanjing.myqcloud.com/test-1314985928/admin/image.png)

## 4.2 补充：

> JWKS (JSON Web Key Set) 代表了用于身份验证和授权的一组JSON格式的Web密钥。它是OAuth 2.0和OpenID Connect等身份验证和授权协议中的一部分。JWKS通常用于在安全的方式下公开API的公钥和其他安全参数。
>
> JWKS包含一个包含一组密钥的JSON对象，每个密钥都有一个唯一的标识符（kid），并提供其算法（alg）、密钥类型（kty）和公钥（用于验证签名）或私钥（用于生成签名）等信息。通常，JWKS中的密钥是使用非对称加密算法生成的，如RSA或ECDSA。
>
> 在OAuth 2.0和OpenID Connect中，JWKS用于验证由身份提供者（如认证服务器）签名的访问令牌或ID令牌。客户端可以通过从指定的JWKS端点获取JWKS并验证签名，来确保令牌的身份验证和完整性。
>
> 总结来说，JWKS是一个存储公钥和其他安全参数的JSON对象，用于身份验证和授权协议中的安全验证和签名操作。

## 4.3 前端如何使用

前端请求登录接口，验证成功后会返回token、refreshtoken和对应的过期时间。前端将token值与token过期时间存储在cookie里，refreshToken与refreshToken的过期时间存储到session里，当请求非白名单的接口时，从cookie中取出token值与其过期时间，首先检验是否过期，若过期，则读取refreshToken与其对应的过期时间，若refreshToken没有过期，则携带refreshToken去请求响应的接口，获取新的Token和refreshToken，并携带新的Token去请求接口，若refreshToken也过期，则跳转到默认页面。

# 5. todo list
## 5.1 cookie模式支持
## 5.2 pgsql支持