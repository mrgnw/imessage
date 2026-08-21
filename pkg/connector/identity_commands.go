// mautrix-imessage - A Matrix-iMessage puppeting bridge.
// Copyright (C) 2024 Ludvig Rhodin
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package connector

import (
	"fmt"
	"strconv"
	"strings"

	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/commands"
)

// ============================================================================
// list-identities / set-identity commands
// ============================================================================

// loginLabel returns the display label for a login: the user's preferred
// handle if one is set, otherwise the raw login ID.
func loginLabel(l *bridgev2.UserLogin) string {
	label := string(l.ID)
	if meta, ok := l.Metadata.(*UserLoginMetadata); ok && meta.PreferredHandle != "" {
		label = meta.PreferredHandle
	}
	return label
}

// parseLoginSelection converts a reply back to the matching login. Mirrors
// parseHandleSelection (login.go) but keyed by loginLabel — accepts a bare
// 1-based index ("2"), the full numbered label we emitted ("2. tel:+1..."),
// or a raw label/ID match. Unlike parseHandleSelection, this returns nil
// (rather than falling back to the trimmed input) when nothing matches,
// since callers here always need an actual login object to act on.
func parseLoginSelection(selected string, logins []*bridgev2.UserLogin) *bridgev2.UserLogin {
	trimmed := strings.TrimSpace(selected)
	if n, err := strconv.Atoi(trimmed); err == nil && n >= 1 && n <= len(logins) {
		return logins[n-1]
	}
	for i, l := range logins {
		label := loginLabel(l)
		numbered := fmt.Sprintf("%d. %s", i+1, label)
		if trimmed == label || trimmed == numbered || trimmed == string(l.ID) {
			return l
		}
	}
	return nil
}

// resolvePortalLogin returns the login that owns ce.Portal's room (matched by
// networkid.UserLoginID against the portal's Receiver), or nil if ce.Portal
// is nil or no cached login matches.
func resolvePortalLogin(ce *commands.Event) *bridgev2.UserLogin {
	if ce.Portal == nil {
		return nil
	}
	for _, l := range ce.User.GetCachedUserLogins() {
		if l != nil && l.ID == ce.Portal.Receiver {
			return l
		}
	}
	return nil
}

// formatHandleList renders a numbered list of handles, one per line, marking
// whichever entry equals current with " (current)".
func formatHandleList(handles []string, current string) string {
	var sb strings.Builder
	for i, h := range handles {
		if h == current {
			sb.WriteString(fmt.Sprintf("%d. %s (current)\n", i+1, h))
		} else {
			sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, h))
		}
	}
	return sb.String()
}

// formatHandleBullets renders a bulleted (unnumbered) list of handles, used
// for the read-only grouped list-identities display where there's no
// per-handle input target to number.
func formatHandleBullets(handles []string, current string) string {
	var sb strings.Builder
	for _, h := range handles {
		if h == current {
			sb.WriteString(fmt.Sprintf("  • %s (current)\n", h))
		} else {
			sb.WriteString(fmt.Sprintf("  • %s\n", h))
		}
	}
	return sb.String()
}

// renderLoginNumberedList renders a numbered list of logins (by loginLabel),
// one per line, for account-selection prompts and errors.
func renderLoginNumberedList(logins []*bridgev2.UserLogin) string {
	var sb strings.Builder
	for i, l := range logins {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, loginLabel(l)))
	}
	return sb.String()
}

// invalidAccountMessage builds the error reply for an unrecognized account
// selector, listing the valid accounts.
func invalidAccountMessage(logins []*bridgev2.UserLogin) string {
	var sb strings.Builder
	sb.WriteString("Unknown account. Valid accounts:\n\n")
	sb.WriteString(renderLoginNumberedList(logins))
	return sb.String()
}

var cmdListIdentities = &commands.FullHandler{
	Name: "list-identities",
	Func: fnListIdentities,
	Help: commands.HelpMeta{
		Section:     commands.HelpSectionAuth,
		Description: "List the iMessage handles registered to your Apple ID account(s) and show which one is currently used for outgoing messages.",
	},
	RequiresLogin: true,
}

