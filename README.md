# Npool account gateway

[![Test](https://github.com/NpoolPlatform/account-gateway/actions/workflows/main.yml/badge.svg?branch=master)](https://github.com/NpoolPlatform/account-gateway/actions/workflows/main.yml)

[目录](#目录)
- [功能](#功能)
- [命令](#命令)
- [步骤](#步骤)
- [最佳实践](#最佳实践)
- [关于mysql](#关于mysql)

-----------
### 功能
- [x] 账户管理网关

#### 地址类别
- 平台冷钱包
  - 线下创建的离线地址。支付地址和充值地址达到设置的collect amount，则会归集到平台冷钱包
- 用户冷钱包
  - 线下创建的离线地址。用户热钱包超过限额部分会被转入平台设置的用户冷钱包
- 用户热钱包
  - 每天的商品收益，属于用户部分的金额会被转入平台设置的用户热钱包
- 支付地址
  - 订单交易的支付地址，未支付订单的支付地址具有唯一性，用于检测用户是否支付该笔交易
- 充值地址
  - 用户在平台充值余额使用的充值地址，每个用户的充值地址具有唯一性
- 收益地址
  - 每个商品会绑定一个收益地址，收益分账的流程参考staker-manager仓库的README.md
- 提现地址
  - 用户个人地址，用于提取用户在平台拥有的资产

### 命令
* make init ```初始化仓库，创建go.mod```
* make verify ```验证开发环境与构建环境，检查code conduct```
* make verify-build ```编译目标```
* make test ```单元测试```
* make generate-docker-images ```生成docker镜像```
* make account-gateway ```单独编译服务```
* make account-gateway-image ```单独生成服务镜像```
* make deploy-to-k8s-cluster ```部署到k8s集群```

### 最佳实践
* 每个服务只提供单一可执行文件，有利于docker镜像打包与k8s部署管理
* 每个服务提供http调试接口，通过curl获取调试信息
* 集群内服务间direct call调用通过服务发现获取目标地址进行调用
* 集群内服务间event call调用通过rabbitmq解耦

### 关于mysql
* 创建app后，从app.Mysql()获取本地mysql client
* [文档参考](https://entgo.io/docs/sql-integration)

### API规范(假设目标对象为Object)
# 前端接口或App管理员接口
* 创建 ```CreateObject```
* 批量创建 ```CreateObjects```
* 更新 ```UpdateObject```
* 更新指定域 ```UpdateObjectFields```
* 原子加指定域 ```AddObjectFields```
* 用ID查询 ```GetObject```
* 条件查询单条 ```GetObjectOnly```
* 条件查询多条 ```GetObjects```
* 条件计数 ```CountObjects```
* 查询ID是否存在 ```ExistObject```
* 查询条件是否存在 ```ExistObjectConds```
* ID删除 ```DeleteObject```

# 大后台管理员接口
* 为其他App创建 ```CreateAppObject```
* 批量为其他App创建 ```CreateAppObjects```
* 查询其他App多条 ```GetAppObjects```

# 注意事项
* 网关对请求的http header字段做如下转换(request为http请求body)
  * X-App-ID -> request.AppID,
  * X-App-ID -> request.Info.AppID, 如果request携带Info字段
  * X-User-ID -> request.UserID,
  * X-User-ID -> request.Info.UserID, 如果request携带Info字段
  * X-Lang-ID -> request.LangID,
  * X-Lang-ID -> request.Info.LangID, 如果request携带Info字段
  * X-App-User-Token -> request.Token,
  * X-App-User-Token -> request.Info.Token, 如果request携带Info字段

* 对于通过条件查询的前端API或App管理员API，其request中应总是包含AppID字段，实现时应该总是将该AppID字段添加到查询条件
* 对于批量创建的前端或App管理员API，其request中应总是包含AppID字段，实现时应该总是将该AppID字段覆盖request.Info或request.Infos中的AppID
* 更新API的实现不应更新AppID域
* 对于通过条件查询或创建其他App数据的大后台管理员API，其request中应总是包含TargetAppID，实现时应该总是将TargetAppID字段覆盖request.Info或request.Infos中的AppID
