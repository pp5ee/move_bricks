# Ask Codex Input

## Question

You are a code review expert. Analyze this draft for a Python-to-Go conversion project and provide structured feedback.

Repository context: This is the move_bricks project - a trading bot for huobi exchange. The goal is to convert Python files in codes/ directory to Go and place them in go-impl/ directory. After conversion, delete the original Python files.

Draft content:
# Requirement
将 move_bricks 项目 codes/ 目录下所有 Python 文件完全重写为 Go 语言版本，放入 go-impl/ 目录。重写完成后，删除原始 codes/ 目录中的 Python 文件，使目标分支只包含 Go 实现代码。需实现：go.mod（模块名 move_bricks）、main.go、以及 huobi.go（火币API客户端）、coroutine.go（协程/goroutine调度）、listen.go（价格监听）、price_notice.go（价格通知）、email_sender.go（邮件发送）等模块。

---
## Standard Deliverables (mandatory for every project)
- **README.md** — must be included at the project root with: project title & description, prerequisites, installation steps, usage examples with code snippets, configuration options, and project structure overview.
- **Git commits** — use conventional commit prefix  for all commits.

Also note: The go-impl/ directory already contains some Go files: main.go, huobi.go, coroutine.go, listen.go, price_notice.go, email_sender.go, go.mod

Critique the draft by providing:
1. CORE_RISKS: highest-risk assumptions and potential failure modes
2. MISSING_REQUIREMENTS: likely omitted requirements or edge cases
3. TECHNICAL_GAPS: feasibility or architecture gaps
4. ALTERNATIVE_DIRECTIONS: viable alternatives with tradeoffs
5. QUESTIONS_FOR_USER: questions that need explicit human decisions
6. CANDIDATE_CRITERIA: candidate acceptance criteria suggestions

## Configuration

- Model: claude-sonnet-4.6
- Effort: high
- Timeout: 3600s
- Timestamp: 2026-04-14_00-43-27
