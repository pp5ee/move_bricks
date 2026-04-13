# Requirement

将 move_bricks 项目 codes/ 目录下所有 Python 文件完全重写为 Go 语言版本，放入 go-impl/ 目录。重写完成后，删除原始 codes/ 目录中的 Python 文件，使目标分支只包含 Go 实现代码。需实现：go.mod（模块名 move_bricks）、main.go、以及 huobi.go（火币API客户端）、coroutine.go（协程/goroutine调度）、listen.go（价格监听）、price_notice.go（价格通知）、email_sender.go（邮件发送）等模块。

---

## Standard Deliverables (mandatory for every project)

- **README.md** — must be included at the project root with: project title & description, prerequisites, installation steps, usage examples with code snippets, configuration options, and project structure overview.
- **Git commits** — use conventional commit prefix `feat:` for all commits.