func fnListIdentities(ce *commands.Event) {
	if ce.Portal != nil {
		login := resolvePortalLogin(ce)
		if login == nil {
			ce.Reply("Could not resolve the iMessage account for this room.")
			return
		}
		replyIdentityList(ce, login)
		return
	}

	logins := ce.User.GetCachedUserLogins()
	switch len(logins) {
	case 0:
		ce.Reply("You're not signed in to iMessage. Run `$cmdprefix login` first.")
	case 1:
		replyIdentityList(ce, logins[0])
	default:
		var sb strings.Builder
		sb.WriteString("**Logged-in iMessage accounts:**\n\n")
		for _, l := range logins {
			sb.WriteString(fmt.Sprintf("**%s**\n", loginLabel(l)))
			client, ok := l.Client.(*IMClient)
			if !ok || client == nil || client.client == nil {
				sb.WriteString("  _(not connected)_\n\n")
				continue
			}
			sb.WriteString(formatHandleBullets(client.getAllHandles(), client.getHandle()))
			sb.WriteString("\n")
		}
		sb.WriteString("Use `$cmdprefix set-identity` to change an account's outgoing handle.")
		ce.Reply(sb.String())
	}
}

// replyIdentityList sends the flat (single-login) registered-identities list
// used both for the management room with exactly one login and for a portal
// room (which always targets exactly one login).
func replyIdentityList(ce *commands.Event, login *bridgev2.UserLogin) {
	client, ok := login.Client.(*IMClient)
	if !ok || client == nil || client.client == nil {
		ce.Reply("Bridge client not available.")
		return
	}
	handles := client.getAllHandles()
	if len(handles) == 0 {
		ce.Reply("No registered handles found for this account.")
		return
	}
	var sb strings.Builder
	sb.WriteString("**Registered identities:**\n\n")
	sb.WriteString(formatHandleList(handles, client.getHandle()))
	sb.WriteString("\nUse `$cmdprefix set-identity` to change which one outgoing messages send from.")
	ce.Reply(sb.String())
}

var cmdSetIdentity = &commands.FullHandler{
	Name:    "set-identity",
	Aliases: []string{"identity"},
	Func:    fnSetIdentity,
	Help: commands.HelpMeta{
		Section:     commands.HelpSectionAuth,
		Description: "Change which registered handle outgoing iMessages are sent from. In a portal room, applies to that room's account. In the management room with multiple Apple ID logins, pass the account as the first argument.",
		Args:        "[account] [handle]",
	},
	RequiresLogin: true,
}

func fnSetIdentity(ce *commands.Event) {
	args := strings.Fields(ce.RawArgs)

	if ce.Portal != nil {
		login := resolvePortalLogin(ce)
		if login == nil {
			ce.Reply("Could not resolve the iMessage account for this room.")
			return
		}
		if len(args) > 1 {
			ce.Reply("Usage in a portal room: `$cmdprefix set-identity [handle]`")
			return
		}
		setIdentityForLogin(ce, login, strings.TrimSpace(ce.RawArgs))
		return
	}

	logins := ce.User.GetCachedUserLogins()
	switch len(logins) {
	case 0:
		ce.Reply("You're not signed in to iMessage. Run `$cmdprefix login` first.")
	case 1:
		if len(args) > 1 {
			ce.Reply("Usage: `$cmdprefix set-identity [handle]`")
			return
		}
		setIdentityForLogin(ce, logins[0], strings.TrimSpace(ce.RawArgs))
	default:
		switch len(args) {
		case 0:
			var sb strings.Builder
			sb.WriteString("**Choose an account:**\n\n")
			sb.WriteString(renderLoginNumberedList(logins))
			sb.WriteString("\nReply with a number to choose an account, or `$cmdprefix cancel` to cancel.")
			ce.Reply(sb.String())

			commands.StoreCommandState(ce.User, &commands.CommandState{
				Action: "set-identity-login",
				Next: commands.MinimalCommandHandlerFunc(func(ce *commands.Event) {
					commands.StoreCommandState(ce.User, nil)
					login := parseLoginSelection(strings.TrimSpace(ce.RawArgs), logins)
					if login == nil {
						ce.Reply("Please reply with a number between 1 and %d, or run `$cmdprefix set-identity` again.", len(logins))
						return
					}
					startHandlePicker(ce, login)
				}),
				Cancel: func() {},
			})
		case 1:
			login := parseLoginSelection(args[0], logins)
			if login == nil {
				ce.Reply(invalidAccountMessage(logins))
				return
			}
			startHandlePicker(ce, login)
		case 2:
			login := parseLoginSelection(args[0], logins)
			if login == nil {
				ce.Reply(invalidAccountMessage(logins))
				return
			}
			setIdentityForLogin(ce, login, args[1])
		default:
			ce.Reply("Usage: `$cmdprefix set-identity [account] [handle]`")
		}
	}
}

