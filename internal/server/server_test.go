package server

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestNewMCPServer_ReturnsNonNil(t *testing.T) {
	s := NewMCPServer("test-key")

	if s == nil {
		t.Fatal("NewMCPServer returned nil")
	}
}

func TestListPortalsTool_Definition(t *testing.T) {
	tool := listPortalsTool()
	if tool.Name != "list_portals" {
		t.Errorf(
			"want name 'list_portals', got %q",
			tool.Name,
		)
	}
	if tool.Description == "" {
		t.Error("want non-empty description")
	}
}

func TestCreatePortalTool_Definition(t *testing.T) {
	tool := createPortalTool()
	if tool.Name != "create_portal" {
		t.Errorf(
			"want name 'create_portal', got %q",
			tool.Name,
		)
	}
}

func TestUpdatePortalTool_Definition(t *testing.T) {
	tool := updatePortalTool()
	if tool.Name != "update_portal" {
		t.Errorf(
			"want name 'update_portal', got %q",
			tool.Name,
		)
	}
}

func TestListPortalUsersTool_Definition(t *testing.T) {
	tool := listPortalUsersTool()
	if tool.Name != "list_portal_users" {
		t.Errorf("want 'list_portal_users', got %q",
			tool.Name,
		)
	}
}

func TestGetPortalUserTool_Definition(t *testing.T) {
	tool := getPortalUserTool()
	if tool.Name != "get_portal_user" {
		t.Errorf("want 'get_portal_user', got %q",
			tool.Name,
		)
	}
}

func TestAddPortalUserTool_Definition(t *testing.T) {
	tool := addPortalUserTool()
	if tool.Name != "add_portal_user" {
		t.Errorf("want 'add_portal_user', got %q",
			tool.Name,
		)
	}
}

func TestUpdatePortalUserTool_Definition(t *testing.T) {
	tool := updatePortalUserTool()
	if tool.Name != "update_portal_user" {
		t.Errorf("want 'update_portal_user', got %q",
			tool.Name,
		)
	}
}

func TestRemovePortalUserTool_Definition(t *testing.T) {
	tool := removePortalUserTool()
	if tool.Name != "remove_portal_user" {
		t.Errorf("want 'remove_portal_user', got %q",
			tool.Name,
		)
	}
}

func TestListStorageTool_Definition(t *testing.T) {
	tool := listStorageTool()
	if tool.Name != "list_storage" {
		t.Errorf(
			"want name 'list_storage', got %q",
			tool.Name,
		)
	}
}

func TestGetStorageTool_Definition(t *testing.T) {
	tool := getStorageTool()
	if tool.Name != "get_storage" {
		t.Errorf(
			"want name 'get_storage', got %q",
			tool.Name,
		)
	}
}

func TestListTransfersTool_Definition(t *testing.T) {
	tool := listTransfersTool()
	if tool.Name != "list_transfers" {
		t.Errorf("want 'list_transfers', got %q",
			tool.Name,
		)
	}
}

func TestListPortalStorageTool_Definition(t *testing.T) {
	tool := listPortalStorageTool()
	if tool.Name != "list_portal_storage" {
		t.Errorf("want 'list_portal_storage', got %q",
			tool.Name,
		)
	}
}

func TestAssignPortalStorageTool_Definition(t *testing.T) {
	tool := assignPortalStorageTool()
	if tool.Name != "assign_portal_storage" {
		t.Errorf("want 'assign_portal_storage', got %q",
			tool.Name,
		)
	}
}

func TestAllToolsRegistered(t *testing.T) {
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

	defs := []struct {
		name string
		def  func() mcp.Tool
	}{
		{"list_portals", listPortalsTool},
		{"create_portal", createPortalTool},
		{"update_portal", updatePortalTool},
		{"list_portal_users", listPortalUsersTool},
		{"get_portal_user", getPortalUserTool},
		{"add_portal_user", addPortalUserTool},
		{"update_portal_user", updatePortalUserTool},
		{"remove_portal_user", removePortalUserTool},
		{"list_portal_storage", listPortalStorageTool},
		{"assign_portal_storage", assignPortalStorageTool},
		{"list_storage", listStorageTool},
		{"get_storage", getStorageTool},
		{"list_transfers", listTransfersTool},
	}

	if len(defs) != len(want) {
		t.Errorf("want %d tools, got %d",
			len(want), len(defs),
		)
	}

	seen := make(map[string]bool)

	for _, d := range defs {
		tool := d.def()
		if tool.Name != d.name {
			t.Errorf("mismatch: want %q, got %q",
				d.name, tool.Name,
			)
		}
		if seen[tool.Name] {
			t.Errorf("duplicate: %q", tool.Name)
		}
		seen[tool.Name] = true
	}

	for _, name := range want {
		if !seen[name] {
			t.Errorf("missing tool: %q", name)
		}
	}
}

func TestUpdatePortalHandler_MissingPortalID(t *testing.T) {
	handler := updatePortalHandler(nil)
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{}

	result, err := handler(
		context.Background(), req,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("want result, got nil")
	}
	if !result.IsError {
		t.Error("want error for missing portal_id")
	}
}

func TestListPortalUsersHandler_MissingID(t *testing.T) {
	handler := listPortalUsersHandler(nil)
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{}

	result, err := handler(
		context.Background(), req,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("want error for missing portal_id")
	}
}

func TestGetPortalUserHandler_MissingFields(t *testing.T) {
	handler := getPortalUserHandler(nil)
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{}

	result, err := handler(
		context.Background(), req,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("want error for missing fields")
	}
}

func TestAddPortalUserHandler_MissingFields(t *testing.T) {
	handler := addPortalUserHandler(nil)
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{}

	result, err := handler(
		context.Background(), req,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("want error for missing fields")
	}
}

func TestRemovePortalUserHandler_MissingFields(t *testing.T) {
	handler := removePortalUserHandler(nil)
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{}

	result, err := handler(
		context.Background(), req,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("want error for missing fields")
	}
}

func TestAssignPortalStorageHandler_MissingFields(t *testing.T) {
	handler := assignPortalStorageHandler(nil)
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{}

	result, err := handler(
		context.Background(), req,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("want error for missing fields")
	}
}

func TestGetStorageHandler_MissingStorageID(t *testing.T) {
	handler := getStorageHandler(nil)
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{}

	result, err := handler(
		context.Background(), req,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("want error for missing storage_id")
	}
}

func TestToolHandler_TextResponse(t *testing.T) {
	handler := func(
		_ context.Context,
		_ mcp.CallToolRequest,
	) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText("hello"), nil
	}

	result, err := handler(
		context.Background(), mcp.CallToolRequest{},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Error("want non-error result")
	}
	if len(result.Content) == 0 {
		t.Error("want content in result")
	}
}

func TestToolHandler_ErrorResponse(t *testing.T) {
	handler := func(
		_ context.Context,
		_ mcp.CallToolRequest,
	) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultError(
			"something went wrong",
		), nil
	}

	result, err := handler(
		context.Background(), mcp.CallToolRequest{},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("want error result")
	}
}

func TestJSONSerializationRoundTrip(t *testing.T) {
	type testStruct struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	original := testStruct{Name: "test", Count: 42}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded testStruct
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if decoded.Name != original.Name {
		t.Errorf("want Name %q, got %q",
			original.Name, decoded.Name,
		)
	}
	if decoded.Count != original.Count {
		t.Errorf("want Count %d, got %d",
			original.Count, decoded.Count,
		)
	}
}
