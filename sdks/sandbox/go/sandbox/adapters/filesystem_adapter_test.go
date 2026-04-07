package adapters

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/config"
	execdapi "github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/internal/openapi/execd"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
)

func TestFilesystemAdapterWriteAndReadFile(t *testing.T) {
	var requests []httpRequestCapture
	httpClient := &http.Client{
		Transport: roundTripCaptureFunc(func(req *http.Request) (*http.Response, error) {
			body := []byte{}
			if req.Body != nil {
				var err error
				body, err = io.ReadAll(req.Body)
				if err != nil {
					return nil, err
				}
			}
			requests = append(requests, httpRequestCapture{
				Method: req.Method,
				Path:   req.URL.Path,
				Query:  req.URL.RawQuery,
				Body:   string(body),
				APIKey: req.Header.Get("OPEN-SANDBOX-API-KEY"),
				Header: req.Header.Get("X-Test-Header"),
				Range:  req.Header.Get("Range"),
				Type:   req.Header.Get("Content-Type"),
			})

			switch req.URL.Path {
			case "/files/upload":
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(`{}`)),
				}, nil
			case "/files/download":
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"text/plain; charset=utf-8"}},
					Body:       io.NopCloser(strings.NewReader("hello from sandbox")),
				}, nil
			default:
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(`{"code":"not_found","message":"not found"}`)),
				}, nil
			}
		}),
	}

	adapter := NewFilesystemAdapter(&fakeFilesystemClient{}, "http://sandbox.example:44772", &config.ConnectionConfig{
		APIKey:     "api-key-1",
		Headers:    map[string]string{"X-Test-Header": "value-1"},
		HTTPClient: httpClient,
	})

	err := adapter.WriteFiles(context.Background(), []models.WriteEntry{
		{Path: "/tmp/hello.txt", Data: "hello from sandbox", Mode: 0644, Encoding: "utf-8"},
	})
	if err != nil {
		t.Fatalf("write files: %v", err)
	}

	content, err := adapter.ReadFile(context.Background(), "/tmp/hello.txt", &models.ReadFileOptions{
		Encoding: "utf-8",
		Range:    "bytes=0-4",
	})
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if content != "hello from sandbox" {
		t.Fatalf("unexpected file content: %q", content)
	}

	if len(requests) != 2 {
		t.Fatalf("expected 2 direct http requests, got %d", len(requests))
	}
	if requests[0].Path != "/files/upload" || requests[0].Method != http.MethodPost {
		t.Fatalf("unexpected upload request: %+v", requests[0])
	}
	if !strings.Contains(requests[0].Type, "multipart/form-data") {
		t.Fatalf("unexpected upload content type: %q", requests[0].Type)
	}
	if !strings.Contains(requests[0].Body, `"path":"/tmp/hello.txt"`) {
		t.Fatalf("expected upload metadata body, got %q", requests[0].Body)
	}
	verifyMultipartUploadShape(t, requests[0].Type, requests[0].Body)
	if requests[1].Path != "/files/download" || requests[1].Method != http.MethodGet {
		t.Fatalf("unexpected download request: %+v", requests[1])
	}
	if requests[1].Range != "bytes=0-4" {
		t.Fatalf("unexpected range header: %q", requests[1].Range)
	}
	if requests[1].APIKey != "api-key-1" || requests[1].Header != "value-1" {
		t.Fatalf("expected propagated headers, got %+v", requests[1])
	}
}

func verifyMultipartUploadShape(t *testing.T, contentType, body string) {
	t.Helper()

	parts := strings.Split(contentType, "boundary=")
	if len(parts) != 2 {
		t.Fatalf("missing multipart boundary in content type: %q", contentType)
	}

	reader := multipart.NewReader(strings.NewReader(body), strings.TrimSpace(parts[1]))

	metadataPart, err := reader.NextPart()
	if err != nil {
		t.Fatalf("read metadata part: %v", err)
	}
	if got := metadataPart.FormName(); got != "metadata" {
		t.Fatalf("unexpected metadata form name: %q", got)
	}
	if got := metadataPart.FileName(); got != "metadata" {
		t.Fatalf("unexpected metadata filename: %q", got)
	}
	if got := metadataPart.Header.Get("Content-Type"); got != "application/json" {
		t.Fatalf("unexpected metadata content type: %q", got)
	}

	filePart, err := reader.NextPart()
	if err != nil {
		t.Fatalf("read file part: %v", err)
	}
	if got := filePart.FormName(); got != "file" {
		t.Fatalf("unexpected file form name: %q", got)
	}
	if got := filePart.FileName(); got != "hello.txt" {
		t.Fatalf("unexpected file filename: %q", got)
	}
	if got := filePart.Header.Get("Content-Type"); got != "text/plain; charset=utf-8" {
		t.Fatalf("unexpected file content type: %q", got)
	}

	if _, err := reader.NextPart(); err != io.EOF {
		t.Fatalf("expected exactly two multipart parts, got err=%v", err)
	}
}