// setIdentityForLogin implements "Step B": given a resolved login and a
// handle selector (which may be empty), either starts the interactive handle
// picker (empty selector) or applies the selection directly (one-shot path).
func setIdentityForLogin(ce *commands.Event, login *bridgev2.UserLogin, handleArg string) {
	client, ok := login.Client.(*IMClient)
	if !ok || client == nil || client.client == nil {
		ce.Reply("Bridge client not available for that account.")
		return
	}
	handles := client.getAllHandles()
	if len(handles) == 0 {
		ce.Reply("No registered handles found for that account.")
		return
	}
	if handleArg == "" {
		startHandlePicker(ce, login)
		return
	}
	applyIdentitySelection(ce, login, client, handleArg, handles)
}

// startHandlePicker prompts for which of a login's registered handles to
// send from, then stores command state to apply the reply.
func startHandlePicker(ce *commands.Event, login *bridgev2.UserLogin) {
	client, ok := login.Client.(*IMClient)
	if !ok || client == nil || client.client == nil {
		ce.Reply("Bridge client not available for that account.")
		return
	}
	handles := client.getAllHandles()
	if len(handles) == 0 {
		ce.Reply("No registered handles found for that account.")
		return
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("**Choose a handle for %s:**\n\n", loginLabel(login)))
	sb.WriteString(formatHandleList(handles, client.getHandle()))
	sb.WriteString("\nReply with a number, or `$cmdprefix cancel` to cancel.")
	ce.Reply(sb.String())

	commands.StoreCommandState(ce.User, &commands.CommandState{
		Action: "set-identity-handle",
		Next: commands.MinimalCommandHandlerFunc(func(ce *commands.Event) {
			commands.StoreCommandState(ce.User, nil)
			applyIdentitySelection(ce, login, client, strings.TrimSpace(ce.RawArgs), handles)
		}),
		Cancel: func() {},
	})
}

// applyIdentitySelection resolves selector against handles (a bare index, a
// numbered label, or a raw handle) and applies it as the login's send-from
// identity. If the resolved handle isn't currently registered, it forces a
// re-registration with Apple and retries once against the refreshed handle
// set before giving up.
func applyIdentitySelection(ce *commands.Event, login *bridgev2.UserLogin, client *IMClient, selector string, handles []string) {
	// Re-check even though callers already checked at prompt time: this can run
	// well after that check, from a stored CommandState's Next handler, and the
	// login may have disconnected (client.client set back to nil) in the interim.
	if client == nil || client.client == nil {
		ce.Reply("Bridge client is no longer available for that account. Please try again.")
		return
	}
	resolved := parseHandleSelection(selector, handles)

	registered := false
	for _, h := range handles {
		if h == resolved {
			registered = true
			break
		}
	}

	if !registered {
		_, err := client.client.ForceReregisterIdentity()
		if err != nil {
			ce.Reply(
				"%q is not a currently registered handle, and refreshing the account's registration failed: %v. Registered handles: %s",
				resolved, err, strings.Join(handles, ", "),
			)
			return
		}
		refreshed := client.client.GetHandles()
		client.setIdentity(client.getHandle(), refreshed)
		handles = refreshed

		registered = false
		for _, h := range handles {
			if h == resolved {
				registered = true
				break
			}
		}
		if !registered {
			ce.Reply(
				"%q is not a registered handle for this account. Registered handles: %s",
				resolved, strings.Join(handles, ", "),
			)
			return
		}
	}

	client.setIdentity(resolved, handles)
	if meta, ok := login.Metadata.(*UserLoginMetadata); ok {
		meta.PreferredHandle = resolved
		_ = login.Save(ce.Ctx)
	}
	ce.Reply("Now sending as **%s**. This is what recipients will see your messages \"from\".", resolved)
}
