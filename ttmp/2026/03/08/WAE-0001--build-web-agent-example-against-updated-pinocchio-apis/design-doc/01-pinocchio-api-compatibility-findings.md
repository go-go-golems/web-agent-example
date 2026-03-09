---
Title: Pinocchio API compatibility findings
Ticket: WAE-0001
Status: active
Topics:
    - pinocchio
    - webchat
    - go
    - backend
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: pinocchio/cmd/web-chat/main.go
      Note: Current pinocchio web-chat entrypoint shows the new Glazed/Geppetto/webchat integration model to migrate toward
    - Path: web-agent-example/cmd/web-agent-example/engine_from_req.go
      Note: Builds runtime state from incoming chat requests and pinocchio webchat types
    - Path: web-agent-example/cmd/web-agent-example/main.go
      Note: CLI entrypoint wires the webchat server and chat handler against pinocchio APIs
    - Path: web-agent-example/cmd/web-agent-example/sink_wrapper.go
      Note: Adapts sink events into the pinocchio webchat streaming surface
    - Path: web-agent-example/go.mod
      Note: Module manifest currently has no require directives
    - Path: web-agent-example/pkg/discodialogue/timeline.go
      Note: Uses pinocchio timeline/webchat APIs that may have moved or changed
    - Path: web-agent-example/pkg/thinkingmode/timeline.go
      Note: Uses pinocchio timeline/webchat APIs that may have moved or changed
ExternalSources: []
Summary: Track build and test results for web-agent-example against the current workspace pinocchio checkout.
LastUpdated: 2026-03-08T16:11:40.958117618-04:00
WhatFor: ""
WhenToUse: ""
---



# Pinocchio API compatibility findings

## Executive Summary