func TestFilesystemAdapterSearchDeleteAndInfoConversion(t *testing.T) {
	now := time.Date(2026, 3, 25, 10, 0, 0, 0, time.UTC)
	client := &fakeFilesystemClient{
		searchResponse: &execdapi.SearchFilesResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &[]execdapi.FileInfo{
				{
					Path:       "/tmp/app.log",
					Size:       42,
					Mode:       0644,
					Owner:      "root",
					Group:      "root",
					CreatedAt:  now,
					ModifiedAt: now,
				},
			},
		},
		infoResponse: &execdapi.GetFilesInfoResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &map[string]execdapi.FileInfo{
				"/tmp/app.log": {
					Path:       "/tmp/app.log",
					Size:       42,
					Mode:       0644,
					Owner:      "root",
					Group:      "root",
					CreatedAt:  now,
					ModifiedAt: now,
				},
			},
		},
		removeFilesResponse: &execdapi.RemoveFilesResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
		},
	}

	adapter := NewFilesystemAdapter(client, "http://sandbox.example:44772", &config.ConnectionConfig{})

	results, err := adapter.Search(context.Background(), models.SearchEntry{
		Path:    "/tmp",
		Pattern: "*.log",
	})
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(results) != 1 || results[0].Path != "/tmp/app.log" {
		t.Fatalf("unexpected search results: %+v", results)
	}

	info, err := adapter.GetFileInfo(context.Background(), []string{"/tmp/app.log"})
	if err != nil {
		t.Fatalf("get file info: %v", err)
	}
	if info["/tmp/app.log"].Owner != "root" {
		t.Fatalf("unexpected file owner: %+v", info["/tmp/app.log"])
	}

	if err := adapter.DeleteFiles(context.Background(), []string{"/tmp/app.log"}); err != nil {
		t.Fatalf("delete files: %v", err)
	}
	if len(client.removeFilesPaths) != 1 || client.removeFilesPaths[0] != "/tmp/app.log" {
		t.Fatalf("unexpected delete paths: %+v", client.removeFilesPaths)
	}
}

func TestFilesystemAdapterMapsReplaceMovePermissionsAndDirectories(t *testing.T) {
	client := &fakeFilesystemClient{
		replaceContentResponse: &execdapi.ReplaceContentResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
		},
		renameFilesResponse: &execdapi.RenameFilesResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
		},
		chmodFilesResponse: &execdapi.ChmodFilesResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
		},
		makeDirsResponse: &execdapi.MakeDirsResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
		},
		removeDirsResponse: &execdapi.RemoveDirsResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
		},
	}

	adapter := NewFilesystemAdapter(client, "http://sandbox.example:44772", &config.ConnectionConfig{})

	if err := adapter.ReplaceContents(context.Background(), []models.ContentReplaceEntry{
		{Path: "/tmp/app.log", OldContent: "old", NewContent: "new"},
	}); err != nil {
		t.Fatalf("replace contents: %v", err)
	}
	if got := client.replaceBody["/tmp/app.log"].Old; got != "old" {
		t.Fatalf("unexpected replace old content: %q", got)
	}

	if err := adapter.MoveFiles(context.Background(), []models.MoveEntry{
		{Src: "/tmp/a.txt", Dest: "/tmp/b.txt"},
	}); err != nil {
		t.Fatalf("move files: %v", err)
	}
	if len(client.renameBody) != 1 || client.renameBody[0].Dest != "/tmp/b.txt" {
		t.Fatalf("unexpected move body: %+v", client.renameBody)
	}

	if err := adapter.SetPermissions(context.Background(), []models.SetPermissionEntry{
		{Path: "/tmp/b.txt", Mode: 0600, Owner: "sandbox", Group: "sandbox"},
	}); err != nil {
		t.Fatalf("set permissions: %v", err)
	}
	if got := client.chmodBody["/tmp/b.txt"].Mode; got != 0600 {
		t.Fatalf("unexpected chmod mode: %d", got)
	}

	if err := adapter.CreateDirectories(context.Background(), []models.WriteEntry{
		{Path: "/tmp/data", Mode: 0755, Owner: "sandbox", Group: "sandbox"},
	}); err != nil {
		t.Fatalf("create directories: %v", err)
	}
	if got := client.makeDirsBody["/tmp/data"].Mode; got != 0755 {
		t.Fatalf("unexpected mkdir mode: %d", got)
	}

	if err := adapter.DeleteDirectories(context.Background(), []string{"/tmp/data"}); err != nil {
		t.Fatalf("delete directories: %v", err)
	}
	if len(client.removeDirsPaths) != 1 || client.removeDirsPaths[0] != "/tmp/data" {
		t.Fatalf("unexpected remove dirs paths: %+v", client.removeDirsPaths)
	}
}

