---
Title: Build web-agent-example against updated pinocchio APIs
Ticket: WAE-0001
Status: active
Topics:
    - pinocchio
    - webchat
    - go
    - backend
DocType: index
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2026-03-08T16:10:46.605959386-04:00
WhatFor: ""
WhenToUse: ""
---

# Build web-agent-example against updated pinocchio APIs

## Overview

This ticket tracks whether `web-agent-example` still builds against the current workspace version of Pinocchio and adjacent libraries. The initial investigation showed that it did not compile because the app still targeted older Geppetto, Glazed, webchat, SEM, and timeline APIs.

The port is now implemented and verified under the workspace `go.work`. The example builds again against the updated Pinocchio APIs, with the remaining open question limited to whether the module also needs a standalone `GOWORK=off` dependency path.

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field
- **Compatibility Findings**: [design-doc/01-pinocchio-api-compatibility-findings.md](./design-doc/01-pinocchio-api-compatibility-findings.md)
- **Implementation Diary**: [reference/01-diary.md](./reference/01-diary.md)

## Status

Current status: **active**. Build and smoke validation are passing under the workspace checkout.

## Topics

- pinocchio
- webchat
- go
- backend

## Tasks

See [tasks.md](./tasks.md) for the current task list.

## Changelog

See [changelog.md](./changelog.md) for recent changes and decisions.

## Structure

- design/ - Architecture and design documents
- reference/ - Prompt packs, API contracts, context summaries
- playbooks/ - Command sequences and test procedures
- scripts/ - Temporary code and tooling
- various/ - Working notes and research
- archive/ - Deprecated or reference-only artifacts
