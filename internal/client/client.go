package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const DefaultBaseURL = "https://api.mediashuttle.com/v1"

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient returns a Client configured for the given API key.
func NewClient(apiKey string) *Client {
	return &Client{
		baseURL:    DefaultBaseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

func (c *Client) doRequest(
	method, path string, body any,
) ([]byte, error) {
	var bodyReader io.Reader

	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf(
				"marshal request body: %w", err,
			)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(
		method, c.baseURL+path, bodyReader,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create request: %w", err,
		)
	}

	req.Header.Set("Authorization", c.apiKey)
	req.Header.Set("Accept", "application/json")

	if body != nil {
		req.Header.Set(
			"Content-Type", "application/json",
		)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(
			"execute request: %w", err,
		)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(
			"read response body: %w", err,
		)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr APIError
		if json.Unmarshal(data, &apiErr) == nil &&
			apiErr.Message != "" {
			return nil, fmt.Errorf(
				"API error %d: %s",
				resp.StatusCode, apiErr.Message,
			)
		}
		return nil, fmt.Errorf(
			"API error %d: %s",
			resp.StatusCode, string(data),
		)
	}

	return data, nil
}

// ListPortals returns all portals, optionally filtered by URL.
func (c *Client) ListPortals(
	urlFilter string,
) (*PortalList, error) {
	path := "/portals"

	if urlFilter != "" {
		path += "?url=" + url.QueryEscape(urlFilter)
	}

	data, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result PortalList
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf(
			"decode response: %w", err,
		)
	}

	return &result, nil
}

// CreatePortal creates a new portal and returns it.
func (c *Client) CreatePortal(
	portal Portal,
) (*Portal, error) {
	data, err := c.doRequest("POST", "/portals", portal)
	if err != nil {
		return nil, err
	}

	var result Portal
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf(
			"decode response: %w", err,
		)
	}

	return &result, nil
}

// UpdatePortal patches a portal and returns the result.
func (c *Client) UpdatePortal(
	portalID string, update UpdatePortal,
) (*Portal, error) {
	data, err := c.doRequest(
		"PATCH", "/portals/"+portalID, update,
	)
	if err != nil {
		return nil, err
	}

	var result Portal
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf(
			"decode response: %w", err,
		)
	}

	return &result, nil
}

// ListPortalUsers returns all members of a portal.
func (c *Client) ListPortalUsers(
	portalID string,
) (*PortalMemberResponse, error) {
	data, err := c.doRequest(
		"GET", "/portals/"+portalID+"/users", nil,
	)
	if err != nil {
		return nil, err
	}

	var result PortalMemberResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf(
			"decode response: %w", err,
		)
	}

	return &result, nil
}

// GetPortalUser returns a single portal member.
func (c *Client) GetPortalUser(
	portalID, email string,
) (*ResponseForPortalMember, error) {
	path := "/portals/" + portalID +
		"/users/" + url.PathEscape(email)

	data, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result ResponseForPortalMember
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf(
			"decode response: %w", err,
		)
	}

	return &result, nil
}

// AddPortalUser adds a user to a portal.
func (c *Client) AddPortalUser(
	portalID string, member PortalMember,
) error {
	path := "/portals/" + portalID + "/users"
	_, err := c.doRequest("POST", path, member)

	return err
}

// UpdatePortalUser updates a portal member's role or permissions.
func (c *Client) UpdatePortalUser(
	portalID, email string, member PortalMember,
) error {
	path := "/portals/" + portalID +
		"/users/" + url.PathEscape(email)
	_, err := c.doRequest("PUT", path, member)

	return err
}

// RemovePortalUser removes a user from a portal.
func (c *Client) RemovePortalUser(
	portalID, email string,
) error {
	path := "/portals/" + portalID +
		"/users/" + url.PathEscape(email)
	_, err := c.doRequest("DELETE", path, nil)

	return err
}

// ListPortalStorage returns storage assigned to a portal.
func (c *Client) ListPortalStorage(
	portalID string,
) (*PortalStorageList, error) {
	path := "/portals/" + portalID + "/storage"

	data, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result PortalStorageList
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf(
			"decode response: %w", err,
		)
	}

	return &result, nil
}