type fakeFilesystemClient struct {
	searchResponse         *execdapi.SearchFilesResponse
	infoResponse           *execdapi.GetFilesInfoResponse
	removeFilesResponse    *execdapi.RemoveFilesResponse
	replaceContentResponse *execdapi.ReplaceContentResponse
	renameFilesResponse    *execdapi.RenameFilesResponse
	chmodFilesResponse     *execdapi.ChmodFilesResponse
	makeDirsResponse       *execdapi.MakeDirsResponse
	removeDirsResponse     *execdapi.RemoveDirsResponse

	removeFilesPaths []string
	replaceBody      execdapi.ReplaceContentJSONRequestBody
	renameBody       execdapi.RenameFilesJSONRequestBody
	chmodBody        execdapi.ChmodFilesJSONRequestBody
	makeDirsBody     execdapi.MakeDirsJSONRequestBody
	removeDirsPaths  []string
}

func (f *fakeFilesystemClient) SearchFilesWithResponse(ctx context.Context, params *execdapi.SearchFilesParams, reqEditors ...execdapi.RequestEditorFn) (*execdapi.SearchFilesResponse, error) {
	return f.searchResponse, nil
}

func (f *fakeFilesystemClient) GetFilesInfoWithResponse(ctx context.Context, params *execdapi.GetFilesInfoParams, reqEditors ...execdapi.RequestEditorFn) (*execdapi.GetFilesInfoResponse, error) {
	return f.infoResponse, nil
}

func (f *fakeFilesystemClient) RemoveFilesWithResponse(ctx context.Context, params *execdapi.RemoveFilesParams, reqEditors ...execdapi.RequestEditorFn) (*execdapi.RemoveFilesResponse, error) {
	f.removeFilesPaths = append([]string(nil), params.Path...)
	return f.removeFilesResponse, nil
}

func (f *fakeFilesystemClient) ReplaceContentWithResponse(ctx context.Context, body execdapi.ReplaceContentJSONRequestBody, reqEditors ...execdapi.RequestEditorFn) (*execdapi.ReplaceContentResponse, error) {
	f.replaceBody = body
	return f.replaceContentResponse, nil
}

func (f *fakeFilesystemClient) RenameFilesWithResponse(ctx context.Context, body execdapi.RenameFilesJSONRequestBody, reqEditors ...execdapi.RequestEditorFn) (*execdapi.RenameFilesResponse, error) {
	f.renameBody = body
	return f.renameFilesResponse, nil
}

func (f *fakeFilesystemClient) ChmodFilesWithResponse(ctx context.Context, body execdapi.ChmodFilesJSONRequestBody, reqEditors ...execdapi.RequestEditorFn) (*execdapi.ChmodFilesResponse, error) {
	f.chmodBody = body
	return f.chmodFilesResponse, nil
}

func (f *fakeFilesystemClient) MakeDirsWithResponse(ctx context.Context, body execdapi.MakeDirsJSONRequestBody, reqEditors ...execdapi.RequestEditorFn) (*execdapi.MakeDirsResponse, error) {
	f.makeDirsBody = body
	return f.makeDirsResponse, nil
}

func (f *fakeFilesystemClient) RemoveDirsWithResponse(ctx context.Context, params *execdapi.RemoveDirsParams, reqEditors ...execdapi.RequestEditorFn) (*execdapi.RemoveDirsResponse, error) {
	f.removeDirsPaths = append([]string(nil), params.Path...)
	return f.removeDirsResponse, nil
}

type httpRequestCapture struct {
	Method string
	Path   string
	Query  string
	Body   string
	APIKey string
	Header string
	Range  string
	Type   string
}

type roundTripCaptureFunc func(req *http.Request) (*http.Response, error)

func (f roundTripCaptureFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
