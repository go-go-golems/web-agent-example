#!/usr/bin/env bash
set -euo pipefail

if rg -n '<style>\{' src/components src/routes src -g '*.tsx' -g '!**/*.stories.tsx' > /tmp/inline-style-matches.txt; then
  echo "Runtime inline <style>{...} blocks are forbidden. Matches:"
  cat /tmp/inline-style-matches.txt
  exit 1
fi

echo "OK: no runtime inline <style>{...} blocks found."
