# AI 辅助研发成效看板

一个用于展示 AI 辅助研发成效的单页看板项目。前端是单个 `index.html`，后端使用 Go 标准库托管页面，并使用 SQLite 存储看板数据。

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

- Go 1.22+
- SQLite
- `modernc.org/sqlite` 纯 Go SQLite 驱动
- 原生 HTML / CSS / JavaScript

## 本地运行

```bash
go run .
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

## 配置

环境变量：

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `PORT` | `8080` | HTTP 服务端口 |
| `SQLITE_PATH` | `dashboard.db` | SQLite 数据库文件路径 |

示例：

```bash
PORT=9000 SQLITE_PATH=./data/dashboard.db go run .
```

Windows PowerShell：

```powershell
$env:PORT="9000"
$env:SQLITE_PATH="./data/dashboard.db"
go run .
```

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

## 构建

```bash
go build -o ai-rd-dashboard .
```

运行构建产物：

```bash
./ai-rd-dashboard
```
