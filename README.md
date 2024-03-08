# YuDai

取自古代身份验证令牌“鱼袋”

![鱼袋](https://static.aihanfu.net/uploadfile/2014/1128/20141128110743761.jpg)

## 功能介绍

YuDai 是一款轻量级 `OIDC` 身份认证服务，实现了[基于 TOKEN (隐式模式)](https://docs.fireboom.io/ji-chu-ke-shi-hua-kai-fa/shen-fen-yan-zheng/yin-shi-mo-shi) 的登录逻辑，支持账户密码、手机号验证码、微信生态等相关登录方式，可作为 `Casdoor` 的轻量级替代。

> OIDC 是一种用于身份验证和授权的开放式协议。它建立在OAuth 2.0基础上，并为第三方应用程序提供了一种方便的方法来验证用户身份并获取用户信息，例如名称、邮件地址等。OIDC还支持单点登录（SSO），以便用户只需在一个地方登录，就可以访问多个应用程序。

- 账户密码登录
- 手机验证码登录
- 微信小程序登录
- 微信公众号登录
- 微信扫码登录
- ...

## 为什么开发 YuDai

开发 YuDai 主要目的是：为 [Fireboom](https://www.fireboom.cloud) 提供一款轻量级 OIDC 服务，满足 基于 Fireboom 开发 web/app/小程序时的简单登录需求，避免依赖 身份验证供应商，如 auth0 等，以及功能特别庞大的 OIDC 服务，如 `Casdoor` 等，降低系统的复杂度。


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

目前`YuDai`支持 `mysql` 和 `postgres` 数据库。

该项目依赖如下表结构：

db schema：[./docs/oidc.sql](./docs/oidc.sql)

or 

prisma schema：[./docs/oidc.prisma](./docs/oidc.prisma)

有2种方式初始化数据库：自动迁移和脚本导入。

- 自动迁移：`YuDai` 服务启动后会自动同步表信息到数据库，详情见 [代码](https://github.com/fireboomio/yudai-oidc/blob/43982032d7f34493c9b96eaf9a32be191134b3ec/object/adapter.go#L48)
- 脚本导入：使用  `db schema` 或 `prisma schema` （在 fireboom 控制台数据建模操作）文件新建表结构


### 环境变量配置


`YuDai` 目前仅支持 环境变量 配置，后续将支持 `.env` 文件配置。

支持如下环境变量：

- `YUDAI_PORT` ：服务监听端口
- `YUDAI_DB_URL` ：数据库连接字符串
    - `postgres`：`postgres://user:password@host:port/dbname?sslmode=disable`
    - `mysql`：`mysql://user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local`
    >更多参数请分别查看文档[postgres driver](https://pkg.go.dev/github.com/lib/pq)和[mysql driver](https://github.com/go-sql-driver/mysql)
- `YUDAI_DB_PREFIX` ：数据库表前缀，可以用来区分多个服务或者和业务表区分开
- `YUDAI_PC_APP_ID` ：微信pc扫码登录app_id
- `YUDAI_PC_APP_SECRET` ：微信pc扫码登录app_secret
- `YUDAI_H5_APP_ID` ：微信公众号h5登录app_id
- `YUDAI_H5_APP_SECRET` ：微信公众号h5登录app_secret
- `YUDAI_APP_APP_ID` ：微信app登录app_id
- `YUDAI_APP_APP_SECRET` ：微信app登录app_secret
- `YUDAI_MINI_APP_ID` ：微信小程序登录app_id
- `YUDAI_MINI_APP_SECRET` ：微信小程序登录app_secret

### 短信配置

在 `provider` 表配置短信供应商：

> 当前仅支持阿里云短信服务。

```sql
INSERT INTO "provider" ("owner", "name", "created_at", "type", "client_id", "client_secret", "sign_name", "template_code") VALUES ('fireboom', 'provider_sms', '2023-01-17 01:22:33', 'Aliyun SMS', 'xxxxxxxxxxxxxxxx', 'xxxxxxxxxxxxxxxxxxxxxxxxxxx', 'app_name', 'temp_code');
```

- owner：默认用 fireboom
- name：默认用 provider_sms

- client_id：客户端 ID 
- client_secret：客户端密钥
- sign_name：短信签名
- template_code：短信模板编码


详情参考[代码](https://github.com/fireboomio/yudai-oidc/blob/43982032d7f34493c9b96eaf9a32be191134b3ec/api/verification.go#L61)

### jwt 密钥生成

`YuDai` 服务第一次启动时，会在 `cert` 目录下生成文件： `token_jwt_key.key` 和 `token_jwt_key.pem` ，用于 JWT 签名和验证。

手动删除文件后，会重新生成，但会导致之前签发的 token 失效。因此在正式项目中，需要保持这两个文件不发生变化，例如 k8s 部署时，需要用 Volume 挂载这2个文件。

### 启动

```sh
YUDAI_DB_URL="xxxx" ./yudai
```

启动后服务端口为： `9825`

## 前端如何使用

前端请求登录接口：`/api/login`

```sh
curl --location --request POST 'http://localhost:9825/api/login' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--data-raw '{
    "loginType": "password",
    "username": "admin",
    "password": "password",
    "phone": "string",
    "code": "string",
    "platform": "string",
    "exclusive": true
}'
```
验证成功后会返回: 

```json
{
  "msg": "Login Success",
  "code": 200,
  "success": true,
  "data": {
    "accessToken": "xxxx",
    "refreshToken": "xxxx",
    "expireIn": 1709991576
  }
}
```

前端将上述数据存储到浏览器的 `localStorage` 里。

当前端请求接口时，从 `localStorage` 中取出 `accessToken` 和  `expireIn`。

首先检验是否过期，若过期，则读取 `refreshToken`，并携带`refreshToken` 请求响应的接口，获取新的 `accessToken` 和新的`refreshToken`，并携带新的 `accessToken` 去请求接口。

若`refreshToken`也过期，则跳转到登录页面。

前端也可以不存储过期时间，全部由接口是否返回401状态码来判断 `accessToken` 和 `refreshToken` 是否过期。

最佳实践可以参考：https://github.com/fireboomio/amis-admin/blob/dev/src/routes/login/index.tsx

## 如何和 fireboom 结合

`YuDai` 与 fireboom 的最佳实践可参考该项目：https://github.com/fireboomio/amis-admin

### 配置身份认证器

将 `YuDai` 作为 fireboom 的身份认证器， 让 fireboom 可以验证 由  `YuDai` 签发的 token。

在 fireboom 控制台，点击“身份认证”-> “+”，新建身份认证器。

![YuDai身份认证配置](https://cos.ap-nanjing.myqcloud.com/test-1314985928/admin/image.png)

参考上图填写：

- 供应商id : 任意填写
- appid：YuDai 为单应用 oidc服务，因此可随意填写，其他供应商则要按情况填写
- issuer：YuDai 的 api 端点，默认为：http://localhost:9825，支持环境变量
- 基于cookie： YuDai 不支持，请 关闭
- 基于token：开启
- jwks：选择 URL 模式

配置后，保存到该文件：`https://github.com/fireboomio/amis-admin/blob/dev/backend/store/authentication/casdoor.json` 

```json
{
  "name": "casdoor",
  "enabled": true,
  "createTime": "2023-09-14T17:03:03.516456935+08:00",
  "updateTime": "2023-09-14T17:04:52.827399338+08:00",
  "deleteTime": "",
  "issuer": {
    "kind": 1,
    // OIDC_API_URL 默认值为：http://localhost:9825 或 http://oidc:9825
    "environmentVariableName": "OIDC_API_URL" 
  },
  "oidcConfigEnabled": false,
  "oidcConfig": {
    "issuer": null,
    "clientId": {
      "kind": 0
    },
    "clientSecret": {
      "kind": 0
    },
    "queryParameters": null
  },
  "jwksProviderEnabled": true,
  "jwksProvider": {
    "jwksJson": null,
    "userInfoCacheTtlSeconds": 0
  }
}
```

fireboom 根据 Issuer 生成服务发现地址，例如：`http://127.0.0.1:9825/.well-known/openid-configuration`

该地址返回很多数据，包括但不限于：

    - 用户端点url：`userinfo_endpoint`
    - jwks访问url：`jwks_uri`

### 配置 REST 数据源

> 该步骤可以省略，省略后前端直接访问 9825 端口，详情见 [docs/swagger.json](docs/swagger.json)

将 `YuDai` 作为 fireboom 的 rest api 数据源，以便 fireboom 代理 `YuDai` 的 rest api。

代理后的接口见目录：`https://github.com/fireboomio/amis-admin/blob/dev/backend/store/operation/user/casdoor`


例如“登录接口”的访问地址为：

```sh
curl 'http://localhost:9991/operations/user/casdoor/login' \
  -X POST  \
  -H 'Content-Type: application/json' \
  --data-raw '{"code":null,"password":null,"phone":null,"platform":null,"username":null,"loginType":null}' \
  --compressed
```

在 fireboom 控制台，点击“数据源”-> “+” -> “REST API”，新建 rest 数据源。

- 名称：数据源名称，用 casdoor
- 指定 oas：YuDai 的 swagger 文档，见 `docs/oidc.json`
- rest 端点：YuDai 的 api 端点，默认为：`http://localhost:9825`，支持环境变量

配置后见如下文件：`https://github.com/fireboomio/amis-admin/blob/dev/backend/store/datasource/casdoor.json`

```json
{
  "name": "casdoor",
  "enabled": true,
  "createTime": "2023-09-18T10:45:49.51069+08:00",
  "updateTime": "2023-09-18T10:45:52.006846+08:00",
  "deleteTime": "",
  "cacheEnabled": false,
  "kind": 1,
  "customRest": {
    "oasFilepath": "casdoor.json",
    "baseUrl": {
      "kind": 1,
    // OIDC_API_URL 默认值为：http://localhost:9825 或 http://oidc:9825
      "environmentVariableName": "OIDC_API_URL"
    },
    "headers": {
    },
    "responseExtractor": {
      "statusCodeJsonpath": "code",
      "errorMessageJsonpath": "msg"
    }
  },
  "customGraphql": null,
  "customDatabase": null
}
```

详情查看[前往阅读](https://docs.fireboom.io/v/v1.0/ji-chu-ke-shi-hua-kai-fa/shen-fen-yan-zheng/yin-shi-mo-shi#oidc-pei-zhi)

## 未来计划

[ ] cookie模式支持
[x] pgsql支持
