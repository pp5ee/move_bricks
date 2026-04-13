# move_bricks Go Rewrite Plan (v2)

## Goal Description
Rewrite all Python files in move_bricks/codes/ to equivalent Go packages in go-impl/ directory.
Target branch: go-rewrite-v2.

## Acceptance Criteria

- AC-1: All Python files have Go equivalents in go-impl/
  - Positive Tests:
    - go-impl/ directory contains .go files for each Python module
    - go-impl/go.mod exists with module definition
  - Negative Tests:
    - No Python code remains in go-impl/

- AC-2: Huobi API client (huobi.py → go-impl/huobi.go)
  - Positive Tests:
    - Implements REST API with HMAC-SHA256 signing
    - API keys from environment variables (HUOBI_ACCESS_KEY, HUOBI_SECRET_KEY)
  - Negative Tests:
    - No hardcoded credentials

- AC-3: WebSocket listener (listen.py → go-impl/listen.go)
  - Positive Tests:
    - Uses github.com/gorilla/websocket
    - Handles reconnection
  - Negative Tests:
    - No blocking calls in goroutines

- AC-4: Async runtime (coroutine.py → go-impl/coroutine.go)
  - Positive Tests:
    - Uses goroutine + channel pattern
    - Context-based cancellation
  - Negative Tests:
    - No global state

- AC-5: Email alerts (email_sender.py → go-impl/email_sender.go)
  - Positive Tests:
    - Uses gomail or net/smtp
    - Config from environment variables
  - Negative Tests:
    - No hardcoded SMTP credentials

- AC-6: go.mod with all dependencies
  - Positive Tests:
    - github.com/gorilla/websocket listed
    - gopkg.in/gomail.v2 or net/smtp used
    - go-resty/resty/v2 for HTTP

- AC-7: Main entry point (move_bricks.py → go-impl/main.go or cmd/main.go)
  - Positive Tests:
    - Compiles without errors
    - Reads config from environment

## Path Boundaries

### Upper Bound
- Implement all modules with full functionality matching Python originals
- Include unit tests for each module

### Lower Bound
- At minimum: all .go files created with correct package structure and function signatures
- go.mod with all required dependencies

## Implementation Notes
- Output directory: go-impl/
- Target branch: go-rewrite-v2
- Key packages: github.com/gorilla/websocket, go-resty/resty/v2, gopkg.in/gomail.v2
- Use goroutine + channel instead of asyncio
- Use conventional commit prefix: feat:
