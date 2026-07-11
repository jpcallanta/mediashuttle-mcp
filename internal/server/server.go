package server

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"mediashuttle-mcp/internal/client"
)

// NewMCPServer returns an MCP server with all tools registered.
func NewMCPServer(apiKey string) *server.MCPServer {
	s := server.NewMCPServer(
		"mediashuttle-mcp",
		"0.1.0",
		server.WithToolCapabilities(true),
	)

	c := client.NewClient(apiKey)

	s.AddTools(
		// Portal tools
		server.ServerTool{
			Tool:    listPortalsTool(),
			Handler: listPortalsHandler(c),
		},
		server.ServerTool{
			Tool:    createPortalTool(),
			Handler: createPortalHandler(c),
		},
		server.ServerTool{
			Tool:    updatePortalTool(),
			Handler: updatePortalHandler(c),
		},

		// Portal user tools
		server.ServerTool{
			Tool:    listPortalUsersTool(),
			Handler: listPortalUsersHandler(c),
		},
		server.ServerTool{
			Tool:    getPortalUserTool(),
			Handler: getPortalUserHandler(c),
		},
		server.ServerTool{
			Tool:    addPortalUserTool(),
			Handler: addPortalUserHandler(c),
		},
		server.ServerTool{
			Tool:    updatePortalUserTool(),
			Handler: updatePortalUserHandler(c),
		},
		server.ServerTool{
			Tool:    removePortalUserTool(),
			Handler: removePortalUserHandler(c),
		},

		// Portal storage tools
		server.ServerTool{
			Tool:    listPortalStorageTool(),
			Handler: listPortalStorageHandler(c),
		},
		server.ServerTool{
			Tool:    assignPortalStorageTool(),
			Handler: assignPortalStorageHandler(c),
		},

		// Storage tools
		server.ServerTool{
			Tool:    listStorageTool(),
			Handler: listStorageHandler(c),
		},
		server.ServerTool{
			Tool:    getStorageTool(),
			Handler: getStorageHandler(c),
		},

		// Transfer tools
		server.ServerTool{
			Tool:    listTransfersTool(),
			Handler: listTransfersHandler(c),
		},
	)

	return s
}

// --- Portal Tools ---

func listPortalsTool() mcp.Tool {
	return mcp.NewTool("list_portals",
		mcp.WithDescription(
			"List all portals in the account.",
		),
		mcp.WithString("url",
			mcp.Description(
				"Filter by portal URL",
			),
		),
	)
}

func listPortalsHandler(
	c *client.Client,
) server.ToolHandlerFunc {
	return func(
		ctx context.Context,
		req mcp.CallToolRequest,
	) (*mcp.CallToolResult, error) {
		args, _ := req.Params.Arguments.(map[string]any)
		urlFilter := ""

		if v, ok := args["url"].(string); ok {
			urlFilter = v
		}

		portals, err := c.ListPortals(urlFilter)
		if err != nil {
			return mcp.NewToolResultError(
				err.Error(),
			), nil
		}

		return mcp.NewToolResultText(
			client.FormatJSON(portals),
		), nil
	}
}

func createPortalTool() mcp.Tool {
	return mcp.NewTool("create_portal",
		mcp.WithDescription(
			"Create a new portal.",
		),
		mcp.WithString("name",
			mcp.Description("Portal name"),
			mcp.Required(),
		),
		mcp.WithString("type",
			mcp.Description(
				"Portal type: send, share, or submit",
			),
			mcp.Required(),
			mcp.WithStringEnumItems(
				[]string{"send", "share", "submit"},
			),
		),
		mcp.WithString("url",
			mcp.Description(
				"URL prefix ending in"+
					" .mediashuttle.com",
			),
		),
		mcp.WithBoolean("media_shuttle_auth",
			mcp.Description(
				"Enable MS auth (default true)",
			),
		),
		mcp.WithBoolean("saml_auth",
			mcp.Description("Enable SAML auth"),
		),
		mcp.WithBoolean(
			"allow_unauthenticated_uploads",
			mcp.Description(
				"Allow uploads without login",
			),
		),
		mcp.WithBoolean(
			"allow_unauthenticated_downloads",
			mcp.Description(
				"Allow downloads without login",
			),
		),
		mcp.WithBoolean("notify_members",
			mcp.Description(
				"Notify members on add (default true)",
			),
		),
	)
}

