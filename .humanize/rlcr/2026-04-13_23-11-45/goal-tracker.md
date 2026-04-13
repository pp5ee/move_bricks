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
Rewrite all Python files in move_bricks/codes/ to equivalent Go packages in go-impl/ directory.
Target branch: go-rewrite-v2.

## Acceptance Criteria

### Acceptance Criteria
<!-- Each criterion must be independently verifiable -->
<!-- Claude must extract or define these in Round 0 -->


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
| go-impl/huobi.go - Huobi API client with HMAC-SHA256 signing | AC-2 | completed | coding | claude | REST API client |
| go-impl/listen.go - WebSocket listener with reconnection | AC-3 | completed | coding | claude | Uses gorilla/websocket |
| go-impl/coroutine.go - Async runtime with goroutine+channel | AC-4 | completed | coding | claude | Context-based cancellation |
| go-impl/email_sender.go - Email alerts with SMTP | AC-5 | completed | coding | claude | Uses gomail |
| go-impl/price_notice.go - Price alert functionality | AC-1 | completed | coding | claude | Additional module |
| go-impl/main.go - Main entry point | AC-7 | completed | coding | claude | Compiles, reads config |
| go-impl/go.mod - Module with dependencies | AC-6 | completed | coding | claude | All deps declared |

### Completed and Verified
<!-- Only move tasks here after Codex verification -->
| AC | Task | Completed Round | Verified Round | Evidence |
|----|------|-----------------|----------------|----------|

### Explicitly Deferred
<!-- Items here require strong justification -->
| Task | Original AC | Deferred Since | Justification | When to Reconsider |
|------|-------------|----------------|---------------|-------------------|

### Open Issues
<!-- Issues discovered during implementation -->
| Issue | Discovered Round | Blocking AC | Resolution Path |
|-------|-----------------|-------------|-----------------|
