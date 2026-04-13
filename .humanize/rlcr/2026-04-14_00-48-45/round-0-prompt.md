Read and execute below with ultrathink

## Goal Tracker Setup (REQUIRED FIRST STEP)

Before starting implementation, you MUST initialize the Goal Tracker:

1. Read @/app/workspaces/31ff3b85-c5b1-4b64-b1ff-104978a78a57/.humanize/rlcr/2026-04-14_00-48-45/goal-tracker.md
2. If the "Ultimate Goal" section says "[To be extracted...]", extract a clear goal statement from the plan
3. If the "Acceptance Criteria" section says "[To be defined...]", define 3-7 specific, testable criteria
4. Populate the "Active Tasks" table with tasks from the plan, mapping each to an AC and filling Tag/Owner
5. Write the updated goal-tracker.md

**IMPORTANT**: The IMMUTABLE SECTION can only be modified in Round 0. After this round, it becomes read-only.

---

## Implementation Plan

For all tasks that need to be completed, please use the Task system (TaskCreate, TaskUpdate, TaskList) to track each item in order of importance.
You are strictly prohibited from only addressing the most important issues - you MUST create Tasks for ALL discovered issues and attempt to resolve each one.

## Task Tag Routing (MUST FOLLOW)

Each task must have one routing tag from the plan: `coding` or `analyze`.

- Tag `coding`: Claude executes the task directly.
- Tag `analyze`: Claude must execute via `/humanize:ask-codex`, then integrate Codex output.
- Keep Goal Tracker "Active Tasks" columns **Tag** and **Owner** aligned with execution (`coding -> claude`, `analyze -> codex`).
- If a task has no explicit tag, default to `coding` (Claude executes directly).

# Plan: Python to Go Rewrite (move_bricks)

## Goal Description
将 move_bricks 项目 codes/ 目录下的 Python 交易机器人代码完全重写为 Go 语言版本，输出到 go-impl/ 目录。重写后删除原始 Python 文件，使目标分支 python-to-go-v3 只包含 Go 实现代码。

## Acceptance Criteria

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

- AC-7: codes/ 目录中的 Python 文件已被删除
  - Positive Tests: codes/ 目录不存在或为空
  - Negative Tests: 不含 .py 文件

## Path Boundaries

### Upper Bound (Maximum Acceptable Scope)
- 可以添加 README 更新
- 可以增加基础单元测试框架

### Lower Bound (Minimum Acceptable Scope)
- 最少需要 7 个 Go 文件和 go.mod

### Allowed Choices
- Can use: gorilla/websocket, gomail, fasthttp, standard library
- Cannot use: Python 依赖或 CGo

## Feasibility Hints and Suggestions

### Conceptual Approach
1. 分析 codes/ 中的 Python 文件（huobi.py, coroutine.py, listen.py, price_notice.py, email_sender.py）
2. 将每个 Python 模块转换为对应的 Go 文件
3. 使用 goroutine/channel 替代 Python asyncio
4. 删除原始 Python 文件
5. 提交所有变更

### Relevant References
- 火币 API: https://huobiapi.github.io/docs/spot/v1/en/
- gorilla/websocket: https://github.com/gorilla/websocket
- gomail: https://github.com/go-mail/gomail

## Dependencies and Sequence
1. 分析 Python 代码结构
2. 创建 go.mod 和依赖
3. 实现各模块
4. 删除 Python 文件
5. 提交

---

## BitLesson Selection (REQUIRED FOR EACH TASK)

Before executing each task or sub-task, you MUST:

1. Read @/app/workspaces/31ff3b85-c5b1-4b64-b1ff-104978a78a57/.humanize/bitlesson.md
2. Run `bitlesson-selector` for each task/sub-task to select relevant lesson IDs
3. Follow the selected lesson IDs (or `NONE`) during implementation

Include a `## BitLesson Delta` section in your summary with:
- Action: none|add|update
- Lesson ID(s): NONE or comma-separated IDs
- Notes: what changed and why (required if action is add or update)

Reference: @/app/workspaces/31ff3b85-c5b1-4b64-b1ff-104978a78a57/.humanize/bitlesson.md

---

## Goal Tracker Rules

Throughout your work, you MUST maintain the Goal Tracker:

1. **Before starting a task**: Mark it as "in_progress" in Active Tasks
   - Confirm Tag/Owner routing is correct before execution
2. **After completing a task**: Move it to "Completed and Verified" with evidence (but mark as "pending verification")
3. **If you discover the plan has errors**:
   - Do NOT silently change direction
   - Add entry to "Plan Evolution Log" with justification
   - Explain how the change still serves the Ultimate Goal
4. **If you need to defer a task**:
   - Move it to "Explicitly Deferred" section
   - Provide strong justification
   - Explain impact on Acceptance Criteria
5. **If you discover new issues**: Add to "Open Issues" table

---

Note: You MUST NOT try to exit `start-rlcr-loop` loop by lying or edit loop state file or try to execute `cancel-rlcr-loop`

After completing the work, please:
0. If you have access to the `code-simplifier` agent, use it to review and optimize the code you just wrote
1. Finalize @/app/workspaces/31ff3b85-c5b1-4b64-b1ff-104978a78a57/.humanize/rlcr/2026-04-14_00-48-45/goal-tracker.md (this is Round 0, so you are initializing it - see "Goal Tracker Setup" above)
2. Commit your changes with a descriptive commit message
3. Write your work summary into @/app/workspaces/31ff3b85-c5b1-4b64-b1ff-104978a78a57/.humanize/rlcr/2026-04-14_00-48-45/round-0-summary.md

Note: Since `--push-every-round` is enabled, you must push your commits to remote after each round.
