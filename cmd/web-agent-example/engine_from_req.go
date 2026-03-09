package main

import (
	"encoding/json"
	"net/http"

	gepprofiles "github.com/go-go-golems/geppetto/pkg/profiles"
	"github.com/google/uuid"

	webhttp "github.com/go-go-golems/pinocchio/pkg/webchat/http"
)

const (
	defaultRuntimeKey         = "default"
	defaultRuntimeFingerprint = "web-agent-example/default"
	defaultProfileVersion     = 1
)

type staticRequestResolver struct {
	runtime *gepprofiles.RuntimeSpec
}

func newStaticRequestResolver() *staticRequestResolver {
	return &staticRequestResolver{runtime: defaultRuntimeSpec()}
}

func (r *staticRequestResolver) Resolve(req *http.Request) (webhttp.ResolvedConversationRequest, error) {
	if req == nil {
		return webhttp.ResolvedConversationRequest{}, &webhttp.RequestResolutionError{Status: http.StatusBadRequest, ClientMsg: "bad request"}
	}

	switch req.Method {
	case http.MethodGet:
		return r.resolveWS(req)
	case http.MethodPost:
		return r.resolveChat(req)
	default:
		return webhttp.ResolvedConversationRequest{}, &webhttp.RequestResolutionError{Status: http.StatusMethodNotAllowed, ClientMsg: "method not allowed"}
	}
}

func (r *staticRequestResolver) resolveWS(req *http.Request) (webhttp.ResolvedConversationRequest, error) {
	convID := req.URL.Query().Get("conv_id")
	if convID == "" {
		return webhttp.ResolvedConversationRequest{}, &webhttp.RequestResolutionError{Status: http.StatusBadRequest, ClientMsg: "missing conv_id"}
	}

	return webhttp.ResolvedConversationRequest{
		ConvID:             convID,
		RuntimeKey:         defaultRuntimeKey,
		RuntimeFingerprint: defaultRuntimeFingerprint,
		ProfileVersion:     defaultProfileVersion,
		ResolvedRuntime:    cloneRuntimeSpec(r.runtime),
	}, nil
}

func (r *staticRequestResolver) resolveChat(req *http.Request) (webhttp.ResolvedConversationRequest, error) {
	var body webhttp.ChatRequestBody
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		return webhttp.ResolvedConversationRequest{}, &webhttp.RequestResolutionError{Status: http.StatusBadRequest, ClientMsg: "bad request", Err: err}
	}
	if body.Prompt == "" && body.Text != "" {
		body.Prompt = body.Text
	}
	if body.Prompt == "" {
		return webhttp.ResolvedConversationRequest{}, &webhttp.RequestResolutionError{Status: http.StatusBadRequest, ClientMsg: "missing prompt"}
	}

	convID := body.ConvID
	if convID == "" {
		convID = uuid.NewString()
	}

	return webhttp.ResolvedConversationRequest{
		ConvID:             convID,
		RuntimeKey:         defaultRuntimeKey,
		RuntimeFingerprint: defaultRuntimeFingerprint,
		ProfileVersion:     defaultProfileVersion,
		ResolvedRuntime:    cloneRuntimeSpec(r.runtime),
		Overrides:          cloneMap(body.RequestOverrides),
		Prompt:             body.Prompt,
		IdempotencyKey:     body.IdempotencyKey,
	}, nil
}
