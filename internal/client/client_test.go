package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParsePortalArgs_EmptyArgs(t *testing.T) {
	got := ParsePortalArgs(map[string]any{})

	if got != (Portal{}) {
		t.Errorf("want empty Portal, got %+v", got)
	}
}

func TestParsePortalArgs_NameAndType(t *testing.T) {
	args := map[string]any{
		"name": "Test Portal",
		"type": "share",
	}
	got := ParsePortalArgs(args)

	if got.Name != "Test Portal" {
		t.Errorf("want Name 'Test Portal', got %q", got.Name)
	}
	if got.Type != "share" {
		t.Errorf("want Type 'share', got %q", got.Type)
	}
}

func TestParsePortalArgs_AllFields(t *testing.T) {
	args := map[string]any{
		"name":                            "Full Portal",
		"url":                             "test.mediashuttle.com",
		"type":                            "send",
		"notify_members":                  true,
		"media_shuttle_auth":              true,
		"saml_auth":                       false,
		"allow_unauthenticated_uploads":   true,
		"allow_unauthenticated_downloads": false,
	}
	got := ParsePortalArgs(args)

	if got.Name != "Full Portal" {
		t.Errorf("want Name, got %q", got.Name)
	}
	if got.URL != "test.mediashuttle.com" {
		t.Errorf("want URL, got %q", got.URL)
	}
	if got.Type != "send" {
		t.Errorf("want Type, got %q", got.Type)
	}
	if got.NotifyMembers == nil || !*got.NotifyMembers {
		t.Error("want NotifyMembers true")
	}
	if got.Authentication == nil {
		t.Fatal("want Authentication non-nil")
	}
	if !got.Authentication.MediaShuttle {
		t.Error("want MediaShuttle true")
	}
	if got.LinkAuth == nil {
		t.Fatal("want LinkAuth non-nil")
	}
	if got.LinkAuth.AllowUnauthUploads == nil ||
		!*got.LinkAuth.AllowUnauthUploads {
		t.Error("want AllowUnauthUploads true")
	}
}

func TestParseMemberArgs_EmptyArgs(t *testing.T) {
	got := ParseMemberArgs(map[string]any{})

	if got != (PortalMember{}) {
		t.Errorf("want empty, got %+v", got)
	}
}

func TestParseMemberArgs_EmailOnly(t *testing.T) {
	args := map[string]any{
		"email": "user@example.com",
	}
	got := ParseMemberArgs(args)

	if got.Email != "user@example.com" {
		t.Errorf("want email, got %q", got.Email)
	}
}

func TestParseMemberArgs_AllFields(t *testing.T) {
	args := map[string]any{
		"email":                     "admin@example.com",
		"role":                      "Ops",
		"expires_on":                "2025-12-31T23:59:59.000Z",
		"can_send_to_member":        true,
		"can_send_to_non_member":    false,
		"can_receive":               true,
		"can_send_from_share":       true,
		"can_deliver_automatically": true,
		"can_submit":                false,
	}
	got := ParseMemberArgs(args)

	if got.Email != "admin@example.com" {
		t.Errorf("want email, got %q", got.Email)
	}
	if got.Role != "Ops" {
		t.Errorf("want role Ops, got %q", got.Role)
	}
	if got.Perms == nil {
		t.Fatal("want Perms non-nil")
	}
	if got.Perms.CanSendToMember == nil ||
		!*got.Perms.CanSendToMember {
		t.Error("want CanSendToMember true")
	}
}

func TestClient_ListPortals_Success(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter, r *http.Request,
		) {
			if r.Header.Get("Authorization") != "test-key" {
				t.Errorf(
					"want auth 'test-key', got %q",
					r.Header.Get("Authorization"),
				)
				http.Error(w,
					`{"message":"unauthorized"}`,
					http.StatusUnauthorized,
				)
				return
			}

			if r.URL.Path != "/portals" {
				t.Errorf(
					"want path /portals, got %s",
					r.URL.Path,
				)
				http.Error(w,
					`{"message":"not found"}`,
					http.StatusNotFound,
				)
				return
			}

			w.Header().Set(
				"Content-Type", "application/json",
			)
			json.NewEncoder(w).Encode(PortalList{
				Items: []Portal{
					{ID: "p1", Name: "P1", Type: "share"},
					{ID: "p2", Name: "P2", Type: "send"},
				},
			})
		}),
	)
	defer srv.Close()

	c := &Client{
		baseURL:    srv.URL,
		apiKey:     "test-key",
		httpClient: srv.Client(),
	}
	portals, err := c.ListPortals("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(portals.Items) != 2 {
		t.Fatalf(
			"want 2 portals, got %d",
			len(portals.Items),
		)
	}
	if portals.Items[0].Name != "P1" {
		t.Errorf(
			"want first name 'P1', got %q",
			portals.Items[0].Name,
		)
	}
}

