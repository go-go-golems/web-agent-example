# Changelog

## 2026-03-08

- Initial workspace created


## 2026-03-08

Ran go test ./... against the workspace checkout. Build fails immediately because web-agent-example/go.mod has no require directives, and the code still references older Geppetto/Glazed/webchat APIs that have moved or been redesigned.

### Related Files

- /home/manuel/workspaces/2026-03-02/deliver-mento-1/web-agent-example/cmd/web-agent-example/main.go — Imports and profile/router wiring target older Geppetto/Glazed/webchat APIs
- /home/manuel/workspaces/2026-03-02/deliver-mento-1/web-agent-example/go.mod — Missing module requirements prevent import resolution before deeper compile errors surface
- /home/manuel/workspaces/2026-03-02/deliver-mento-1/web-agent-example/pkg/discodialogue/timeline.go — Timeline projection still targets TimelineEntityV1 and old SEM middleware protos


## 2026-03-09

Ported web-agent-example to the current workspace Geppetto/Glazed/Pinocchio APIs. The app now uses the current sections/fields/values command model, a current webchat request resolver/runtime composer, JSON SEM payloads for custom middleware events, and TimelineEntityV2 projections. Verified with go test ./... and go run ./cmd/web-agent-example serve --help under the workspace go.work.

### Related Files

- /home/manuel/workspaces/2026-03-02/deliver-mento-1/web-agent-example/cmd/web-agent-example/engine_from_req.go — Replaced the old engine builder with the current webchat request resolver contract
- /home/manuel/workspaces/2026-03-02/deliver-mento-1/web-agent-example/cmd/web-agent-example/main.go — Replaced legacy command and router wiring with current webchat server composition
- /home/manuel/workspaces/2026-03-02/deliver-mento-1/web-agent-example/cmd/web-agent-example/runtime_composer.go — Added a current runtime composer that maps resolved runtime middleware into Geppetto middleware
- /home/manuel/workspaces/2026-03-02/deliver-mento-1/web-agent-example/pkg/discodialogue/sem.go — Replaced removed shared SEM proto usage with local JSON payloads
- /home/manuel/workspaces/2026-03-02/deliver-mento-1/web-agent-example/pkg/discodialogue/timeline.go — Ported custom timeline projection to TimelineEntityV2
- /home/manuel/workspaces/2026-03-02/deliver-mento-1/web-agent-example/pkg/thinkingmode/sem.go — Replaced removed shared SEM proto usage with local JSON payloads
- /home/manuel/workspaces/2026-03-02/deliver-mento-1/web-agent-example/pkg/thinkingmode/timeline.go — Ported custom timeline projection to TimelineEntityV2


## 2026-03-09

Recorded the completed port in the implementation diary and linked the code commit 191a2713d7d009d038472767e72f6426c320106f after verifying go test ./... and go run ./cmd/web-agent-example serve --help under the workspace go.work.

### Related Files

- /home/manuel/workspaces/2026-03-02/deliver-mento-1/web-agent-example/ttmp/2026/03/08/WAE-0001--build-web-agent-example-against-updated-pinocchio-apis/reference/01-diary.md — Implementation diary entry for the completed port