func createPortalHandler(
	c *client.Client,
) server.ToolHandlerFunc {
	return func(
		ctx context.Context,
		req mcp.CallToolRequest,
	) (*mcp.CallToolResult, error) {
		args, _ := req.Params.Arguments.(map[string]any)
		portal := client.ParsePortalArgs(args)

		result, err := c.CreatePortal(portal)
		if err != nil {
			return mcp.NewToolResultError(
				err.Error(),
			), nil
		}

		return mcp.NewToolResultText(
			client.FormatJSON(result),
		), nil
	}
}

func updatePortalTool() mcp.Tool {
	return mcp.NewTool("update_portal",
		mcp.WithDescription(
			"Update an existing portal.",
		),
		mcp.WithString("portal_id",
			mcp.Description("Portal UUID"),
			mcp.Required(),
		),
		mcp.WithString("name",
			mcp.Description("New portal name"),
		),
		mcp.WithString("url",
			mcp.Description("New URL prefix"),
		),
		mcp.WithBoolean("media_shuttle_auth",
			mcp.Description("Enable MS auth"),
		),
		mcp.WithBoolean("saml_auth",
			mcp.Description("Enable SAML auth"),
		),
		mcp.WithBoolean(
			"allow_unauthenticated_uploads",
			mcp.Description(
				"Allow uploads without login",
			),
		),
		mcp.WithBoolean(
			"allow_unauthenticated_downloads",
			mcp.Description(
				"Allow downloads without login",
			),
		),
		mcp.WithBoolean("notify_members",
			mcp.Description("Notify members on add"),
		),
	)
}

func updatePortalHandler(
	c *client.Client,
) server.ToolHandlerFunc {
	return func(
		ctx context.Context,
		req mcp.CallToolRequest,
	) (*mcp.CallToolResult, error) {
		args, _ := req.Params.Arguments.(map[string]any)
		portalID, _ := args["portal_id"].(string)

		if portalID == "" {
			return mcp.NewToolResultError(
				"portal_id is required",
			), nil
		}

		update := client.UpdatePortal{}

		if v, ok := args["name"].(string); ok {
			update.Name = v
		}
		if v, ok := args["url"].(string); ok {
			update.URL = v
		}

		auth := &client.Authentication{}
		hasAuth := false

		if v, ok := args["media_shuttle_auth"].(bool); ok {
			auth.MediaShuttle = v
			hasAuth = true
		}
		if v, ok := args["saml_auth"].(bool); ok {
			auth.SAML = v
			hasAuth = true
		}

		if hasAuth {
			update.Authentication = auth
		}

		linkAuth := &client.LinkAuth{}
		hasLinkAuth := false

		k := "allow_unauthenticated_uploads"
		if v, ok := args[k].(bool); ok {
			linkAuth.AllowUnauthUploads = &v
			hasLinkAuth = true
		}

		k = "allow_unauthenticated_downloads"
		if v, ok := args[k].(bool); ok {
			linkAuth.AllowUnauthDownloads = &v
			hasLinkAuth = true
		}

		if hasLinkAuth {
			update.LinkAuth = linkAuth
		}

		if v, ok := args["notify_members"].(bool); ok {
			update.NotifyMembers = &v
		}

		result, err := c.UpdatePortal(portalID, update)
		if err != nil {
			return mcp.NewToolResultError(
				err.Error(),
			), nil
		}

		return mcp.NewToolResultText(
			client.FormatJSON(result),
		), nil
	}
}

// --- Portal User Tools ---

func listPortalUsersTool() mcp.Tool {
	return mcp.NewTool("list_portal_users",
		mcp.WithDescription(
			"List all members of a portal.",
		),
		mcp.WithString("portal_id",
			mcp.Description("Portal UUID"),
			mcp.Required(),
		),
	)
}

