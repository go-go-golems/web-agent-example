package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/go-go-golems/pinocchio/pkg/webchat"
)

type noCookieEngineFromReqBuilder struct {
}

func newNoCookieEngineFromReqBuilder() *noCookieEngineFromReqBuilder {
	return &noCookieEngineFromReqBuilder{}
}

func (b *noCookieEngineFromReqBuilder) BuildEngineFromReq(req *http.Request) (webchat.EngineBuildInput, *webchat.ChatRequestBody, error) {
	if req == nil {
		return webchat.EngineBuildInput{}, nil, &webchat.RequestBuildError{Status: http.StatusBadRequest, ClientMsg: "bad request"}
	}

	switch req.Method {
	case http.MethodGet:
		in, err := b.buildFromWSReq(req)
		return in, nil, err
	case http.MethodPost:
		return b.buildFromChatReq(req)
	default:
		return webchat.EngineBuildInput{}, nil, &webchat.RequestBuildError{Status: http.StatusMethodNotAllowed, ClientMsg: "method not allowed"}
	}
}

func (b *noCookieEngineFromReqBuilder) buildFromWSReq(req *http.Request) (webchat.EngineBuildInput, error) {
	convID := strings.TrimSpace(req.URL.Query().Get("conv_id"))
	if convID == "" {
		return webchat.EngineBuildInput{}, &webchat.RequestBuildError{Status: http.StatusBadRequest, ClientMsg: "missing conv_id"}
	}

	_ = strings.TrimSpace(req.URL.Query().Get("profile"))
	return webchat.EngineBuildInput{ConvID: convID, ProfileSlug: "default"}, nil
}

func (b *noCookieEngineFromReqBuilder) buildFromChatReq(req *http.Request) (webchat.EngineBuildInput, *webchat.ChatRequestBody, error) {
	var body webchat.ChatRequestBody
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		return webchat.EngineBuildInput{}, nil, &webchat.RequestBuildError{Status: http.StatusBadRequest, ClientMsg: "bad request", Err: err}
	}
	if body.Prompt == "" && body.Text != "" {
		body.Prompt = body.Text
	}

	convID := strings.TrimSpace(body.ConvID)
	if convID == "" {
		convID = uuid.NewString()
		body.ConvID = convID
	}

	_ = strings.TrimSpace(profileSlugFromChatRequest(req))
	return webchat.EngineBuildInput{ConvID: convID, ProfileSlug: "default", Overrides: body.Overrides}, &body, nil
}

// profileSlugFromChatRequest replicates the webchat behavior but keeps it local to this package.
func profileSlugFromChatRequest(req *http.Request) string {
	if req == nil {
		return ""
	}
	path := req.URL.Path
	if path == "" {
		return ""
	}
	if idx := strings.Index(path, "/chat/"); idx >= 0 {
		rest := path[idx+len("/chat/"):]
		if rest == "" {
			return ""
		}
		if i := strings.Index(rest, "/"); i >= 0 {
			rest = rest[:i]
		}
		return strings.TrimSpace(rest)
	}
	return ""
}
