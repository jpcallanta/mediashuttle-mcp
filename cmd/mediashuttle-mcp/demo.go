package main

import (
	"fmt"

	"mediashuttle-mcp/internal/client"
)

// RunDemo exercises the Media Shuttle API through the MCP tool
// surface, printing results to stdout.
func RunDemo(c *client.Client) {
	fmt.Println("✦ Media Shuttle MCP — Capability Demo")
	fmt.Println("─ Connecting to", client.DefaultBaseURL)
	fmt.Println()

	// ── Portal Tools ──
	fmt.Println("═══ Portal Tools ═══")
	fmt.Println()

	fmt.Println("• Tool: list_portals")
	fmt.Println("  GET /portals")
	portals, err := c.ListPortals("")
	if err != nil {
		fmt.Println("  " + err.Error())
	}
	fmt.Println()

	if err == nil && len(portals.Items) > 0 {
		pid := portals.Items[0].ID
		fmt.Println("• Using first portal:", pid)

		fmt.Println("• Tool: list_portal_users")
		fmt.Println("  GET /portals/{id}/users")
		users, err := c.ListPortalUsers(pid)
		if err != nil {
			fmt.Println("  " + err.Error())
		} else {
			fmt.Printf("  ✓ Found %d user(s)\n",
				len(users.Items))
			fmt.Println(client.FormatJSON(users))
		}
		fmt.Println()

		fmt.Println("• Tool: list_portal_storage")
		fmt.Println("  GET /portals/{id}/storage")
		ps, err := c.ListPortalStorage(pid)
		if err != nil {
			fmt.Println("  " + err.Error())
		} else {
			fmt.Printf(
				"  ✓ Found %d storage assignment(s)\n",
				len(ps.Items))
			fmt.Println(client.FormatJSON(ps))
		}
		fmt.Println()

		fmt.Println("• Tool: list_transfers (filtered by portal)")
		fmt.Println("  GET /transfers?state=active&portalId=...")
		tx, err := c.ListTransfers(pid)
		if err != nil {
			fmt.Println("  " + err.Error())
		} else {
			fmt.Printf("  ✓ Found %d transfer(s)\n",
				len(tx.Items))
			fmt.Println(client.FormatJSON(tx))
		}
		fmt.Println()

		fmt.Println("• Write tools (skipped — read-only demo):")
		fmt.Println("    create_portal, update_portal,")
		fmt.Println("    add_portal_user, update_portal_user,")
		fmt.Println("    remove_portal_user, get_portal_user,")
		fmt.Println("    assign_portal_storage")
		fmt.Println()
	}

	// ── Storage Tools ──
	fmt.Println("═══ Storage Tools ═══")
	fmt.Println()

	fmt.Println("• Tool: list_storage")
	fmt.Println("  GET /storage")
	storageList, err := c.ListStorage()
	if err != nil {
		fmt.Println("  " + err.Error())
	}
	fmt.Println()

	if err == nil && len(storageList.Items) > 0 {
		sid := storageList.Items[0].ID
		fmt.Println("• Using first storage:", sid)

		fmt.Println("• Tool: get_storage")
		fmt.Println("  GET /storage/{id}")
		storage, err := c.GetStorage(sid)
		if err != nil {
			fmt.Println("  " + err.Error())
		} else {
			fmt.Println("  ✓ OK")
			fmt.Println(client.FormatJSON(storage))
		}
		fmt.Println()
	}

	// ── Transfer Tools ──
	fmt.Println("═══ Transfer Tools ═══")
	fmt.Println()

	fmt.Println("• Tool: list_transfers (all)")
	fmt.Println("  GET /transfers?state=active")
	allTx, err := c.ListTransfers("")
	if err != nil {
		fmt.Println("  " + err.Error())
	} else {
		fmt.Printf("  ✓ Found %d transfer(s)\n",
			len(allTx.Items))
		fmt.Println(client.FormatJSON(allTx))
	}
	fmt.Println()

	fmt.Println("─ Demo complete.")
}