func listPortalUsersHandler(
	c *client.Client,
) server.ToolHandlerFunc {
	return func(
		ctx context.Context,
		req mcp.CallToolRequest,
	) (*mcp.CallToolResult, error) {
		args, _ := req.Params.Arguments.(map[string]any)
		portalID, _ := args["portal_id"].(string)

		if portalID == "" {
			return mcp.NewToolResultError(
				"portal_id is required",
			), nil
		}

		users, err := c.ListPortalUsers(portalID)
		if err != nil {
			return mcp.NewToolResultError(
				err.Error(),
			), nil
		}

		return mcp.NewToolResultText(
			client.FormatJSON(users),
		), nil
	}
}

func getPortalUserTool() mcp.Tool {
	return mcp.NewTool("get_portal_user",
		mcp.WithDescription(
			"Get portal member details.",
		),
		mcp.WithString("portal_id",
			mcp.Description("Portal UUID"),
			mcp.Required(),
		),
		mcp.WithString("email",
			mcp.Description("Member email"),
			mcp.Required(),
		),
	)
}

func getPortalUserHandler(
	c *client.Client,
) server.ToolHandlerFunc {
	return func(
		ctx context.Context,
		req mcp.CallToolRequest,
	) (*mcp.CallToolResult, error) {
		args, _ := req.Params.Arguments.(map[string]any)
		portalID, _ := args["portal_id"].(string)
		email, _ := args["email"].(string)

		if portalID == "" || email == "" {
			return mcp.NewToolResultError(
				"portal_id and email required",
			), nil
		}

		user, err := c.GetPortalUser(portalID, email)
		if err != nil {
			return mcp.NewToolResultError(
				err.Error(),
			), nil
		}

		return mcp.NewToolResultText(
			client.FormatJSON(user),
		), nil
	}
}

func addPortalUserTool() mcp.Tool {
	return mcp.NewTool("add_portal_user",
		mcp.WithDescription(
			"Add a user to a portal.",
		),
		mcp.WithString("portal_id",
			mcp.Description("Portal UUID"),
			mcp.Required(),
		),
		mcp.WithString("email",
			mcp.Description("User email"),
			mcp.Required(),
		),
		mcp.WithString("role",
			mcp.Description(
				"Role: Member or Ops",
			),
			mcp.WithStringEnumItems(
				[]string{"Member", "Ops"},
			),
		),
		mcp.WithString("expires_on",
			mcp.Description(
				"ISO 8601 expiry date",
			),
		),
		mcp.WithBoolean("can_send_to_member",
			mcp.Description(
				"Allow sending to members",
			),
		),
		mcp.WithBoolean("can_send_to_non_member",
			mcp.Description(
				"Allow sending to non-members",
			),
		),
		mcp.WithBoolean("can_receive",
			mcp.Description(
				"Allow receiving content",
			),
		),
		mcp.WithBoolean("can_send_from_share",
			mcp.Description(
				"Allow sending from share portal",
			),
		),
		mcp.WithBoolean("can_deliver_automatically",
			mcp.Description(
				"Allow auto-delivery",
			),
		),
		mcp.WithBoolean("can_submit",
			mcp.Description(
				"Allow submitting content",
			),
		),
	)
}

func addPortalUserHandler(
	c *client.Client,
) server.ToolHandlerFunc {
	return func(
		ctx context.Context,
		req mcp.CallToolRequest,
	) (*mcp.CallToolResult, error) {
		args, _ := req.Params.Arguments.(map[string]any)
		portalID, _ := args["portal_id"].(string)

		if portalID == "" {
			return mcp.NewToolResultError(
				"portal_id is required",
			), nil
		}

		member := client.ParseMemberArgs(args)

		if member.Email == "" {
			return mcp.NewToolResultError(
				"email is required",
			), nil
		}

		if err := c.AddPortalUser(
			portalID, member,
		); err != nil {
			return mcp.NewToolResultError(
				err.Error(),
			), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf(
			"User %s added to portal %s",
			member.Email, portalID,
		)), nil
	}
}

