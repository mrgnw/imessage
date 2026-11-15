package connector

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"gopkg.in/yaml.v3"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/event"

	"github.com/lrhodin/corten-matrix/pkg/rustpushgo"
)

type countingTransport struct{ requests atomic.Int32 }

func (t *countingTransport) RoundTrip(*http.Request) (*http.Response, error) {
	t.requests.Add(1)
	return nil, errors.New("test transport: no network")
}

func countOutboundPreviewRequests(t *testing.T) *countingTransport {
	transport := &countingTransport{}
	original := ogHTTPClient
	ogHTTPClient = &http.Client{Transport: transport}
	t.Cleanup(func() { ogHTTPClient = original })
	return transport
}

func TestURLPreviewsSwitchGatesLiveFetches(t *testing.T) {
	const body = "look at https://example.com/article"
	for _, enabled := range []bool{false, true} {
		transport := countOutboundPreviewRequests(t)
		connector := &IMConnector{Config: IMConfig{URLPreviews: enabled}}
		connector.Bridge = &bridgev2.Bridge{Network: connector}
		client := &IMClient{Main: connector}
		portal := &bridgev2.Portal{Portal: &database.Portal{}, Bridge: connector.Bridge}
		ctx := context.Background()

		outbound := client.convertURLPreviewToIMessage(ctx, &event.MessageEventContent{Body: body})
		inbound := convertURLPreviewToBeeper(ctx, portal, nil, &rustpushgo.WrappedMessage{}, body)

		requests := transport.requests.Load()
		if !enabled {
			if requests != 0 || outbound != body || inbound != nil {
				t.Fatalf("url_previews off: %d requests, outbound %q, inbound %v", requests, outbound, inbound)
			}
			continue
		}
		if requests == 0 || len(inbound) != 1 || !strings.HasPrefix(outbound, "\x00RL\x01") {
			t.Fatalf("url_previews on: %d requests, outbound %q, inbound %v", requests, outbound, inbound)
		}
	}
}

func TestURLPreviewsConfig(t *testing.T) {
	var c IMConfig
	if err := yaml.Unmarshal([]byte(ExampleConfig), &c); err != nil {
		t.Fatal(err)
	}
	if c.URLPreviews || c.BackfillURLPreviews() {
		t.Fatal("example config must default url_previews to off")
	}
	c.URLPreviewsInBackfill = true
	if c.BackfillURLPreviews() {
		t.Fatal("url_previews_in_backfill must not fetch while url_previews is off")
	}
	c.URLPreviews = true
	if !c.BackfillURLPreviews() {
		t.Fatal("both switches on should fetch during backfill")
	}
}
