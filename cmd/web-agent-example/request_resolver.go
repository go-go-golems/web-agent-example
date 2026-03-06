package main

import (
	"encoding/json"
	"net/http"
	"strings"

	webhttp "github.com/go-go-golems/pinocchio/pkg/webchat/http"
	"github.com/google/uuid"
)

type noCookieRequestResolver struct{}

func newNoCookieRequestResolver() *noCookieRequestResolver {
	return &noCookieRequestResolver{}
}

func (r *noCookieRequestResolver) Resolve(req *http.Request) (webhttp.ResolvedConversationRequest, error) {
	if req == nil {
		return webhttp.ResolvedConversationRequest{}, &webhttp.RequestResolutionError{
			Status:    http.StatusBadRequest,
			ClientMsg: "bad request",
		}
	}

	switch req.Method {
	case http.MethodGet:
		return r.resolveWS(req)
	case http.MethodPost:
		return r.resolveChat(req)
	default:
		return webhttp.ResolvedConversationRequest{}, &webhttp.RequestResolutionError{
			Status:    http.StatusMethodNotAllowed,
			ClientMsg: "method not allowed",
		}
	}
}

func (r *noCookieRequestResolver) resolveWS(req *http.Request) (webhttp.ResolvedConversationRequest, error) {
	convID := strings.TrimSpace(req.URL.Query().Get("conv_id"))
	if convID == "" {
		return webhttp.ResolvedConversationRequest{}, &webhttp.RequestResolutionError{
			Status:    http.StatusBadRequest,
			ClientMsg: "missing conv_id",
		}
	}

	return webhttp.ResolvedConversationRequest{
		ConvID:             convID,
		RuntimeKey:         defaultRuntimeKey,
		RuntimeFingerprint: defaultRuntimeKey,
	}, nil
}

func (r *noCookieRequestResolver) resolveChat(req *http.Request) (webhttp.ResolvedConversationRequest, error) {
	var body webhttp.ChatRequestBody
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		return webhttp.ResolvedConversationRequest{}, &webhttp.RequestResolutionError{
			Status:    http.StatusBadRequest,
			ClientMsg: "bad request",
			Err:       err,
		}
	}
	if body.Prompt == "" && body.Text != "" {
		body.Prompt = body.Text
	}

	convID := strings.TrimSpace(body.ConvID)
	if convID == "" {
		convID = uuid.NewString()
	}

	return webhttp.ResolvedConversationRequest{
		ConvID:             convID,
		RuntimeKey:         defaultRuntimeKey,
		RuntimeFingerprint: defaultRuntimeKey,
		Overrides:          copyStringAnyMap(body.RequestOverrides),
		Prompt:             strings.TrimSpace(body.Prompt),
		IdempotencyKey:     strings.TrimSpace(body.IdempotencyKey),
	}, nil
}

func copyStringAnyMap(in map[string]any) map[string]any {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