func updatePortalUserTool() mcp.Tool {
	return mcp.NewTool("update_portal_user",
		mcp.WithDescription(
			"Update a portal member.",
		),
		mcp.WithString("portal_id",
			mcp.Description("Portal UUID"),
			mcp.Required(),
		),
		mcp.WithString("email",
			mcp.Description("Member email"),
			mcp.Required(),
		),
		mcp.WithString("role",
			mcp.Description("New role: Member or Ops"),
			mcp.WithStringEnumItems(
				[]string{"Member", "Ops"},
			),
		),
		mcp.WithString("expires_on",
			mcp.Description("New expiry (ISO 8601)"),
		),
		mcp.WithBoolean("can_send_to_member",
			mcp.Description("Allow sending to members"),
		),
		mcp.WithBoolean("can_send_to_non_member",
			mcp.Description(
				"Allow sending to non-members",
			),
		),
		mcp.WithBoolean("can_receive",
			mcp.Description("Allow receiving content"),
		),
		mcp.WithBoolean("can_send_from_share",
			mcp.Description(
				"Allow sending from share portal",
			),
		),
		mcp.WithBoolean("can_deliver_automatically",
			mcp.Description("Allow auto-delivery"),
		),
		mcp.WithBoolean("can_submit",
			mcp.Description("Allow submitting content"),
		),
	)
}

func updatePortalUserHandler(
	c *client.Client,
) server.ToolHandlerFunc {
	return func(
		ctx context.Context,
		req mcp.CallToolRequest,
	) (*mcp.CallToolResult, error) {
		args, _ := req.Params.Arguments.(map[string]any)
		portalID, _ := args["portal_id"].(string)
		email, _ := args["email"].(string)

		if portalID == "" || email == "" {
			return mcp.NewToolResultError(
				"portal_id and email required",
			), nil
		}

		member := client.ParseMemberArgs(args)

		if err := c.UpdatePortalUser(
			portalID, email, member,
		); err != nil {
			return mcp.NewToolResultError(
				err.Error(),
			), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf(
			"User %s updated in portal %s",
			email, portalID,
		)), nil
	}
}

func removePortalUserTool() mcp.Tool {
	return mcp.NewTool("remove_portal_user",
		mcp.WithDescription(
			"Remove a user from a portal.",
		),
		mcp.WithString("portal_id",
			mcp.Description("Portal UUID"),
			mcp.Required(),
		),
		mcp.WithString("email",
			mcp.Description("Member email"),
			mcp.Required(),
		),
	)
}

func removePortalUserHandler(
	c *client.Client,
) server.ToolHandlerFunc {
	return func(
		ctx context.Context,
		req mcp.CallToolRequest,
	) (*mcp.CallToolResult, error) {
		args, _ := req.Params.Arguments.(map[string]any)
		portalID, _ := args["portal_id"].(string)
		email, _ := args["email"].(string)

		if portalID == "" || email == "" {
			return mcp.NewToolResultError(
				"portal_id and email required",
			), nil
		}

		if err := c.RemovePortalUser(
			portalID, email,
		); err != nil {
			return mcp.NewToolResultError(
				err.Error(),
			), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf(
			"User %s removed from portal %s",
			email, portalID,
		)), nil
	}
}

// --- Portal Storage Tools ---

func listPortalStorageTool() mcp.Tool {
	return mcp.NewTool("list_portal_storage",
		mcp.WithDescription(
			"List storage assigned to a portal.",
		),
		mcp.WithString("portal_id",
			mcp.Description("Portal UUID"),
			mcp.Required(),
		),
	)
}

func listPortalStorageHandler(
	c *client.Client,
) server.ToolHandlerFunc {
	return func(
		ctx context.Context,
		req mcp.CallToolRequest,
	) (*mcp.CallToolResult, error) {
		args, _ := req.Params.Arguments.(map[string]any)
		portalID, _ := args["portal_id"].(string)

		if portalID == "" {
			return mcp.NewToolResultError(
				"portal_id is required",
			), nil
		}

		storage, err := c.ListPortalStorage(portalID)
		if err != nil {
			return mcp.NewToolResultError(
				err.Error(),
			), nil
		}

		return mcp.NewToolResultText(
			client.FormatJSON(storage),
		), nil
	}
}