`go test ./...` does not build `web-agent-example` against the current workspace checkout. The first blocker is mechanical: [web-agent-example/go.mod](/home/manuel/workspaces/2026-03-02/deliver-mento-1/web-agent-example/go.mod#L1) only declares the module and Go version, so none of the imported workspace packages are required. After checking the current workspace package layout, the app is also clearly behind several Geppetto, Glazed, and Pinocchio API moves.

The project is not completely obsolete. Core `pinocchio/pkg/webchat` router/server pieces still exist, but the app needs a migration pass for dependency wiring, request resolution, profile construction, SEM payload imports, and timeline projection types before it will compile again.

## Problem Statement

We updated Pinocchio APIs and need to know whether `web-agent-example` still builds and what parts are now stale. The risk is that multiple libraries evolved together, so a simple `go get` is unlikely to be enough.

Observed result:

- `go test ./...` fails immediately because several imports cannot be resolved.
- The missing imports are not just absent from `go.mod`; some of those package paths no longer exist in the workspace.

## Proposed Solution

Migrate `web-agent-example` in this order:

1. Restore a valid module manifest so local workspace imports resolve.
2. Port the CLI wiring from old Geppetto/Glazed packages to the current sections/fields/values API.
3. Replace the old webchat request-builder/profile model with the current resolver/runtime/profile APIs.
4. Port SEM and timeline integrations away from the removed shared middleware proto package and `TimelineEntityV1`.
5. Re-run `go test ./...` and then a smoke run of `go run ./cmd/web-agent-example`.

## Design Decisions

- Treat the current `go test` failure as a useful signal, but not the whole story.
  The missing `require` directives in [web-agent-example/go.mod](/home/manuel/workspaces/2026-03-02/deliver-mento-1/web-agent-example/go.mod#L1) explain the first failure, but workspace package inspection shows genuine API/package moves behind it.
- Use current `pinocchio/cmd/web-chat` as the migration reference.
  [pinocchio/cmd/web-chat/main.go](/home/manuel/workspaces/2026-03-02/deliver-mento-1/pinocchio/cmd/web-chat/main.go#L1) reflects the new intended integration model for Glazed, Geppetto profiles, and webchat HTTP/runtime wiring.
- Separate "still present" APIs from removed ones.
  `webchat.NewRouter`, `BuildHTTPServer`, `NewFromRouter`, and `WithEventSinkWrapper` still exist, so the app can likely be revived without rewriting the whole server.

## Alternatives Considered

- Only adding `require` directives to `go.mod`.
  Rejected because it would not fix imports such as `github.com/go-go-golems/geppetto/pkg/layers` or `github.com/go-go-golems/glazed/pkg/cmds/parameters`, which no longer exist at those paths.
- Treating the failure as a pure Pinocchio regression.
  Rejected because the breakage spans Geppetto and Glazed package layout changes as well, not just Pinocchio.

## Implementation Plan

1. Add the missing module requirements or run a controlled `go mod tidy` after import paths are updated.
2. Replace old imports in [cmd/web-agent-example/main.go](/home/manuel/workspaces/2026-03-02/deliver-mento-1/web-agent-example/cmd/web-agent-example/main.go#L10) with current packages:
   - `geppetto/pkg/layers` -> `geppetto/pkg/sections`
   - `glazed/pkg/cmds/parameters` -> `glazed/pkg/cmds/fields`
   - `glazed/pkg/cmds/layers` -> `glazed/pkg/cmds/values` and `glazed/pkg/cmds/schema`
3. Port the request path in [engine_from_req.go](/home/manuel/workspaces/2026-03-02/deliver-mento-1/web-agent-example/cmd/web-agent-example/engine_from_req.go#L20) from the old `EngineBuildInput` / `ChatRequestBody` / `RequestBuildError` model to the current `pinocchio/pkg/webchat/http` resolver contract.
4. Port profile registration in [main.go](/home/manuel/workspaces/2026-03-02/deliver-mento-1/web-agent-example/cmd/web-agent-example/main.go#L83) from `webchat.Profile` and root-level `MiddlewareUse` to `geppetto/pkg/profiles.Profile` with `RuntimeSpec` and `PolicySpec`.
5. Port the event sink wrapper in [sink_wrapper.go](/home/manuel/workspaces/2026-03-02/deliver-mento-1/web-agent-example/cmd/web-agent-example/sink_wrapper.go#L14) to the current `EventSinkWrapper` signature, which now receives `infruntime.ConversationRuntimeRequest`.
6. Rewrite timeline projection code in [pkg/discodialogue/timeline.go](/home/manuel/workspaces/2026-03-02/deliver-mento-1/web-agent-example/pkg/discodialogue/timeline.go#L67) and [pkg/thinkingmode/timeline.go](/home/manuel/workspaces/2026-03-02/deliver-mento-1/web-agent-example/pkg/thinkingmode/timeline.go#L38) to emit `TimelineEntityV2` instead of removed `TimelineEntityV1` snapshots.
7. Replace imports from `pinocchio/pkg/sem/pb/proto/sem/middleware` with the new package locations or JSON SEM payload strategy used by current `cmd/web-chat/thinkingmode`.

## Open Questions

- Should `web-agent-example` stay on a custom no-cookie request flow, or should it adopt the standard webchat resolver/profile policy wholesale?
- Should the custom SEM payloads continue using generated protobuf types, or should they mirror the newer JSON-serialized SEM translator approach in Pinocchio?
- Is the goal to minimally restore buildability, or to fully align with current profile registry and debug/timeline features?

## References

- `go test ./...` failure from 2026-03-08:
  missing imports for `geppetto/pkg/layers`, `glazed/pkg/cmds/layers`, `glazed/pkg/cmds/parameters`, and `pinocchio/pkg/sem/pb/proto/sem/middleware`
- Current Geppetto sections API: [geppetto/pkg/sections/sections.go](/home/manuel/workspaces/2026-03-02/deliver-mento-1/geppetto/pkg/sections/sections.go#L1)
- Current Glazed fields API: [glazed/pkg/cmds/fields/definitions.go](/home/manuel/workspaces/2026-03-02/deliver-mento-1/glazed/pkg/cmds/fields/definitions.go#L1)
- Current webchat request/http contract: [pinocchio/pkg/webchat/http/api.go](/home/manuel/workspaces/2026-03-02/deliver-mento-1/pinocchio/pkg/webchat/http/api.go#L1)
- Current profile types: [geppetto/pkg/profiles/types.go](/home/manuel/workspaces/2026-03-02/deliver-mento-1/geppetto/pkg/profiles/types.go#L1)
