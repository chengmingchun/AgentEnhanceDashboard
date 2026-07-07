# AI 辅助研发成效看板

用于展示 AI 辅助研发成效的单页看板项目。前端是单个 `index.html`，后端提供 Go 与 Python 两种实现，数据统一存储在 SQLite。

## 功能

- 总览 AI 生成代码行、效率提升、平均评分、已合并 MR、高贡献小组等指标
- 明细表展示需求名称、MR 名称、负责人、小组、AI 生成代码数、效率提升、生成效果评分和问题反馈
- MR 名称支持点击跳转
- 支持按时间、人员、小组和关键词过滤
- 支持按人、按小组统计 AI 使用情况
- 支持浅色 / 深色主题切换
- 支持导出当前筛选视图为 CSV
- SQLite 启动自动建表，空库时自动写入样例数据

## 技术栈

- 前端：原生 HTML / CSS / JavaScript，单文件交付
- 数据库：SQLite
- Go 后端：Go 1.22+，`modernc.org/sqlite`
- Python 后端：Python 3.10+，仅使用标准库

## Python 后端设计

Python 后端在 `python_backend/` 目录下，按职责分层：

| 文件 | 职责 |
| --- | --- |
| `config.py` | 环境变量与运行配置 |
| `models.py` | 领域模型与 API 输出转换 |
| `database.py` | SQLite 连接工厂、迁移、种子数据初始化 |
| `repository.py` | Repository 抽象与 SQLite 实现 |
| `service.py` | 应用服务层，承载用例逻辑 |
| `http_app.py` | HTTP Controller、路由与应用工厂 |
| `seed.py` | 初始化样例数据 |

设计上使用了 Repository、Service Layer、Application Factory、Connection Factory 等模式。后续如果要接入真实研发流水线、MR 平台或用户系统，可以优先扩展 Repository 或 Service，不需要改动页面与 HTTP 协议层。

## 使用 Python 后端运行

```bash
python run_python.py
```

默认访问：

```text
http://localhost:8080
```

健康检查：

```text
http://localhost:8080/healthz
```

数据接口：

```text
http://localhost:8080/api/records
```

Windows PowerShell 自定义端口与数据库路径：

```powershell
$env:PORT="9000"
$env:SQLITE_PATH="./data/dashboard.db"
python run_python.py
```

## 使用 Go 后端运行

```bash
go run .
```

构建：

```bash
go build -o ai-rd-dashboard .
```

## 配置

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `PORT` | `8080` | HTTP 服务端口 |
| `SQLITE_PATH` | `dashboard.db` | SQLite 数据库文件路径 |

## 数据表

启动时会自动创建 `records` 表：

| 字段 | 说明 |
| --- | --- |
| `requirement` | 需求名称 |
| `mr` | MR 名称 |
| `mr_url` | MR 跳转地址 |
| `owner` | 负责人 |
| `group_name` | 所属小组，如 `cos`、`d5` |
| `lines` | AI 生成代码行数 |
| `efficiency` | 效率提升百分比 |
| `score` | 生成效果评分 |
| `problem` | 遇到的问题 |
| `record_date` | 记录日期 |
| `status` | MR 状态 |

## 接口

### `GET /api/records`

返回所有看板记录，按日期倒序排列。

### `GET /healthz`

返回服务健康状态：

```json
{"status":"ok"}
```
