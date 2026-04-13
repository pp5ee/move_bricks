# Requirement

将 move_bricks/codes/ 目录下所有 Python 文件转写为 Go 语言。目标：在 go-impl/ 目录下创建等价的 Go 包，保留原有业务逻辑（交易所 API 调用、协程模型、WebSocket 行情订阅、邮件告警）。使用 Go 的 goroutine + channel 替代 Python asyncio，用 github.com/gorilla/websocket 做 WebSocket，gomail 做邮件，resty 做 HTTP 请求。API key 改为环境变量。输出到 go-impl/ 目录，新建分支 go-rewrite-v2。

---

## Standard Deliverables (mandatory for every project)

- **README.md** — must be included at the project root with: project title & description, prerequisites, installation steps, usage examples with code snippets, configuration options, and project structure overview.
- **Git commits** — use conventional commit prefix `feat:` for all commits.
