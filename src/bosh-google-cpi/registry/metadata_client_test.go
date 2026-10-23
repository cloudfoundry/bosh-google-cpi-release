package registry

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

type metadataTestClient struct{ service *compute.Service }

func (c metadataTestClient) Project() string                  { return "project" }
func (c metadataTestClient) ComputeService() *compute.Service { return c.service }

func writeResponse(t *testing.T, w io.Writer, body string) {
	if _, err := io.WriteString(w, body); err != nil {
		t.Errorf("writing response: %v", err)
	}
}

func newTestComputeService(t *testing.T, server *httptest.Server) *compute.Service {
	service, err := compute.NewService(context.Background(), option.WithHTTPClient(server.Client()))
	if err != nil {
		t.Fatal(err)
	}
	service.BasePath = server.URL + "/"
	return service
}

func TestMetadataRejectsCrossZoneAmbiguityBeforeWriting(t *testing.T) {
	for _, paginated := range []bool{false, true} {
		t.Run(map[bool]string{false: "same page", true: "different pages"}[paginated], func(t *testing.T) {
			writes := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					writes++
					http.Error(w, "unexpected write", 400)
					return
				}
				if got := r.URL.Query().Get("filter"); got != `name = "target"` {
					t.Errorf("filter=%q", got)
				}
				entries := map[string]interface{}{}
				next := ""
				for _, zone := range []string{"a", "b"} {
					if paginated {
						if zone == "b" && r.URL.Query().Get("pageToken") == "" {
							next = "second"
							continue
						}
						if zone == "a" && r.URL.Query().Get("pageToken") != "" {
							continue
						}
					}
					entries[zone] = map[string]interface{}{"instances": []map[string]interface{}{{"name": "target", "zone": "zones/" + zone, "metadata": map[string]interface{}{"fingerprint": "fixture"}}}}
				}
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(map[string]interface{}{"items": entries, "nextPageToken": next}); err != nil {
					t.Errorf("encoding response: %v", err)
				}
			}))
			defer server.Close()
			service := newTestComputeService(t, server)
			c := MetadataClient{googleClient: metadataTestClient{service}, logger: boshlog.NewWriterLogger(boshlog.LevelDebug, io.Discard), options: ClientOptions{GCEMetadataKey: GCEMetadataKey}}
			err := c.Update("target", AgentSettings{Mbus: "nats://user:fixture-secret@host"})
			if err == nil || !strings.Contains(err.Error(), "ambiguous") {
				t.Fatalf("expected ambiguous name rejection, got %v", err)
			}
			if writes != 0 {
				t.Fatalf("wrote metadata to ambiguous resource %d times", writes)
			}
		})
	}
}

func TestMetadataErrorsAndUpdateLogsDoNotContainValues(t *testing.T) {
	for _, badSettings := range []string{"missing", "malformed"} {
		t.Run(badSettings, func(t *testing.T) {
			writes := 0
			var sent compute.Metadata
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if r.Method == http.MethodPost {
					writes++
					if err := json.NewDecoder(r.Body).Decode(&sent); err != nil {
						t.Errorf("decoding request: %v", err)
					}
					writeResponse(t, w, `{"name":"op","zone":"zones/a","status":"DONE"}`)
					return
				}
				if strings.Contains(r.URL.Path, "/operations/") {
					writeResponse(t, w, `{"name":"op","zone":"zones/a","status":"DONE"}`)
					return
				}
				items := []map[string]string{{"key": "other", "value": "fixture-other-secret"}}
				if badSettings == "malformed" {
					items = append(items, map[string]string{"key": GCEMetadataKey, "value": `{"mbus":"nats://user:fixture-secret@host",`})
				}
				if err := json.NewEncoder(w).Encode(map[string]interface{}{"items": map[string]interface{}{"a": map[string]interface{}{"instances": []map[string]interface{}{{"name": "target", "zone": "zones/a", "metadata": map[string]interface{}{"fingerprint": "fixture", "items": items}}}}}}); err != nil {
					t.Errorf("encoding response: %v", err)
				}
			}))
			defer server.Close()
			service := newTestComputeService(t, server)
			var logs bytes.Buffer
			c := MetadataClient{googleClient: metadataTestClient{service}, logger: boshlog.NewWriterLogger(boshlog.LevelDebug, &logs), options: ClientOptions{GCEMetadataKey: GCEMetadataKey}}
			_, err := c.Fetch("target")
			if err == nil {
				t.Fatal("expected missing/malformed metadata error")
			}
			if strings.Contains(err.Error(), "fixture-secret") || strings.Contains(err.Error(), "fixture-other-secret") {
				t.Fatalf("metadata disclosed in error: %v", err)
			}
			err = c.Update("target", AgentSettings{Mbus: "nats://user:fixture-secret@host"})
			if err != nil {
				t.Fatal(err)
			}
			if writes != 1 {
				t.Fatalf("expected one update, got %d", writes)
			}
			if strings.Contains(logs.String(), "fixture-secret") || strings.Contains(logs.String(), "fixture-other-secret") {
				t.Fatalf("metadata disclosed in debug log: %s", logs.String())
			}
			found := false
			for _, item := range sent.Items {
				if item == nil {
					t.Fatal("nil metadata item")
				}
				if item.Key == GCEMetadataKey && strings.Contains(*item.Value, "fixture-secret") {
					found = true
				}
			}
			if !found {
				t.Fatal("actual metadata update lost credentials")
			}
		})
	}
}