func assignPortalStorageTool() mcp.Tool {
	return mcp.NewTool("assign_portal_storage",
		mcp.WithDescription(
			"Assign storage to a portal.",
		),
		mcp.WithString("portal_id",
			mcp.Description("Portal UUID"),
			mcp.Required(),
		),
		mcp.WithString("storage_id",
			mcp.Description("Storage UUID"),
			mcp.Required(),
		),
		mcp.WithString("repository_path",
			mcp.Description("Remote repo path"),
		),
	)
}

func assignPortalStorageHandler(
	c *client.Client,
) server.ToolHandlerFunc {
	return func(
		ctx context.Context,
		req mcp.CallToolRequest,
	) (*mcp.CallToolResult, error) {
		args, _ := req.Params.Arguments.(map[string]any)
		portalID, _ := args["portal_id"].(string)
		storageID, _ := args["storage_id"].(string)

		if portalID == "" || storageID == "" {
			return mcp.NewToolResultError(
				"portal_id and storage_id required",
			), nil
		}

		ps := client.PortalStorage{
			StorageID: storageID,
		}

		if v, ok := args["repository_path"].(string); ok {
			ps.RepositoryPath = v
		}

		if err := c.AssignPortalStorage(
			portalID, storageID, ps,
		); err != nil {
			return mcp.NewToolResultError(
				err.Error(),
			), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf(
			"Storage %s assigned to portal %s",
			storageID, portalID,
		)), nil
	}
}

// --- Storage Tools ---

func listStorageTool() mcp.Tool {
	return mcp.NewTool("list_storage",
		mcp.WithDescription(
			"List all storage locations.",
		),
	)
}

func listStorageHandler(
	c *client.Client,
) server.ToolHandlerFunc {
	return func(
		ctx context.Context,
		req mcp.CallToolRequest,
	) (*mcp.CallToolResult, error) {
		storage, err := c.ListStorage()
		if err != nil {
			return mcp.NewToolResultError(
				err.Error(),
			), nil
		}

		return mcp.NewToolResultText(
			client.FormatJSON(storage),
		), nil
	}
}

func getStorageTool() mcp.Tool {
	return mcp.NewTool("get_storage",
		mcp.WithDescription(
			"Get storage location details.",
		),
		mcp.WithString("storage_id",
			mcp.Description("Storage UUID"),
			mcp.Required(),
		),
	)
}

func getStorageHandler(
	c *client.Client,
) server.ToolHandlerFunc {
	return func(
		ctx context.Context,
		req mcp.CallToolRequest,
	) (*mcp.CallToolResult, error) {
		args, _ := req.Params.Arguments.(map[string]any)
		storageID, _ := args["storage_id"].(string)

		if storageID == "" {
			return mcp.NewToolResultError(
				"storage_id is required",
			), nil
		}

		storage, err := c.GetStorage(storageID)
		if err != nil {
			return mcp.NewToolResultError(
				err.Error(),
			), nil
		}

		return mcp.NewToolResultText(
			client.FormatJSON(storage),
		), nil
	}
}

// --- Transfer Tools ---

func listTransfersTool() mcp.Tool {
	return mcp.NewTool("list_transfers",
		mcp.WithDescription(
			"List active transfers.",
		),
		mcp.WithString("portal_id",
			mcp.Description(
				"Optional portal UUID filter",
			),
		),
	)
}

func listTransfersHandler(
	c *client.Client,
) server.ToolHandlerFunc {
	return func(
		ctx context.Context,
		req mcp.CallToolRequest,
	) (*mcp.CallToolResult, error) {
		args, _ := req.Params.Arguments.(map[string]any)
		portalID, _ := args["portal_id"].(string)

		transfers, err := c.ListTransfers(portalID)
		if err != nil {
			return mcp.NewToolResultError(
				err.Error(),
			), nil
		}

		return mcp.NewToolResultText(
			client.FormatJSON(transfers),
		), nil
	}
}
