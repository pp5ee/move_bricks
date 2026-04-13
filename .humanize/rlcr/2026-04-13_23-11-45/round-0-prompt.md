Read and execute below with ultrathink

## Goal Tracker Setup (REQUIRED FIRST STEP)

Before starting implementation, you MUST initialize the Goal Tracker:

1. Read @/app/workspaces/f3e61245-def3-4c88-b1fa-bf264386b545/.humanize/rlcr/2026-04-13_23-11-45/goal-tracker.md
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

---

## BitLesson Selection (REQUIRED FOR EACH TASK)

Before executing each task or sub-task, you MUST:

1. Read @/app/workspaces/f3e61245-def3-4c88-b1fa-bf264386b545/.humanize/bitlesson.md
2. Run `bitlesson-selector` for each task/sub-task to select relevant lesson IDs
3. Follow the selected lesson IDs (or `NONE`) during implementation

Include a `## BitLesson Delta` section in your summary with:
- Action: none|add|update
- Lesson ID(s): NONE or comma-separated IDs
- Notes: what changed and why (required if action is add or update)

Reference: @/app/workspaces/f3e61245-def3-4c88-b1fa-bf264386b545/.humanize/bitlesson.md

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
1. Finalize @/app/workspaces/f3e61245-def3-4c88-b1fa-bf264386b545/.humanize/rlcr/2026-04-13_23-11-45/goal-tracker.md (this is Round 0, so you are initializing it - see "Goal Tracker Setup" above)
2. Commit your changes with a descriptive commit message
3. Write your work summary into @/app/workspaces/f3e61245-def3-4c88-b1fa-bf264386b545/.humanize/rlcr/2026-04-13_23-11-45/round-0-summary.md

Note: Since `--push-every-round` is enabled, you must push your commits to remote after each round.
