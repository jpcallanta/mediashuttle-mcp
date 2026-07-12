package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mcphttpserver "github.com/mark3labs/mcp-go/server"
)

const testSessionHeader = "Mcp-Session-Id"

func sendMCPPost(
	t *testing.T,
	ts *httptest.Server,
	sessionID string,
	body map[string]any,
) (*http.Response, string) {
	t.Helper()

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	req, err := http.NewRequest(
		"POST", ts.URL+"/mcp",
		bytes.NewReader(data),
	)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if sessionID != "" {
		req.Header.Set(testSessionHeader, sessionID)
	}

	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer resp.Body.Close()

	var result json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(
		&result,
	); err != nil {
		return resp, ""
	}

	return resp, string(result)
}

func initSession(
	t *testing.T,
	ts *httptest.Server,
) string {
	t.Helper()

	resp, body := sendMCPPost(t, ts, "", map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params": map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]any{},
			"clientInfo": map[string]any{
				"name":    "test",
				"version": "1.0",
			},
		},
	})
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("init failed: %d: %s",
			resp.StatusCode, body,
		)
	}

	sessionID := resp.Header.Get(testSessionHeader)
	if sessionID == "" {
		t.Fatal("want session ID in response")
	}

	_, _ = sendMCPPost(t, ts, sessionID, map[string]any{
		"jsonrpc": "2.0",
		"method":  "notifications/initialized",
	})

	return sessionID
}

