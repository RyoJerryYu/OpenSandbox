package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/config"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/internal/convert"
	execdapi "github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/internal/openapi/execd"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
)

type FilesystemClient interface {
	SearchFilesWithResponse(ctx context.Context, params *execdapi.SearchFilesParams, reqEditors ...execdapi.RequestEditorFn) (*execdapi.SearchFilesResponse, error)
	GetFilesInfoWithResponse(ctx context.Context, params *execdapi.GetFilesInfoParams, reqEditors ...execdapi.RequestEditorFn) (*execdapi.GetFilesInfoResponse, error)
	RemoveFilesWithResponse(ctx context.Context, params *execdapi.RemoveFilesParams, reqEditors ...execdapi.RequestEditorFn) (*execdapi.RemoveFilesResponse, error)
	ReplaceContentWithResponse(ctx context.Context, body execdapi.ReplaceContentJSONRequestBody, reqEditors ...execdapi.RequestEditorFn) (*execdapi.ReplaceContentResponse, error)
	RenameFilesWithResponse(ctx context.Context, body execdapi.RenameFilesJSONRequestBody, reqEditors ...execdapi.RequestEditorFn) (*execdapi.RenameFilesResponse, error)
	ChmodFilesWithResponse(ctx context.Context, body execdapi.ChmodFilesJSONRequestBody, reqEditors ...execdapi.RequestEditorFn) (*execdapi.ChmodFilesResponse, error)
	MakeDirsWithResponse(ctx context.Context, body execdapi.MakeDirsJSONRequestBody, reqEditors ...execdapi.RequestEditorFn) (*execdapi.MakeDirsResponse, error)
	RemoveDirsWithResponse(ctx context.Context, params *execdapi.RemoveDirsParams, reqEditors ...execdapi.RequestEditorFn) (*execdapi.RemoveDirsResponse, error)
}

type FilesystemAdapter struct {
	client           FilesystemClient
	baseURL          string
	connectionConfig *config.ConnectionConfig
}

func NewFilesystemAdapter(client FilesystemClient, baseURL string, connectionConfig *config.ConnectionConfig) *FilesystemAdapter {
	return &FilesystemAdapter{
		client:           client,
		baseURL:          baseURL,
		connectionConfig: connectionConfig,
	}
}

func (a *FilesystemAdapter) GetFileInfo(ctx context.Context, paths []string) (map[string]models.EntryInfo, error) {
	resp, err := a.client.GetFilesInfoWithResponse(ctx, &execdapi.GetFilesInfoParams{Path: append([]string(nil), paths...)})
	if err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, NormalizeHTTPError(fmt.Errorf("get file info failed: unexpected status %d", resp.StatusCode()), resp.HTTPResponse)
	}
	return convert.FromExecdFileInfoMap(resp.JSON200), nil
}

func (a *FilesystemAdapter) Search(ctx context.Context, entry models.SearchEntry) ([]models.EntryInfo, error) {
	params := &execdapi.SearchFilesParams{Path: entry.Path}
	if entry.Pattern != "" {
		params.Pattern = &entry.Pattern
	}
	resp, err := a.client.SearchFilesWithResponse(ctx, params)
	if err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, NormalizeHTTPError(fmt.Errorf("search files failed: unexpected status %d", resp.StatusCode()), resp.HTTPResponse)
	}
	return convert.FromExecdFileInfoList(resp.JSON200), nil
}

func (a *FilesystemAdapter) CreateDirectories(ctx context.Context, entries []models.WriteEntry) error {
	body := make(execdapi.MakeDirsJSONRequestBody, len(entries))
	for _, entry := range entries {
		body[entry.Path] = convert.ToExecdPermission(entry.Mode, entry.Owner, entry.Group)
	}
	resp, err := a.client.MakeDirsWithResponse(ctx, body)
	if err != nil {
		return err
	}
	if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
		return nil
	}
	return NormalizeHTTPError(fmt.Errorf("create directories failed: unexpected status %d", resp.StatusCode()), resp.HTTPResponse)
}

func (a *FilesystemAdapter) DeleteDirectories(ctx context.Context, paths []string) error {
	resp, err := a.client.RemoveDirsWithResponse(ctx, &execdapi.RemoveDirsParams{Path: append([]string(nil), paths...)})
	if err != nil {
		return err
	}
	if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
		return nil
	}
	return NormalizeHTTPError(fmt.Errorf("delete directories failed: unexpected status %d", resp.StatusCode()), resp.HTTPResponse)
}

func (a *FilesystemAdapter) WriteFiles(ctx context.Context, entries []models.WriteEntry) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for _, entry := range entries {
		if err := addUploadEntry(writer, entry); err != nil {
			return err
		}
	}
	if err := writer.Close(); err != nil {
		return err
	}

	req, err := a.newDirectRequest(ctx, http.MethodPost, "/files/upload", nil, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := a.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	return NormalizeHTTPError(fmt.Errorf("write files failed: unexpected status %d", resp.StatusCode), resp)
}

func (a *FilesystemAdapter) ReadFile(ctx context.Context, path string, opts *models.ReadFileOptions) (string, error) {
	content, err := a.ReadBytes(ctx, path, opts)
	if err != nil {
		return "", err
	}

	encoding := "utf-8"
	if opts != nil && opts.Encoding != "" {
		encoding = opts.Encoding
	}
	if encoding != "utf-8" {
		return "", fmt.Errorf("unsupported encoding: %s", encoding)
	}
	return string(content), nil
}