func TestClient_ListPortals_WithFilter(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter, r *http.Request,
		) {
			urlFilter := r.URL.Query().Get("url")
			if urlFilter != "test.mediashuttle.com" {
				t.Errorf(
					"want filter, got %q",
					urlFilter,
				)
			}

			w.Header().Set(
				"Content-Type", "application/json",
			)
			json.NewEncoder(w).Encode(
				PortalList{Items: []Portal{}},
			)
		}),
	)
	defer srv.Close()

	c := &Client{
		baseURL:    srv.URL,
		apiKey:     "key",
		httpClient: srv.Client(),
	}
	_, err := c.ListPortals("test.mediashuttle.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_CreatePortal_Success(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter, r *http.Request,
		) {
			if r.Method != "POST" {
				t.Errorf(
					"want POST, got %s", r.Method,
				)
			}

			var body Portal
			if err := json.NewDecoder(
				r.Body,
			).Decode(&body); err != nil {
				t.Fatalf("decode: %v", err)
			}

			if body.Name != "New Portal" {
				t.Errorf(
					"want name 'New Portal', got %q",
					body.Name,
				)
			}

			w.Header().Set(
				"Content-Type", "application/json",
			)
			json.NewEncoder(w).Encode(Portal{
				ID:   "new-id",
				Name: body.Name,
				Type: body.Type,
			})
		}),
	)
	defer srv.Close()

	c := &Client{
		baseURL:    srv.URL,
		apiKey:     "key",
		httpClient: srv.Client(),
	}
	result, err := c.CreatePortal(
		Portal{Name: "New Portal", Type: "share"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "new-id" {
		t.Errorf(
			"want ID 'new-id', got %q", result.ID,
		)
	}
}

func TestClient_UpdatePortal_Success(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter, r *http.Request,
		) {
			if r.Method != "PATCH" {
				t.Errorf(
					"want PATCH, got %s", r.Method,
				)
			}

			want := "/portals/portal-123"
			if r.URL.Path != want {
				t.Errorf("want %s, got %s",
					want, r.URL.Path,
				)
			}

			w.Header().Set(
				"Content-Type", "application/json",
			)
			json.NewEncoder(w).Encode(Portal{
				ID:   "portal-123",
				Name: "Updated",
			})
		}),
	)
	defer srv.Close()

	c := &Client{
		baseURL:    srv.URL,
		apiKey:     "key",
		httpClient: srv.Client(),
	}
	result, err := c.UpdatePortal(
		"portal-123", UpdatePortal{Name: "Updated"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "Updated" {
		t.Errorf(
			"want name 'Updated', got %q",
			result.Name,
		)
	}
}

func TestClient_ListPortalUsers_Success(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter, r *http.Request,
		) {
			want := "/portals/p1/users"
			if r.URL.Path != want {
				t.Errorf("want %s, got %s",
					want, r.URL.Path,
				)
			}

			w.Header().Set(
				"Content-Type", "application/json",
			)
			json.NewEncoder(w).Encode(
				PortalMemberResponse{
					Items: []PortalMemberListItem{
						{
							Email:       "a@test.com",
							LastLoginOn: "2025-01-01T00:00:00Z",
						},
						{Email: "b@test.com"},
					},
				},
			)
		}),
	)
	defer srv.Close()

	c := &Client{
		baseURL:    srv.URL,
		apiKey:     "key",
		httpClient: srv.Client(),
	}
	users, err := c.ListPortalUsers("p1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users.Items) != 2 {
		t.Fatalf(
			"want 2 users, got %d",
			len(users.Items),
		)
	}
}

