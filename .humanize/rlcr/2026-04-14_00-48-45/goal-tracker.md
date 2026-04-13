# Goal Tracker

<!--
This file tracks the ultimate goal, acceptance criteria, and plan evolution.
It prevents goal drift by maintaining a persistent anchor across all rounds.

RULES:
- IMMUTABLE SECTION: Do not modify after initialization
- MUTABLE SECTION: Update each round, but document all changes
- Every task must be in one of: Active, Completed, or Deferred
- Deferred items require explicit justification
-->

## IMMUTABLE SECTION
<!-- Do not modify after initialization -->

### Ultimate Goal
将 move_bricks 项目 codes/ 目录下的 Python 交易机器人代码完全重写为 Go 语言版本，输出到 go-impl/ 目录。重写后删除原始 Python 文件，使目标分支 python-to-go-v3 只包含 Go 实现代码。

## Acceptance Criteria

### Acceptance Criteria
<!-- Each criterion must be independently verifiable -->
<!-- Claude must extract or define these in Round 0 -->


- AC-1: go-impl/ 目录包含所有模块的 Go 实现
  - Positive Tests: go-impl/ 包含 huobi.go, coroutine.go, listen.go, price_notice.go, email_sender.go, main.go
  - Negative Tests: 不含 Python 文件（.py）

- AC-2: go-impl/huobi.go 实现火币 API 客户端
  - Positive Tests: 包含 HMAC-SHA256 签名、订单操作、余额查询方法
  - Negative Tests: 不含硬编码 API key

- AC-3: go-impl/listen.go 实现 WebSocket 价格监听
  - Positive Tests: 使用 gorilla/websocket，支持自动重连
  - Negative Tests: 不使用 Python 协程模式

- AC-4: go-impl/coroutine.go 实现异步任务调度
  - Positive Tests: 使用 goroutine 和 channel 替代 Python asyncio
  - Negative Tests: 不依赖 Python 运行时

- AC-5: go-impl/email_sender.go 实现邮件发送
  - Positive Tests: 使用 Go SMTP 库发送告警邮件
  - Negative Tests: 不依赖 Python smtplib

- AC-6: go-impl/go.mod 存在且模块名为 move_bricks
  - Positive Tests: go.mod 包含 module move_bricks 和必要依赖
  - Negative Tests: 无

---

## MUTABLE SECTION
<!-- Update each round with justification for changes -->

### Plan Version: 1 (Updated: Round 0)

#### Plan Evolution Log
<!-- Document any changes to the plan with justification -->
| Round | Change | Reason | Impact on AC |
|-------|--------|--------|--------------|
| 0 | Initial plan | - | - |

#### Active Tasks
<!-- Map each task to its target Acceptance Criterion and routing tag -->
| Task | Target AC | Status | Tag | Owner | Notes |
|------|-----------|--------|-----|-------|-------|
| 分析 Python 代码结构 | AC-1 | completed | analyze | codex | 分析 codes/ 目录结构 |
| 创建 go.mod 和依赖 | AC-6 | completed | coding | claude | 初始化 Go 模块 |
| 实现 huobi.go | AC-2 | completed | coding | claude | HMAC-SHA256 签名，火币 API |
| 实现 coroutine.go | AC-4 | completed | coding | claude | goroutine 和 channel |
| 实现 listen.go | AC-3 | completed | coding | claude | WebSocket 价格监听 |
| 实现 price_notice.go | AC-1 | completed | coding | claude | 价格告警 |
| 实现 email_sender.go | AC-5 | completed | coding | claude | SMTP 邮件发送 |
| 实现 main.go | AC-1 | completed | coding | claude | 主入口 |
| 删除 Python 文件 | AC-7 | completed | coding | claude | 删除 codes/ 目录文件 |

### Completed and Verified
<!-- Only move tasks here after Codex verification -->
| AC | Task | Completed Round | Verified Round | Evidence |
|----|------|-----------------|----------------|----------|
| AC-1 | go-impl/ 目录包含所有模块 | 0 | - | 7 个 Go 文件已创建 |
| AC-2 | huobi.go 实现 | 0 | - | HMAC-SHA256 签名已实现 |
| AC-3 | listen.go 实现 | 0 | - | gorilla/websocket 已使用 |
| AC-4 | coroutine.go 实现 | 0 | - | goroutine/channel 已使用 |
| AC-5 | email_sender.go 实现 | 0 | - | Go SMTP 库已使用 |
| AC-6 | go.mod 存在 | 0 | - | module move_bricks 已定义 |
| AC-7 | Python 文件已删除 | 0 | - | 已提交删除 11 个文件 |

### Explicitly Deferred
<!-- Items here require strong justification -->
| Task | Original AC | Deferred Since | Justification | When to Reconsider |
|------|-------------|----------------|---------------|-------------------|

### Open Issues
<!-- Issues discovered during implementation -->
| Issue | Discovered Round | Blocking AC | Resolution Path |
|-------|-----------------|-------------|-----------------|
