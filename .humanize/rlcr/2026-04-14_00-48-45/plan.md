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