func (a *FilesystemAdapter) ReadBytes(ctx context.Context, path string, opts *models.ReadFileOptions) ([]byte, error) {
	query := url.Values{}
	query.Set("path", path)

	req, err := a.newDirectRequest(ctx, http.MethodGet, "/files/download", query, nil)
	if err != nil {
		return nil, err
	}
	if opts != nil && opts.Range != "" {
		req.Header.Set("Range", opts.Range)
	}

	resp, err := a.httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, NormalizeHTTPError(fmt.Errorf("read file failed: unexpected status %d", resp.StatusCode), resp)
	}
	return io.ReadAll(resp.Body)
}

func (a *FilesystemAdapter) DeleteFiles(ctx context.Context, paths []string) error {
	resp, err := a.client.RemoveFilesWithResponse(ctx, &execdapi.RemoveFilesParams{Path: append([]string(nil), paths...)})
	if err != nil {
		return err
	}
	if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
		return nil
	}
	return NormalizeHTTPError(fmt.Errorf("delete files failed: unexpected status %d", resp.StatusCode()), resp.HTTPResponse)
}

func (a *FilesystemAdapter) MoveFiles(ctx context.Context, entries []models.MoveEntry) error {
	body := make(execdapi.RenameFilesJSONRequestBody, 0, len(entries))
	for _, entry := range entries {
		body = append(body, execdapi.RenameFileItem{
			Src:  entry.Src,
			Dest: entry.Dest,
		})
	}
	resp, err := a.client.RenameFilesWithResponse(ctx, body)
	if err != nil {
		return err
	}
	if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
		return nil
	}
	return NormalizeHTTPError(fmt.Errorf("move files failed: unexpected status %d", resp.StatusCode()), resp.HTTPResponse)
}

func (a *FilesystemAdapter) ReplaceContents(ctx context.Context, entries []models.ContentReplaceEntry) error {
	body := make(execdapi.ReplaceContentJSONRequestBody, len(entries))
	for _, entry := range entries {
		body[entry.Path] = execdapi.ReplaceFileContentItem{
			Old: entry.OldContent,
			New: entry.NewContent,
		}
	}
	resp, err := a.client.ReplaceContentWithResponse(ctx, body)
	if err != nil {
		return err
	}
	if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
		return nil
	}
	return NormalizeHTTPError(fmt.Errorf("replace contents failed: unexpected status %d", resp.StatusCode()), resp.HTTPResponse)
}

func (a *FilesystemAdapter) SetPermissions(ctx context.Context, entries []models.SetPermissionEntry) error {
	body := make(execdapi.ChmodFilesJSONRequestBody, len(entries))
	for _, entry := range entries {
		body[entry.Path] = convert.ToExecdPermission(entry.Mode, entry.Owner, entry.Group)
	}
	resp, err := a.client.ChmodFilesWithResponse(ctx, body)
	if err != nil {
		return err
	}
	if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
		return nil
	}
	return NormalizeHTTPError(fmt.Errorf("set permissions failed: unexpected status %d", resp.StatusCode()), resp.HTTPResponse)
}

func (a *FilesystemAdapter) newDirectRequest(ctx context.Context, method, path string, query url.Values, body io.Reader) (*http.Request, error) {
	target := a.baseURL + path
	if len(query) > 0 {
		target += "?" + query.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, method, target, body)
	if err != nil {
		return nil, err
	}
	applyFilesystemConnectionHeaders(req, a.connectionConfig)
	return req, nil
}

func (a *FilesystemAdapter) httpClient() *http.Client {
	if a.connectionConfig != nil && a.connectionConfig.HTTPClient != nil {
		return a.connectionConfig.HTTPClient
	}
	return http.DefaultClient
}

func addUploadEntry(writer *multipart.Writer, entry models.WriteEntry) error {
	metadata := map[string]any{
		"path":  entry.Path,
		"mode":  entry.Mode,
		"owner": entry.Owner,
		"group": entry.Group,
	}
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	metadataHeader := textproto.MIMEHeader{}
	metadataHeader.Set("Content-Disposition", `form-data; name="metadata"; filename="metadata"`)
	metadataHeader.Set("Content-Type", "application/json")
	metadataPart, err := writer.CreatePart(metadataHeader)
	if err != nil {
		return err
	}
	if _, err := metadataPart.Write(metadataBytes); err != nil {
		return err
	}

	fileHeader := textproto.MIMEHeader{}
	fileHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename=%q`, entry.Path))
	fileHeader.Set("Content-Type", uploadContentType(entry))
	filePart, err := writer.CreatePart(fileHeader)
	if err != nil {
		return err
	}
	switch data := entry.Data.(type) {
	case string:
		_, err = io.WriteString(filePart, data)
	case []byte:
		_, err = filePart.Write(data)
	case io.Reader:
		_, err = io.Copy(filePart, data)
	case nil:
		err = nil
	default:
		return fmt.Errorf("unsupported write entry data type %T", entry.Data)
	}
	return err
}

func uploadContentType(entry models.WriteEntry) string {
	switch entry.Data.(type) {
	case string:
		encoding := entry.Encoding
		if encoding == "" {
			encoding = "utf-8"
		}
		return "text/plain; charset=" + encoding
	default:
		return "application/octet-stream"
	}
}

func applyFilesystemConnectionHeaders(req *http.Request, connectionConfig *config.ConnectionConfig) {
	if connectionConfig == nil {
		return
	}
	if connectionConfig.APIKey != "" {
		req.Header.Set("OPEN-SANDBOX-API-KEY", connectionConfig.APIKey)
	}
	for k, v := range connectionConfig.Headers {
		req.Header.Set(k, v)
	}
}
