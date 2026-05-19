GO-WATCH 链上众筹监控服务

配套 FundingMonitor 合约的后端服务，用于监听链上众筹事件、存储数据、提供 API 接口，实现实时资金监控与排行榜功能。⚠️当前仓库为Go配套仓库，适配的Solidity代码仓库在https://github.com/wyxct/Solidity_project


✨ 项目亮点

事件实时监听：通过轮询 RPC 节点，捕获 Funded/ReFunded/Distribute 事件并写入数据库
数据持久化：GORM + PostgreSQL 存储每笔交易记录，同时维护用户捐款总额统计表
排行榜接口：自动聚合用户捐款数据，提供按金额排序的实时排行榜
标准分层架构：controller → service → dao → models，结构清晰易维护
轻量配置：YAML 配置文件管理 RPC、数据库等敏感信息，不硬编码


🛠️ 技术栈
表格模块技术 / 工具语言Go 1.21+链上交互go-ethereum/ethclientORMGORM v2数据库PostgreSQL配置Viper/YAMLAPI标准 net/http/ Gin（如使用）

📂 项目结构plaintext.
├── abi/                  # 合约 ABI 定义
│   ├── abi.go            # ABI 解析封装
│   └── Funding.json      # 合约 ABI 文件
├── config/               # 配置文件与加载
│   ├── config.go         # 配置结构体定义
│   └── config.yaml       # 项目配置（RPC、DB 等）
├── controller/           # API 控制器层
│   └── funding_controller.go
├── dao/                  # 数据访问层（所有 DB 操作）
│   ├── fund_record_dao.go          # 捐款记录 DAO
│   ├── refund_record_dao.go        # 退款记录 DAO
│   ├── distribute_record_dao.go   # 分账记录 DAO
│   └── fund_total_record_dao.go   # 用户总额统计 DAO
├── models/               # 数据库模型定义
│   ├── FundRecord.go              # 捐款记录模型
│   ├── ReFundRecord.go            # 退款记录模型
│   ├── DistributeRecord.go       # 分账记录模型
│   ├── FundingTotalRecord.go     # 用户总额统计模型
│   └── requestModel.go           # API 请求结构体
├── routes/               # 路由定义
│   └── routes.go
├── service/              # 业务逻辑层
│   ├── funding_service.go        # 众筹相关业务逻辑
│   └── watcher.go                # 链上事件监听服务
├── go.mod
├── go.sum
└── main.go               # 程序入口


🚀 核心功能说明
1. 链上事件监听（watcher.go）

轮询 RPC 节点，从上次处理的块高开始，捕获新的 Funded/ReFunded/Distribute 事件
解析事件数据，自动更新：
fund_records：每笔捐款明细
refund_records：退款明细
distribute_records：分账明细
fund_total_records：用户捐款总额（自动累加 / 扣减）



2. 数据存储与统计

使用 GORM 操作 PostgreSQL，确保数据安全写入
自动维护用户捐款总额表，支持实时查询与排序
所有交易记录永久保存，支持对账与历史回溯

3. API 接口
提供以下核心接口（示例）：

GET /api/funding/list：分页查询捐款记录
GET /api/funding/rank：按捐款总额排序的排行榜
GET /api/refund/list：查询退款记录
GET /api/distribute/list：查询分账记录


📦 运行与部署
1. 配置文件 config.yamlyamlchain:
  rpc_url: "https://rpc-url"
  contract: "0x..."

db:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "yourpass"
  name: "funding_watch"

2. 安装依赖bash运行go mod tidy

3. 运行服务bash运行go run main.go


⚠️ 安全说明

敏感配置（RPC、数据库密码）通过 config.yaml 管理，需加入 .gitignore 避免泄露
所有数据库操作通过 DAO 层封装，避免 SQL 注入风险
链上事件监听逻辑已处理网络波动、RPC 超时等异常情况
