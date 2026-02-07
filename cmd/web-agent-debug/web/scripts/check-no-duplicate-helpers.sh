#!/usr/bin/env bash
set -euo pipefail

PATTERN='function\s+(getEventTypeInfo|getKindInfo|getKindIcon|truncateText|formatPhase|formatPhaseName)\b|const\s+(getEventTypeInfo|getKindInfo|getKindIcon|truncateText|formatPhase|formatPhaseName)\s*='

if rg -n "$PATTERN" src/components src/routes -g '*.{ts,tsx}' > /tmp/duplicate-helper-matches.txt; then
  echo "Duplicate helper regression detected in runtime components/routes:"
  cat /tmp/duplicate-helper-matches.txt
  exit 1
fi

echo "OK: no duplicate helper signature regressions found in src/components and src/routes."