func TestClient_AddPortalUser_Success(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter, r *http.Request,
		) {
			if r.Method != "POST" {
				t.Errorf(
					"want POST, got %s", r.Method,
				)
			}

			want := "/portals/p1/users"
			if r.URL.Path != want {
				t.Errorf("want %s, got %s",
					want, r.URL.Path,
				)
			}

			var body PortalMember
			if err := json.NewDecoder(
				r.Body,
			).Decode(&body); err != nil {
				t.Fatalf("decode: %v", err)
			}

			if body.Email != "new@test.com" {
				t.Errorf("want email, got %q",
					body.Email,
				)
			}

			w.WriteHeader(http.StatusOK)
		}),
	)
	defer srv.Close()

	c := &Client{
		baseURL:    srv.URL,
		apiKey:     "key",
		httpClient: srv.Client(),
	}
	err := c.AddPortalUser("p1", PortalMember{
		Email: "new@test.com",
		Role:  "Member",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_RemovePortalUser_Success(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter, r *http.Request,
		) {
			if r.Method != "DELETE" {
				t.Errorf(
					"want DELETE, got %s", r.Method,
				)
			}

			want := "/portals/p1/users/user@test.com"
			if r.URL.Path != want {
				t.Errorf("want %s, got %s",
					want, r.URL.Path,
				)
			}

			w.WriteHeader(http.StatusOK)
		}),
	)
	defer srv.Close()

	c := &Client{
		baseURL:    srv.URL,
		apiKey:     "key",
		httpClient: srv.Client(),
	}
	err := c.RemovePortalUser("p1", "user@test.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_APIError(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter, r *http.Request,
		) {
			w.Header().Set(
				"Content-Type", "application/json",
			)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(APIError{
				StatusCode: 404,
				Error:      "Not Found",
				Message:    "Portal not found",
			})
		}),
	)
	defer srv.Close()

	c := &Client{
		baseURL:    srv.URL,
		apiKey:     "key",
		httpClient: srv.Client(),
	}
	_, err := c.ListPortals("")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestClient_ListTransfers_Success(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter, r *http.Request,
		) {
			if r.URL.Path != "/transfers" {
				t.Errorf(
					"want /transfers, got %s",
					r.URL.Path,
				)
			}

			state := r.URL.Query().Get("state")
			if state != "active" {
				t.Errorf(
					"want state=active, got %s",
					state,
				)
			}

			w.Header().Set(
				"Content-Type", "application/json",
			)
			json.NewEncoder(w).Encode(TransferList{
				Items: []Transfer{{
					ID:        "t1",
					Direction: "upload",
					State:     "active",
				}},
			})
		}),
	)
	defer srv.Close()

	c := &Client{
		baseURL:    srv.URL,
		apiKey:     "key",
		httpClient: srv.Client(),
	}
	transfers, err := c.ListTransfers("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(transfers.Items) != 1 {
		t.Fatalf(
			"want 1 transfer, got %d",
			len(transfers.Items),
		)
	}
	if transfers.Items[0].Direction != "upload" {
		t.Errorf(
			"want direction 'upload', got %q",
			transfers.Items[0].Direction,
		)
	}
}

func TestClient_ListStorage_Success(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter, r *http.Request,
		) {
			if r.URL.Path != "/storage" {
				t.Errorf(
					"want /storage, got %s",
					r.URL.Path,
				)
			}

			w.Header().Set(
				"Content-Type", "application/json",
			)
			json.NewEncoder(w).Encode(StorageList{
				Items: []Storage{{
					ID:     "s1",
					Type:   "s3",
					Status: "available",
				}},
			})
		}),
	)
	defer srv.Close()

	c := &Client{
		baseURL:    srv.URL,
		apiKey:     "key",
		httpClient: srv.Client(),
	}
	storage, err := c.ListStorage()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(storage.Items) != 1 {
		t.Fatalf(
			"want 1 storage, got %d",
			len(storage.Items),
		)
	}
}

func TestFormatJSON_Success(t *testing.T) {
	type testStruct struct {
		Name string `json:"name"`
		ID   int    `json:"id"`
	}

	got := FormatJSON(testStruct{
		Name: "test",
		ID:   42,
	})
	want := `{
  "name": "test",
  "id": 42
}`
	if got != want {
		t.Errorf("mismatch:\ngot:\n%s\nwant:\n%s",
			got, want,
		)
	}
}

func boolPtr(b bool) *bool { return &b }
