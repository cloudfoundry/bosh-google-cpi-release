package image

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"
	computebeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/storage/v1"
)

type completedImageOperation struct{}

func (completedImageOperation) Waiter(op *compute.Operation, _ string, _ string) (*compute.Operation, error) {
	return op, nil
}
func (completedImageOperation) WaiterB(op *computebeta.Operation, _ string, _ string) (*computebeta.Operation, error) {
	return op, nil
}

func writeResponse(t *testing.T, w io.Writer, body string) {
	if _, err := io.WriteString(w, body); err != nil {
		t.Errorf("writing response: %v", err)
	}
}

func TestImageUploadKeepsArchivePrivateAndImportsBeforeCleanup(t *testing.T) {
	var calls []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/upload/"):
			calls = append(calls, "upload")
			if r.URL.Query().Get("predefinedAcl") != "private" {
				t.Errorf("upload ACL is not explicitly private: %s", r.URL.RawQuery)
			}
			_, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
			if err != nil {
				t.Error(err)
				return
			}
			reader := multipart.NewReader(r.Body, params["boundary"])
			metadata, err := reader.NextPart()
			if err != nil {
				t.Error(err)
				return
			}
			var object storage.Object
			if err := json.NewDecoder(metadata).Decode(&object); err != nil {
				t.Error(err)
			}
			for _, acl := range object.Acl {
				if acl.Entity == "allUsers" || acl.Entity == "allAuthenticatedUsers" {
					t.Errorf("public object ACL requested: %v", acl)
				}
			}
			body, err := reader.NextPart()
			if err != nil {
				t.Error(err)
				return
			}
			image, err := io.ReadAll(body)
			if err != nil {
				t.Error(err)
				return
			}
			if !bytes.Equal(image, []byte("fixture archive")) {
				t.Error("archive data changed")
			}
			writeResponse(t, w, `{"name":"private-image","mediaLink":"https://storage.googleapis.com/private-bucket/private-image"}`)
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/b"):
			calls = append(calls, "bucket")
			var bucket storage.Bucket
			if err := json.NewDecoder(r.Body).Decode(&bucket); err != nil {
				t.Error(err)
			}
			if bucket.IamConfiguration == nil || bucket.IamConfiguration.PublicAccessPrevention != "enforced" {
				t.Error("temporary bucket permits public sharing")
			}
			if err := json.NewEncoder(w).Encode(bucket); err != nil {
				t.Error(err)
			}
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/images"):
			calls = append(calls, "import")
			var image compute.Image
			if err := json.NewDecoder(r.Body).Decode(&image); err != nil {
				t.Error(err)
			}
			if image.RawDisk == nil || image.RawDisk.Source != "https://storage.googleapis.com/private-bucket/private-image" {
				t.Error("image import lost private archive URL")
			}
			writeResponse(t, w, `{"name":"operation","status":"DONE"}`)
		case r.Method == http.MethodDelete:
			calls = append(calls, "cleanup")
			writeResponse(t, w, `{}`)
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			http.Error(w, "unexpected", 400)
		}
	}))
	defer server.Close()
	computeService, err := compute.NewService(context.Background(), option.WithHTTPClient(server.Client()))
	if err != nil {
		t.Fatal(err)
	}
	computeService.BasePath = server.URL + "/compute/v1/"
	storageService, err := storage.NewService(context.Background(), option.WithHTTPClient(server.Client()))
	if err != nil {
		t.Fatal(err)
	}
	storageService.BasePath = server.URL + "/storage/v1/"
	file := filepath.Join(t.TempDir(), "archive.tgz")
	if err := os.WriteFile(file, []byte("fixture archive"), 0600); err != nil {
		t.Fatal(err)
	}
	service := GoogleImageService{project: "project", computeService: computeService, storageService: storageService, operationService: completedImageOperation{}, uuidGen: boshuuid.NewGenerator(), logger: boshlog.NewWriterLogger(boshlog.LevelDebug, io.Discard)}
	if _, err := service.CreateFromTarball(file, "fixture"); err != nil {
		t.Fatal(err)
	}
	if strings.Join(calls, ",") != "bucket,upload,import,cleanup,cleanup" {
		t.Fatalf("wrong operation order: %v", calls)
	}
}