// AssignPortalStorage assigns storage to a portal.
func (c *Client) AssignPortalStorage(
	portalID, storageID string, ps PortalStorage,
) error {
	path := "/portals/" + portalID +
		"/storage/" + storageID
	_, err := c.doRequest("PUT", path, ps)

	return err
}

// ListStorage returns all storage locations.
func (c *Client) ListStorage() (*StorageList, error) {
	data, err := c.doRequest("GET", "/storage", nil)
	if err != nil {
		return nil, err
	}

	var result StorageList
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf(
			"decode response: %w", err,
		)
	}

	return &result, nil
}

// GetStorage returns a single storage location.
func (c *Client) GetStorage(
	storageID string,
) (*Storage, error) {
	data, err := c.doRequest(
		"GET", "/storage/"+storageID, nil,
	)
	if err != nil {
		return nil, err
	}

	var result Storage
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf(
			"decode response: %w", err,
		)
	}

	return &result, nil
}

// ListTransfers returns active transfers, optionally filtered by portal.
func (c *Client) ListTransfers(
	portalID string,
) (*TransferList, error) {
	path := "/transfers?state=active"

	if portalID != "" {
		path += "&portalId=" +
			url.QueryEscape(portalID)
	}

	data, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result TransferList
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf(
			"decode response: %w", err,
		)
	}

	return &result, nil
}

// ParsePortalArgs extracts a Portal from MCP tool arguments.
func ParsePortalArgs(args map[string]any) Portal {
	p := Portal{}

	if v, ok := args["name"].(string); ok {
		p.Name = v
	}
	if v, ok := args["url"].(string); ok {
		p.URL = v
	}
	if v, ok := args["type"].(string); ok {
		p.Type = v
	}
	if v, ok := args["notify_members"].(bool); ok {
		p.NotifyMembers = &v
	}

	auth := &Authentication{}
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
		p.Authentication = auth
	}

	linkAuth := &LinkAuth{}
	hasLinkAuth := false

	if v, ok := args["allow_unauthenticated_uploads"].(bool); ok {
		linkAuth.AllowUnauthUploads = &v
		hasLinkAuth = true
	}
	if v, ok := args["allow_unauthenticated_downloads"].(bool); ok {
		linkAuth.AllowUnauthDownloads = &v
		hasLinkAuth = true
	}

	if hasLinkAuth {
		p.LinkAuth = linkAuth
	}

	return p
}

// ParseMemberArgs extracts a PortalMember from MCP tool arguments.
func ParseMemberArgs(args map[string]any) PortalMember {
	m := PortalMember{}

	if v, ok := args["email"].(string); ok {
		m.Email = v
	}
	if v, ok := args["role"].(string); ok {
		m.Role = v
	}
	if v, ok := args["expires_on"].(string); ok {
		m.ExpiresOn = v
	}

	perms := &PortalPermissions{}
	hasPerms := false

	if v, ok := args["can_send_to_member"].(bool); ok {
		perms.CanSendToMember = &v
		hasPerms = true
	}
	if v, ok := args["can_send_to_non_member"].(bool); ok {
		perms.CanSendToNonMember = &v
		hasPerms = true
	}
	if v, ok := args["can_receive"].(bool); ok {
		perms.CanReceive = &v
		hasPerms = true
	}
	if v, ok := args["can_send_from_share"].(bool); ok {
		perms.CanSendFromShare = &v
		hasPerms = true
	}
	if v, ok := args["can_deliver_automatically"].(bool); ok {
		perms.CanAutoDeliver = &v
		hasPerms = true
	}
	if v, ok := args["can_submit"].(bool); ok {
		perms.CanSubmit = &v
		hasPerms = true
	}

	if hasPerms {
		m.Perms = perms
	}

	return m
}

// FormatJSON returns pretty-printed JSON or a fallback string.
func FormatJSON(v any) string {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", v)
	}

	return string(data)
}

// JoinNonEmpty joins non-empty strings with ", ".
func JoinNonEmpty(parts ...string) string {
	nonEmpty := make([]string, 0, len(parts))

	for _, p := range parts {
		if p != "" {
			nonEmpty = append(nonEmpty, p)
		}
	}

	return strings.Join(nonEmpty, ", ")
}
