# Ask Codex Input

## Question

Analyze the following draft for a Python to Go language rewrite project. Provide feedback on feasibility, risks, and potential issues.

## Repository Context
- Project location: /app/workspaces/f3e61245-def3-4c88-b1fa-bf264386b545
- Current code: Python files in move_bricks/codes/ including:
  - coroutine.py, email_sender.py, httpserver.py, huobi.py, listen.py, move_bricks.py, price_notice.py, test_trade_huobi.py
- Business logic: exchange API calls, coroutine model, WebSocket market data subscription, email alerts
- Target: Go implementation in go-impl/ directory

## Draft Content
# Requirement

将 move_bricks/codes/ 目录下所有 Python 文件转写为 Go 语言。目标：在 go-impl/ 目录下创建等价的 Go 包，保留原有业务逻辑（交易所 API 调用、协程模型、WebSocket 行情订阅、邮件告警）。使用 Go 的 goroutine + channel 替代 Python asyncio，用 github.com/gorilla/websocket 做 WebSocket，gomail 做邮件，resty 做 HTTP 请求。API key 改为环境变量。输出到 go-impl/ 目录，新建分支 go-rewrite-v2。

---

## Standard Deliverables (mandatory for every project)

- **README.md** — must be included at the project root with: project title & description, prerequisites, installation steps, usage examples with code snippets, configuration options, and project structure overview.
- **Git commits** — use conventional commit prefix feat: for all commits.

## Required Output Format
Provide your analysis in this format:
- CORE_RISKS: highest-risk assumptions and potential failure modes
- MISSING_REQUIREMENTS: likely omitted requirements or edge cases
- TECHNICAL_GAPS: feasibility or architecture gaps
- ALTERNATIVE_DIRECTIONS: viable alternatives with tradeoffs
- QUESTIONS_FOR_USER: questions that need explicit human decisions
- CANDIDATE_CRITERIA: candidate acceptance criteria suggestions

## Configuration

- Model: claude-sonnet-4.6
- Effort: high
- Timeout: 3600s
- Timestamp: 2026-04-13_23-00-41