func TestStreamableHTTPServer_Initialize(t *testing.T) {
	s := NewMCPServer("test-key")
	httpSrv := mcphttpserver.NewStreamableHTTPServer(s)
	ts := httptest.NewServer(httpSrv)
	defer ts.Close()

	resp, body := sendMCPPost(t, ts, "", map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params": map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]any{},
			"clientInfo": map[string]any{
				"name":    "test",
				"version": "1.0",
			},
		},
	})

	if resp.StatusCode != http.StatusOK {
		t.Fatalf(
			"want 200, got %d: %s",
			resp.StatusCode, body,
		)
	}

	sessionID := resp.Header.Get(testSessionHeader)
	if sessionID == "" {
		t.Error("want session ID header")
	}

	var initResp struct {
		Result struct {
			ProtocolVersion string `json:"protocolVersion"`
			ServerInfo      struct {
				Name    string `json:"name"`
				Version string `json:"version"`
			} `json:"serverInfo"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(body), &initResp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if initResp.Result.ProtocolVersion != "2024-11-05" {
		t.Errorf(
			"want 2024-11-05, got %s",
			initResp.Result.ProtocolVersion,
		)
	}
	if initResp.Result.ServerInfo.Name != "mediashuttle-mcp" {
		t.Errorf(
			"want mediashuttle-mcp, got %s",
			initResp.Result.ServerInfo.Name,
		)
	}
	if initResp.Result.ServerInfo.Version != "0.1.0" {
		t.Errorf(
			"want 0.1.0, got %s",
			initResp.Result.ServerInfo.Version,
		)
	}
}

func TestStreamableHTTPServer_ListTools(t *testing.T) {
	s := NewMCPServer("test-key")
	httpSrv := mcphttpserver.NewStreamableHTTPServer(s)
	ts := httptest.NewServer(httpSrv)
	defer ts.Close()

	sessionID := initSession(t, ts)

	resp, body := sendMCPPost(t, ts, sessionID, map[string]any{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "tools/list",
	})
	if resp.StatusCode != http.StatusOK {
		t.Fatalf(
			"tools/list: %d: %s",
			resp.StatusCode, body,
		)
	}

	var listResp struct {
		Result struct {
			Tools []struct {
				Name        string `json:"name"`
				Description string `json:"description"`
			} `json:"tools"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(body), &listResp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	want := []string{
		"list_portals",
		"create_portal",
		"update_portal",
		"list_portal_users",
		"get_portal_user",
		"add_portal_user",
		"update_portal_user",
		"remove_portal_user",
		"list_portal_storage",
		"assign_portal_storage",
		"list_storage",
		"get_storage",
		"list_transfers",
	}

	if len(listResp.Result.Tools) != len(want) {
		t.Errorf(
			"want %d tools, got %d",
			len(want),
			len(listResp.Result.Tools),
		)
	}

	seen := make(map[string]bool)
	for _, tool := range listResp.Result.Tools {
		if tool.Description == "" {
			t.Errorf(
				"tool %q: empty description",
				tool.Name,
			)
		}
		seen[tool.Name] = true
	}
	for _, name := range want {
		if !seen[name] {
			t.Errorf("missing tool: %q", name)
		}
	}
}

func TestStreamableHTTPServer_CallTool_MissingArgs(t *testing.T) {
	s := NewMCPServer("test-key")
	httpSrv := mcphttpserver.NewStreamableHTTPServer(s)
	ts := httptest.NewServer(httpSrv)
	defer ts.Close()

	sessionID := initSession(t, ts)

	resp, body := sendMCPPost(t, ts, sessionID, map[string]any{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "tools/call",
		"params": map[string]any{
			"name":      "list_portal_users",
			"arguments": map[string]any{},
		},
	})
	if resp.StatusCode != http.StatusOK {
		t.Fatalf(
			"tools/call: %d: %s",
			resp.StatusCode, body,
		)
	}

	var callResp struct {
		Result struct {
			IsError bool `json:"isError"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(body), &callResp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if !callResp.Result.IsError {
		t.Error("want error for missing portal_id")
	}
}

func TestStreamableHTTPServer_PutMethod(t *testing.T) {
	s := NewMCPServer("test-key")
	httpSrv := mcphttpserver.NewStreamableHTTPServer(s)
	ts := httptest.NewServer(httpSrv)
	defer ts.Close()

	req, err := http.NewRequest(
		"PUT", ts.URL+"/mcp", nil,
	)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("do: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("want 404, got %d", resp.StatusCode)
	}
}

func TestStreamableHTTPServer_BadContentType(t *testing.T) {
	s := NewMCPServer("test-key")
	httpSrv := mcphttpserver.NewStreamableHTTPServer(s)
	ts := httptest.NewServer(httpSrv)
	defer ts.Close()

	req, err := http.NewRequest(
		"POST", ts.URL+"/mcp",
		bytes.NewReader([]byte(`{}`)),
	)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("do: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("want 400, got %d", resp.StatusCode)
	}
}

func TestStreamableHTTPServer_MultipleSessions(t *testing.T) {
	s := NewMCPServer("test-key")
	httpSrv := mcphttpserver.NewStreamableHTTPServer(s)
	ts := httptest.NewServer(httpSrv)
	defer ts.Close()

	sid1 := initSession(t, ts)
	sid2 := initSession(t, ts)

	if sid1 == sid2 {
		t.Error("want different session IDs")
	}

	listCount := func(sessionID string) int {
		t.Helper()

		resp, body := sendMCPPost(t, ts, sessionID, map[string]any{
			"jsonrpc": "2.0",
			"id":      2,
			"method":  "tools/list",
		})
		if resp.StatusCode != http.StatusOK {
			t.Fatalf(
				"tools/list: %d: %s",
				resp.StatusCode, body,
			)
		}
		var listResp struct {
			Result struct {
				Tools []any `json:"tools"`
			} `json:"result"`
		}
		if err := json.Unmarshal(
			[]byte(body), &listResp,
		); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		return len(listResp.Result.Tools)
	}

	if n := listCount(sid1); n != 13 {
		t.Errorf("session 1: want 13, got %d", n)
	}
	if n := listCount(sid2); n != 13 {
		t.Errorf("session 2: want 13, got %d", n)
	}
}

func TestStreamableHTTPServer_NotFoundWithoutInit(t *testing.T) {
	s := NewMCPServer("test-key")
	httpSrv := mcphttpserver.NewStreamableHTTPServer(s)
	ts := httptest.NewServer(httpSrv)
	defer ts.Close()

	resp, body := sendMCPPost(t, ts, "nonexistent-session", map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/list",
	})
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf(
			"want 404, got %d: %s",
			resp.StatusCode, body,
		)
	}
}

func TestStreamableHTTPServer_GetMethodRejectedWithoutSession(t *testing.T) {
	s := NewMCPServer("test-key")
	httpSrv := mcphttpserver.NewStreamableHTTPServer(s)
	ts := httptest.NewServer(httpSrv)
	defer ts.Close()

	resp, err := ts.Client().Get(ts.URL + "/mcp")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf(
			"want 200 (SSE stream), got %d",
			resp.StatusCode,
		)
	}

	if resp.Header.Get("Content-Type") !=
		"text/event-stream" {
		t.Errorf(
			"want text/event-stream, got %s",
			resp.Header.Get("Content-Type"),
		)
	}
}
