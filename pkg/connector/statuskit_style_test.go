package connector

import (
	"context"
	"testing"

	"gopkg.in/yaml.v3"
	"maunium.net/go/mautrix/bridgev2/networkid"
)

func TestStatusKitNotificationStyleDefaultsToTitle(t *testing.T) {
	var c IMConfig
	if err := yaml.Unmarshal([]byte(ExampleConfig), &c); err != nil {
		t.Fatal(err)
	}
	if c.StatusKitNotificationStyle != "title" || c.StatusKitTopicStyle() {
		t.Fatalf("example style = %q, want title", c.StatusKitNotificationStyle)
	}
	for _, style := range []string{"", "notice", "Topic"} {
		c.StatusKitNotificationStyle = style
		if c.StatusKitTopicStyle() {
			t.Fatalf("style %q must fall back to title", style)
		}
	}
}

func TestStatusKitTopic(t *testing.T) {
	mode := func(s string) *string { return &s }
	cases := []struct {
		mode *string
		want string
	}{
		{nil, ""},
		{mode(""), ""},
		{mode("com.apple.donotdisturb.mode.default"), "🔕 Do Not Disturb"},
		{mode("com.apple.focus.mode.driving"), "🔕 Driving"},
		{mode("5C1E0A8E-custom-focus-uuid"), "🔕 Focus"},
	}
	for _, tc := range cases {
		if got := statusKitTopic(tc.mode); got != tc.want {
			t.Fatalf("statusKitTopic(%v) = %q, want %q", tc.mode, got, tc.want)
		}
	}
}

func TestTopicStyleKeepsMoonOutOfTitle(t *testing.T) {
	client := &IMClient{Main: &IMConnector{Config: IMConfig{StatusKitNotificationStyle: "topic"}}}
	client.statusKitPresenceByPortal.Store(networkid.PortalID("tel:+15555550123"), "com.apple.donotdisturb.mode.default")
	if client.dmTitleShowsFocus(context.Background(), "tel:+15555550123") {
		t.Fatal("topic style must not put the moon in the DM title")
	}
}
