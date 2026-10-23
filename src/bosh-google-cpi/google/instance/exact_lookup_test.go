package instance

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	computebeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

func TestExactLookupAndAmbiguousZones(t *testing.T) {
	cases := []struct {
		name      string
		names     []string
		paginated bool
		found     bool
		ambiguous bool
	}{
		{name: "absent"},
		{name: "suffix collision", names: []string{"other-target"}},
		{name: "exact with unrelated", names: []string{"other-target", "target"}, found: true},
		{name: "duplicate zones", names: []string{"target", "target"}, ambiguous: true},
		{name: "duplicate pages", names: []string{"target", "target"}, paginated: true, ambiguous: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if got := r.URL.Query().Get("filter"); got != `name = "target"` {
					t.Errorf("lookup filter=%q", got)
				}
				items := map[string]interface{}{}
				pageToken := ""
				for n, name := range tc.names {
					if tc.paginated {
						if r.URL.Query().Get("pageToken") == "" && n > 0 {
							pageToken = "second"
							continue
						}
						if r.URL.Query().Get("pageToken") != "" && n == 0 {
							continue
						}
					}
					items[fmt.Sprintf("zones/zone%d", n)] = map[string]interface{}{"instances": []map[string]interface{}{{"name": name, "zone": fmt.Sprintf("zones/zone%d", n)}}}
				}
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(map[string]interface{}{"items": items, "nextPageToken": pageToken}); err != nil {
					t.Errorf("encoding response: %v", err)
				}
			}))
			defer server.Close()
			service, err := compute.NewService(context.Background(), option.WithHTTPClient(server.Client()))
			if err != nil {
				t.Fatal(err)
			}
			service.BasePath = server.URL + "/"
			beta, err := computebeta.NewService(context.Background(), option.WithHTTPClient(server.Client()))
			if err != nil {
				t.Fatal(err)
			}
			beta.BasePath = server.URL + "/"
			logger := boshlog.NewWriterLogger(boshlog.LevelDebug, io.Discard)
			i := GoogleInstanceService{computeService: service, computeServiceB: beta, project: "project", logger: logger}
			_, found, err := i.Find("target", "")
			_, betaFound, betaErr := i.FindBeta("target", "")
			if betaFound != tc.found || (betaErr != nil) != tc.ambiguous {
				t.Fatalf("beta found=%v err=%v", betaFound, betaErr)
			}
			if (err != nil) != tc.ambiguous {
				t.Fatalf("ambiguous=%v, got err=%v", tc.ambiguous, err)
			}
			if found != tc.found {
				t.Fatalf("found=%v expected %v", found, tc.found)
			}
		})
	}
}
