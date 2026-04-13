#!/bin/sh
cd '/app/workspaces/31ff3b85-c5b1-4b64-b1ff-104978a78a57' || exit 1
exec 'claude' '--dangerously-skip-permissions' '--print' '/humanize:gen-plan --input requirements/draft.md --output docs/plan.md --direct'
