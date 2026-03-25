# Go Sandbox SDK Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a Go Sandbox SDK that matches the existing Python, JavaScript, Kotlin, and C# SDKs in capabilities and core concepts, including `Sandbox`, `SandboxManager`, readiness handling, execd command streaming, filesystem operations, metrics, egress policy, and lifecycle management.

**Architecture:** The Go SDK will follow the repository convention of separating generated OpenAPI transport from handwritten SDK business logic. Generated lifecycle/execd/egress clients will live under `internal/openapi`, adapters will convert wire models into stable Go SDK models, and the public API will expose high-level `Sandbox` and `SandboxManager` entry points with Go-style `context.Context` and `Close()` semantics.

**Tech Stack:** Go, `net/http`, OpenAPI code generation for Go, Server-Sent Events parsing, `testing`, `httptest`

---

## File Structure

### New directories

- `sdks/sandbox/go/`
- `sdks/sandbox/go/cmd/generate-openapi/`
- `sdks/sandbox/go/sandbox/`
- `sdks/sandbox/go/sandbox/config/`
- `sdks/sandbox/go/sandbox/models/`
- `sdks/sandbox/go/sandbox/services/`
- `sdks/sandbox/go/sandbox/adapters/`
- `sdks/sandbox/go/sandbox/factory/`
- `sdks/sandbox/go/sandbox/errors/`
- `sdks/sandbox/go/sandbox/internal/openapi/lifecycle/`
- `sdks/sandbox/go/sandbox/internal/openapi/execd/`
- `sdks/sandbox/go/sandbox/internal/openapi/egress/`
- `sdks/sandbox/go/sandbox/internal/sse/`
- `sdks/sandbox/go/sandbox/internal/transport/`
- `sdks/sandbox/go/sandbox/internal/convert/`
- `sdks/sandbox/go/examples/quickstart/`
- `sdks/sandbox/go/examples/streaming/`
- `sdks/sandbox/go/examples/manager/`
- `sdks/sandbox/go/tests/unit/`
- `sdks/sandbox/go/tests/integration/`
- `sdks/sandbox/go/tests/e2e/`

### Public API files

- Create: `sdks/sandbox/go/sandbox/doc.go`
- Create: `sdks/sandbox/go/sandbox/exports.go`
- Create: `sdks/sandbox/go/sandbox/sandbox.go`
- Create: `sdks/sandbox/go/sandbox/sandbox_manager.go`
- Create: `sdks/sandbox/go/sandbox/options.go`
- Create: `sdks/sandbox/go/sandbox/readiness.go`

### Config files

- Create: `sdks/sandbox/go/sandbox/config/connection.go`
- Create: `sdks/sandbox/go/sandbox/config/defaults.go`

### Model files

- Create: `sdks/sandbox/go/sandbox/models/sandboxes.go`
- Create: `sdks/sandbox/go/sandbox/models/execd.go`
- Create: `sdks/sandbox/go/sandbox/models/execution.go`
- Create: `sdks/sandbox/go/sandbox/models/filesystem.go`
- Create: `sdks/sandbox/go/sandbox/models/egress.go`

### Service interface files

- Create: `sdks/sandbox/go/sandbox/services/sandboxes.go`
- Create: `sdks/sandbox/go/sandbox/services/commands.go`
- Create: `sdks/sandbox/go/sandbox/services/filesystem.go`
- Create: `sdks/sandbox/go/sandbox/services/health.go`
- Create: `sdks/sandbox/go/sandbox/services/metrics.go`
- Create: `sdks/sandbox/go/sandbox/services/egress.go`

### Adapter and conversion files

- Create: `sdks/sandbox/go/sandbox/adapters/sandboxes_adapter.go`
- Create: `sdks/sandbox/go/sandbox/adapters/commands_adapter.go`
- Create: `sdks/sandbox/go/sandbox/adapters/filesystem_adapter.go`
- Create: `sdks/sandbox/go/sandbox/adapters/health_adapter.go`
- Create: `sdks/sandbox/go/sandbox/adapters/metrics_adapter.go`
- Create: `sdks/sandbox/go/sandbox/adapters/egress_adapter.go`
- Create: `sdks/sandbox/go/sandbox/adapters/errors.go`
- Create: `sdks/sandbox/go/sandbox/internal/convert/sandboxes.go`
- Create: `sdks/sandbox/go/sandbox/internal/convert/execd.go`
- Create: `sdks/sandbox/go/sandbox/internal/convert/filesystem.go`
- Create: `sdks/sandbox/go/sandbox/internal/convert/egress.go`

### Transport and factory files

- Create: `sdks/sandbox/go/sandbox/factory/adapter_factory.go`
- Create: `sdks/sandbox/go/sandbox/factory/default_adapter_factory.go`
- Create: `sdks/sandbox/go/sandbox/internal/transport/http_client_provider.go`
- Create: `sdks/sandbox/go/sandbox/internal/transport/base_url.go`
- Create: `sdks/sandbox/go/sandbox/internal/transport/headers.go`

### SSE and error files

- Create: `sdks/sandbox/go/sandbox/internal/sse/parser.go`
- Create: `sdks/sandbox/go/sandbox/internal/sse/event_stream.go`
- Create: `sdks/sandbox/go/sandbox/errors/errors.go`

### Build and generation files

- Create: `sdks/sandbox/go/go.mod`
- Create: `sdks/sandbox/go/README.md`
- Create: `sdks/sandbox/go/Makefile`
- Create: `sdks/sandbox/go/cmd/generate-openapi/main.go`
- Create: `sdks/sandbox/go/scripts/generate_api.sh`
- Create: `sdks/sandbox/go/tools.go`

### Examples and tests

- Create: `sdks/sandbox/go/examples/quickstart/main.go`
- Create: `sdks/sandbox/go/examples/streaming/main.go`
- Create: `sdks/sandbox/go/examples/manager/main.go`
- Create: `sdks/sandbox/go/tests/unit/connection_config_test.go`
- Create: `sdks/sandbox/go/tests/unit/errors_test.go`
- Create: `sdks/sandbox/go/tests/unit/readiness_test.go`
- Create: `sdks/sandbox/go/tests/unit/sse_parser_test.go`
- Create: `sdks/sandbox/go/tests/unit/commands_adapter_test.go`
- Create: `sdks/sandbox/go/tests/unit/filesystem_adapter_test.go`
- Create: `sdks/sandbox/go/tests/unit/sandboxes_adapter_test.go`
- Create: `sdks/sandbox/go/tests/unit/models_test.go`
- Create: `sdks/sandbox/go/tests/integration/http_stack_test.go`
- Create: `sdks/sandbox/go/tests/e2e/sandbox_e2e_test.go`
- Create: `sdks/sandbox/go/tests/e2e/sandbox_manager_e2e_test.go`
- Create: `sdks/sandbox/go/tests/e2e/code_interpreter_e2e_test.go`

## Milestones

1. Scaffold the Go SDK package, generation pipeline, and public model/config surface.
2. Implement lifecycle transport and the public `Sandbox` / `SandboxManager` orchestration.
3. Implement execd adapters for commands, filesystem, health, and metrics.
4. Implement readiness, SSE streaming, egress policy, and error normalization.
5. Add examples, unit tests, integration tests, and E2E parity coverage.
6. Polish documentation and package ergonomics for external users.

## Current Execution Status

- Completed: Task 1, Task 2, Task 3, Task 4
- Completed with initial scope: Task 6
- In progress: Task 5, Task 7, Task 8, Task 14

### Progress Notes

- Landed: Go module bootstrap, generator entrypoint, `Makefile`, tools pinning, and initial unit tests under `sdks/sandbox/go/`.
- Landed: public `ConnectionConfig`, stable lifecycle models, top-level exports, service interfaces, adapter factory contracts, transport helpers, and public error types.
- Landed: initial `SandboxesAdapter`, initial `SandboxManager`, and initial `Sandbox` lifecycle shell (`GetInfo`, `GetEndpoint`, `Renew`, `Pause`, `Kill`, `Close`).
- Landed: generated clients under `sdks/sandbox/go/sandbox/internal/openapi/{lifecycle,execd,egress}` and `go.mod` runtime dependency for generated code.
- Added during execution: generation-time lifecycle spec preprocessing that rewrites OpenAPI 3.1 `oneOf + null` patterns into a temporary 3.0.3-compatible spec for `oapi-codegen`.
- Current blocker removed: `oapi-codegen` is installed and generation now succeeds locally, but default factory wiring still needs to be moved from fake/injected clients to real generated lifecycle transport.
- Remaining high-value path: connect generated lifecycle client to `DefaultAdapterFactory`, then replace placeholder lifecycle wiring in `SandboxManager` / `Sandbox` create-connect-resume flows before moving to execd/egress adapters.

## Task Plan

### Task 1: Bootstrap the Go SDK module and generation pipeline

**Status:** Completed

**Files:**
- Create: `sdks/sandbox/go/go.mod`
- Create: `sdks/sandbox/go/Makefile`
- Create: `sdks/sandbox/go/tools.go`
- Create: `sdks/sandbox/go/cmd/generate-openapi/main.go`
- Create: `sdks/sandbox/go/scripts/generate_api.sh`
- Create: `sdks/sandbox/go/sandbox/internal/openapi/lifecycle/.gitkeep`
- Create: `sdks/sandbox/go/sandbox/internal/openapi/execd/.gitkeep`
- Create: `sdks/sandbox/go/sandbox/internal/openapi/egress/.gitkeep`
- Test: `sdks/sandbox/go/Makefile`

- [ ] **Step 1: Create the failing generator bootstrap check**

```go
package main

import "testing"

func TestGeneratorSpecsExist(t *testing.T) {
    t.Fatal("generator bootstrap not implemented")
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd sdks/sandbox/go && go test ./...`
Expected: FAIL with missing generator bootstrap implementation

- [ ] **Step 3: Add module, tool pinning, and generator entrypoint**

```go
module github.com/alibaba/opensandbox/sdks/sandbox/go

go 1.25
```

```go
//go:build tools

package tools

import _ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
```

```go
package main

func main() {
    // Resolve repo-root specs and invoke generation for lifecycle, execd, and egress.
}
```

- [ ] **Step 4: Add repeatable generate commands**

```makefile
generate:
	go run ./cmd/generate-openapi

test:
	go test ./...
```

- [ ] **Step 5: Run test to verify the module is wired**

Run: `cd sdks/sandbox/go && go test ./...`
Expected: PASS or FAIL only on not-yet-created packages, not on module/bootstrap wiring

- [ ] **Step 6: Commit**

```bash
git add sdks/sandbox/go
git commit -m "feat(go-sdk): bootstrap module and api generation"
```

### Task 2: Define public models, defaults, and connection configuration

**Status:** Completed for the initial lifecycle/config surface. Additional model files for execd/filesystem/egress are still pending.

**Files:**
- Create: `sdks/sandbox/go/sandbox/doc.go`
- Create: `sdks/sandbox/go/sandbox/exports.go`
- Create: `sdks/sandbox/go/sandbox/options.go`
- Create: `sdks/sandbox/go/sandbox/config/connection.go`
- Create: `sdks/sandbox/go/sandbox/config/defaults.go`
- Create: `sdks/sandbox/go/sandbox/models/sandboxes.go`
- Create: `sdks/sandbox/go/sandbox/models/execd.go`
- Create: `sdks/sandbox/go/sandbox/models/execution.go`
- Create: `sdks/sandbox/go/sandbox/models/filesystem.go`
- Create: `sdks/sandbox/go/sandbox/models/egress.go`
- Test: `sdks/sandbox/go/tests/unit/connection_config_test.go`
- Test: `sdks/sandbox/go/tests/unit/models_test.go`

- [ ] **Step 1: Write failing config and model tests**

```go
func TestDefaultConnectionConfigUsesLocalhost(t *testing.T) {
    cfg := config.DefaultConnectionConfig()
    if got := cfg.BaseURL(); got != "http://localhost:8080" {
        t.Fatalf("unexpected base url: %s", got)
    }
}

func TestNilTimeoutMeansManualCleanup(t *testing.T) {
    var timeout *time.Duration
    opts := sandbox.SandboxCreateOptions{Timeout: timeout}
    if opts.Timeout != nil {
        t.Fatal("expected nil timeout")
    }
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd sdks/sandbox/go && go test ./tests/unit -run 'Test(DefaultConnectionConfigUsesLocalhost|NilTimeoutMeansManualCleanup)'`
Expected: FAIL with missing config/model types

- [ ] **Step 3: Implement defaults and stable public model types**

```go
const (
    DefaultDomain      = "localhost:8080"
    DefaultProtocol    = ProtocolHTTP
    DefaultExecdPort   = 44772
    DefaultEgressPort  = 18080
)
```

```go
type ConnectionConfig struct {
    Domain         string
    Protocol       Protocol
    APIKey         string
    Headers        map[string]string
    RequestTimeout time.Duration
    UseServerProxy bool
    HTTPClient     *http.Client
    SSEHTTPClient  *http.Client
}
```

- [ ] **Step 4: Export the top-level package surface**

```go
package sandbox

type SandboxCreateOptions = models.SandboxCreateOptions
```

- [ ] **Step 5: Run targeted tests**

Run: `cd sdks/sandbox/go && go test ./tests/unit -run 'Test(DefaultConnectionConfigUsesLocalhost|NilTimeoutMeansManualCleanup)'`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add sdks/sandbox/go
git commit -m "feat(go-sdk): add public config and model surface"
```

### Task 3: Define service interfaces and adapter factory boundaries

**Status:** Completed for the initial lifecycle-oriented interfaces and stack boundaries.

**Files:**
- Create: `sdks/sandbox/go/sandbox/services/sandboxes.go`
- Create: `sdks/sandbox/go/sandbox/services/commands.go`
- Create: `sdks/sandbox/go/sandbox/services/filesystem.go`
- Create: `sdks/sandbox/go/sandbox/services/health.go`
- Create: `sdks/sandbox/go/sandbox/services/metrics.go`
- Create: `sdks/sandbox/go/sandbox/services/egress.go`
- Create: `sdks/sandbox/go/sandbox/factory/adapter_factory.go`
- Create: `sdks/sandbox/go/sandbox/factory/default_adapter_factory.go`
- Test: `sdks/sandbox/go/tests/unit/models_test.go`

- [ ] **Step 1: Write the failing interface compilation test**

```go
func TestFactoryExposesLifecycleExecdAndEgressStacks(t *testing.T) {
    var _ factory.AdapterFactory = (*factory.DefaultAdapterFactory)(nil)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd sdks/sandbox/go && go test ./tests/unit -run TestFactoryExposesLifecycleExecdAndEgressStacks`
Expected: FAIL with undefined factory symbols

- [ ] **Step 3: Define interface contracts**

```go
type AdapterFactory interface {
    CreateLifecycleStack(opts CreateLifecycleStackOptions) (*LifecycleStack, error)
    CreateExecdStack(opts CreateExecdStackOptions) (*ExecdStack, error)
    CreateEgressStack(opts CreateEgressStackOptions) (*EgressStack, error)
}
```

- [ ] **Step 4: Add placeholder default factory**

```go
type DefaultAdapterFactory struct{}
```

- [ ] **Step 5: Run targeted test**

Run: `cd sdks/sandbox/go && go test ./tests/unit -run TestFactoryExposesLifecycleExecdAndEgressStacks`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add sdks/sandbox/go
git commit -m "feat(go-sdk): define service and factory boundaries"
```

### Task 4: Implement internal transport helpers and HTTP client ownership rules

**Status:** Completed

**Files:**
- Create: `sdks/sandbox/go/sandbox/internal/transport/http_client_provider.go`
- Create: `sdks/sandbox/go/sandbox/internal/transport/base_url.go`
- Create: `sdks/sandbox/go/sandbox/internal/transport/headers.go`
- Modify: `sdks/sandbox/go/sandbox/config/connection.go`
- Test: `sdks/sandbox/go/tests/unit/connection_config_test.go`
- Test: `sdks/sandbox/go/tests/integration/http_stack_test.go`

- [ ] **Step 1: Write failing transport tests**

```go
func TestMergeHeadersPrefersEndpointHeaders(t *testing.T) {}
func TestCloseOnlyClosesOwnedClients(t *testing.T) {}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd sdks/sandbox/go && go test ./tests/unit ./tests/integration -run 'Test(MergeHeadersPrefersEndpointHeaders|CloseOnlyClosesOwnedClients)'`
Expected: FAIL with missing transport helpers

- [ ] **Step 3: Implement HTTP client provider and header merging**

```go
func MergeHeaders(base, endpoint map[string]string) map[string]string {
    // endpoint overrides base keys
}
```

- [ ] **Step 4: Implement config cloning and local cleanup behavior**

```go
func (c *ConnectionConfig) CloneWithHTTPClientIfMissing() *ConnectionConfig
func (c *ConnectionConfig) Close() error
```

- [ ] **Step 5: Run targeted tests**

Run: `cd sdks/sandbox/go && go test ./tests/unit ./tests/integration -run 'Test(MergeHeadersPrefersEndpointHeaders|CloseOnlyClosesOwnedClients)'`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add sdks/sandbox/go
git commit -m "feat(go-sdk): add transport helpers and client lifecycle"
```

### Task 5: Generate lifecycle, execd, and egress OpenAPI clients

**Status:** In progress

**Execution Note:** The original `oapi-codegen` invocation failed on OpenAPI 3.1 nullable unions in `specs/sandbox-lifecycle.yml`. The generator has been updated to preprocess the lifecycle spec into a temporary 3.0.3-compatible file for code generation. Generated clients now exist on disk and compile after adding `github.com/oapi-codegen/runtime`, but the default factory is not yet wired to use them.

**Files:**
- Modify: `sdks/sandbox/go/cmd/generate-openapi/main.go`
- Create: `sdks/sandbox/go/sandbox/internal/openapi/lifecycle/*.go`
- Create: `sdks/sandbox/go/sandbox/internal/openapi/execd/*.go`
- Create: `sdks/sandbox/go/sandbox/internal/openapi/egress/*.go`
- Test: `sdks/sandbox/go/tests/integration/http_stack_test.go`

- [ ] **Step 1: Write failing generation smoke test**

```go
func TestGeneratedPackagesBuild(t *testing.T) {
    t.Fatal("generated packages not available")
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd sdks/sandbox/go && go test ./tests/integration -run TestGeneratedPackagesBuild`
Expected: FAIL with missing generated packages

- [ ] **Step 3: Implement generator for all three specs**

```go
// generate lifecycle from specs/sandbox-lifecycle.yml
// generate execd from specs/execd-api.yaml
// generate egress from specs/egress-api.yaml
```

- [ ] **Step 4: Run generation and format**

Run: `cd sdks/sandbox/go && make generate`
Expected: generated `.go` files under `sandbox/internal/openapi/...`

- [ ] **Step 5: Run smoke test**

Run: `cd sdks/sandbox/go && go test ./tests/integration -run TestGeneratedPackagesBuild`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add sdks/sandbox/go
git commit -m "feat(go-sdk): generate openapi transport clients"
```

### Task 6: Implement error normalization and public error types

**Status:** Completed for the initial HTTP/status/request-id normalization path

**Files:**
- Create: `sdks/sandbox/go/sandbox/errors/errors.go`
- Create: `sdks/sandbox/go/sandbox/adapters/errors.go`
- Test: `sdks/sandbox/go/tests/unit/errors_test.go`

- [ ] **Step 1: Write failing error mapping tests**

```go
func TestSandboxErrorCarriesRequestID(t *testing.T) {}
func TestReadyTimeoutErrorMatchesErrorsAs(t *testing.T) {}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd sdks/sandbox/go && go test ./tests/unit -run 'Test(SandboxErrorCarriesRequestID|ReadyTimeoutErrorMatchesErrorsAs)'`
Expected: FAIL with missing error types

- [ ] **Step 3: Implement public errors**

```go
type SandboxError struct {
    Code string
    Message string
    RequestID string
    StatusCode int
    Cause error
}
```

- [ ] **Step 4: Implement adapter-side HTTP/OpenAPI to SDK error conversion**

```go
func NormalizeError(err error, resp *http.Response) error
```

- [ ] **Step 5: Run targeted tests**

Run: `cd sdks/sandbox/go && go test ./tests/unit -run 'Test(SandboxErrorCarriesRequestID|ReadyTimeoutErrorMatchesErrorsAs)'`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add sdks/sandbox/go
git commit -m "feat(go-sdk): add public errors and adapter normalization"
```

### Task 7: Implement lifecycle conversion and sandboxes adapter

**Status:** In progress

**Execution Note:** An initial handwritten `SandboxesAdapter` and conversion layer are in place and covered by unit tests. The remaining work is to swap the adapter from the current interface-driven fake client wiring to the real generated lifecycle client.

**Files:**
- Create: `sdks/sandbox/go/sandbox/internal/convert/sandboxes.go`
- Create: `sdks/sandbox/go/sandbox/adapters/sandboxes_adapter.go`
- Test: `sdks/sandbox/go/tests/unit/sandboxes_adapter_test.go`

- [ ] **Step 1: Write failing lifecycle adapter tests**

```go
func TestSandboxesAdapterConvertsCreateResponse(t *testing.T) {}
func TestSandboxesAdapterPassesUseServerProxyToEndpointLookup(t *testing.T) {}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd sdks/sandbox/go && go test ./tests/unit -run 'TestSandboxesAdapter'`
Expected: FAIL with missing adapter implementation

- [ ] **Step 3: Implement lifecycle model conversion**

```go
func ToSandboxInfo(apiResp any) (*models.SandboxInfo, error)
```

- [ ] **Step 4: Implement lifecycle adapter methods**

```go
type SandboxesAdapter struct {
    // wraps generated lifecycle client
}
```

- [ ] **Step 5: Run targeted tests**

Run: `cd sdks/sandbox/go && go test ./tests/unit -run 'TestSandboxesAdapter'`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add sdks/sandbox/go
git commit -m "feat(go-sdk): implement lifecycle adapter"
```

### Task 8: Implement `SandboxManager`

**Status:** In progress

**Execution Note:** `SandboxManager` exists and its renew/get/list/pause/resume/kill/close behavior is tested against the service interface. It still needs real generated lifecycle transport wiring through `DefaultAdapterFactory`.

**Files:**
- Create: `sdks/sandbox/go/sandbox/sandbox_manager.go`
- Modify: `sdks/sandbox/go/sandbox/exports.go`
- Test: `sdks/sandbox/go/tests/unit/sandboxes_adapter_test.go`
- Test: `sdks/sandbox/go/tests/e2e/sandbox_manager_e2e_test.go`

- [ ] **Step 1: Write failing manager tests**

```go
func TestManagerRenewUsesNowPlusTimeout(t *testing.T) {}
func TestManagerCloseReleasesLocalResourcesOnly(t *testing.T) {}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd sdks/sandbox/go && go test ./tests/unit -run 'TestManager'`
Expected: FAIL with missing manager implementation

- [ ] **Step 3: Implement `NewSandboxManager` and lifecycle delegation**

```go
func NewSandboxManager(opts SandboxManagerOptions) (*SandboxManager, error)
```

- [ ] **Step 4: Implement renew, list, get, pause, resume, kill, close**

```go
func (m *SandboxManager) RenewSandbox(ctx context.Context, sandboxID string, timeout time.Duration) error
```

- [ ] **Step 5: Run targeted tests**

Run: `cd sdks/sandbox/go && go test ./tests/unit -run 'TestManager'`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add sdks/sandbox/go
git commit -m "feat(go-sdk): add sandbox manager"
```

### Task 9: Implement health and metrics adapters

**Files:**
- Create: `sdks/sandbox/go/sandbox/adapters/health_adapter.go`
- Create: `sdks/sandbox/go/sandbox/adapters/metrics_adapter.go`
- Create: `sdks/sandbox/go/sandbox/internal/convert/execd.go`
- Test: `sdks/sandbox/go/tests/unit/models_test.go`

- [ ] **Step 1: Write failing health and metrics tests**

```go
func TestHealthAdapterPingReturnsTrueOn200(t *testing.T) {}
func TestMetricsAdapterConvertsTimestampAndResourceNumbers(t *testing.T) {}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd sdks/sandbox/go && go test ./tests/unit -run 'Test(HealthAdapter|MetricsAdapter)'`
Expected: FAIL with missing health/metrics adapters

- [ ] **Step 3: Implement execd health adapter**

```go
func (a *HealthAdapter) Ping(ctx context.Context) (bool, error)
```

- [ ] **Step 4: Implement metrics conversion and adapter**

```go
func (a *MetricsAdapter) GetMetrics(ctx context.Context) (*models.SandboxMetrics, error)
```

- [ ] **Step 5: Run targeted tests**

Run: `cd sdks/sandbox/go && go test ./tests/unit -run 'Test(HealthAdapter|MetricsAdapter)'`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add sdks/sandbox/go
git commit -m "feat(go-sdk): add health and metrics adapters"
```

### Task 10: Implement filesystem conversion and adapter

**Files:**
- Create: `sdks/sandbox/go/sandbox/internal/convert/filesystem.go`
- Create: `sdks/sandbox/go/sandbox/adapters/filesystem_adapter.go`
- Test: `sdks/sandbox/go/tests/unit/filesystem_adapter_test.go`

- [ ] **Step 1: Write failing filesystem adapter tests**

```go
func TestFilesystemAdapterWriteAndReadFile(t *testing.T) {}
func TestFilesystemAdapterSearchAndDelete(t *testing.T) {}
func TestFilesystemAdapterReplaceContentMapsEntries(t *testing.T) {}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd sdks/sandbox/go && go test ./tests/unit -run 'TestFilesystemAdapter'`
Expected: FAIL with missing filesystem adapter

- [ ] **Step 3: Implement filesystem model conversion**

```go
func ToEntryInfoList(apiResp any) ([]models.EntryInfo, error)
```

- [ ] **Step 4: Implement filesystem adapter methods**

```go
func (a *FilesystemAdapter) ReadFile(ctx context.Context, path string) (string, error)
```

- [ ] **Step 5: Run targeted tests**

Run: `cd sdks/sandbox/go && go test ./tests/unit -run 'TestFilesystemAdapter'`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add sdks/sandbox/go
git commit -m "feat(go-sdk): add filesystem adapter"
```

### Task 11: Implement SSE parser and command streaming primitives

**Files:**
- Create: `sdks/sandbox/go/sandbox/internal/sse/parser.go`
- Create: `sdks/sandbox/go/sandbox/internal/sse/event_stream.go`
- Test: `sdks/sandbox/go/tests/unit/sse_parser_test.go`

- [ ] **Step 1: Write failing SSE parser tests**

```go
func TestParserDecodesJSONEvents(t *testing.T) {}
func TestParserReturnsErrorOnNonOKResponse(t *testing.T) {}
func TestParserStopsOnContextCancel(t *testing.T) {}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd sdks/sandbox/go && go test ./tests/unit -run 'TestParser'`
Expected: FAIL with missing parser implementation

- [ ] **Step 3: Implement line-oriented SSE parser**

```go
type Event struct {
    Event string
    Data  []byte
}
```

- [ ] **Step 4: Implement stream wrapper with events and terminal error**

```go
type Stream struct {
    Events <-chan Event
    Done   <-chan error
}
```

- [ ] **Step 5: Run targeted tests**

Run: `cd sdks/sandbox/go && go test ./tests/unit -run 'TestParser'`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add sdks/sandbox/go
git commit -m "feat(go-sdk): add sse parser"
```

### Task 12: Implement commands adapter and execution aggregation

**Files:**
- Create: `sdks/sandbox/go/sandbox/adapters/commands_adapter.go`
- Modify: `sdks/sandbox/go/sandbox/internal/convert/execd.go`
- Test: `sdks/sandbox/go/tests/unit/commands_adapter_test.go`

- [ ] **Step 1: Write failing commands adapter tests**

```go
func TestCommandsAdapterRunAggregatesStdoutStderrAndComplete(t *testing.T) {}
func TestCommandsAdapterRunStreamExposesOrderedEvents(t *testing.T) {}
func TestCommandsAdapterBackgroundLogCursorParsing(t *testing.T) {}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd sdks/sandbox/go && go test ./tests/unit -run 'TestCommandsAdapter'`
Expected: FAIL with missing commands adapter

- [ ] **Step 3: Implement event conversion and command request mapping**

```go
func ToRunCommandRequest(command string, opts *models.RunCommandOptions) any
```

- [ ] **Step 4: Implement `Run`, `RunStream`, `Interrupt`, `GetCommandStatus`, and `GetBackgroundCommandLogs`**

```go
func (a *CommandsAdapter) Run(ctx context.Context, command string, opts *models.RunCommandOptions) (*models.CommandExecution, error)
```

- [ ] **Step 5: Run targeted tests**

Run: `cd sdks/sandbox/go && go test ./tests/unit -run 'TestCommandsAdapter'`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add sdks/sandbox/go
git commit -m "feat(go-sdk): add command execution and streaming"
```

### Task 13: Implement egress adapter

**Files:**
- Create: `sdks/sandbox/go/sandbox/internal/convert/egress.go`
- Create: `sdks/sandbox/go/sandbox/adapters/egress_adapter.go`
- Test: `sdks/sandbox/go/tests/unit/models_test.go`

- [ ] **Step 1: Write failing egress tests**

```go
func TestEgressAdapterGetsPolicy(t *testing.T) {}
func TestEgressAdapterPatchesRulesWithoutReplacingDefaultAction(t *testing.T) {}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd sdks/sandbox/go && go test ./tests/unit -run 'TestEgressAdapter'`
Expected: FAIL with missing egress adapter

- [ ] **Step 3: Implement egress conversion**

```go
func ToNetworkPolicy(apiResp any) (*models.NetworkPolicy, error)
```

- [ ] **Step 4: Implement egress adapter methods**

```go
func (a *EgressAdapter) PatchRules(ctx context.Context, rules []models.NetworkRule) error
```

- [ ] **Step 5: Run targeted tests**

Run: `cd sdks/sandbox/go && go test ./tests/unit -run 'TestEgressAdapter'`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add sdks/sandbox/go
git commit -m "feat(go-sdk): add egress policy adapter"
```

### Task 14: Implement `Sandbox` orchestration for create, connect, resume, close, and remote actions

**Status:** In progress

**Execution Note:** The minimal `Sandbox` type and basic lifecycle instance methods are implemented. The create/connect/resume orchestration and readiness flow are still pending.

**Files:**
- Create: `sdks/sandbox/go/sandbox/sandbox.go`
- Create: `sdks/sandbox/go/sandbox/readiness.go`
- Modify: `sdks/sandbox/go/sandbox/options.go`
- Modify: `sdks/sandbox/go/sandbox/exports.go`
- Test: `sdks/sandbox/go/tests/unit/readiness_test.go`
- Test: `sdks/sandbox/go/tests/e2e/sandbox_e2e_test.go`

- [ ] **Step 1: Write failing sandbox orchestration tests**

```go
func TestCreateResolvesExecdAndEgressEndpoints(t *testing.T) {}
func TestCreateCleansUpRemoteSandboxOnInitializationFailure(t *testing.T) {}
func TestResumeReturnsFreshSandboxInstance(t *testing.T) {}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd sdks/sandbox/go && go test ./tests/unit -run 'Test(Create|Resume)'`
Expected: FAIL with missing sandbox orchestration

- [ ] **Step 3: Implement create/connect/resume flow**

```go
func Create(ctx context.Context, opts SandboxCreateOptions) (*Sandbox, error)
func Connect(ctx context.Context, opts SandboxConnectOptions) (*Sandbox, error)
func Resume(ctx context.Context, opts SandboxResumeOptions) (*Sandbox, error)
```

- [ ] **Step 4: Implement instance methods and local cleanup**

```go
func (s *Sandbox) Kill(ctx context.Context) error
func (s *Sandbox) Close() error
func (s *Sandbox) Pause(ctx context.Context) error
func (s *Sandbox) Renew(ctx context.Context, timeout time.Duration) (*models.SandboxRenewResponse, error)
```

- [ ] **Step 5: Run targeted tests**

Run: `cd sdks/sandbox/go && go test ./tests/unit -run 'Test(Create|Resume)'`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add sdks/sandbox/go
git commit -m "feat(go-sdk): add sandbox orchestration"
```

### Task 15: Implement readiness polling, default health check, and endpoint URL helpers

**Files:**
- Modify: `sdks/sandbox/go/sandbox/readiness.go`
- Modify: `sdks/sandbox/go/sandbox/sandbox.go`
- Test: `sdks/sandbox/go/tests/unit/readiness_test.go`

- [ ] **Step 1: Write failing readiness tests**

```go
func TestWaitUntilReadyPollsLifecycleThenHealth(t *testing.T) {}
func TestWaitUntilReadyUsesCustomHealthCheckWhenProvided(t *testing.T) {}
func TestWaitUntilReadyReturnsTimeoutError(t *testing.T) {}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd sdks/sandbox/go && go test ./tests/unit -run 'TestWaitUntilReady'`
Expected: FAIL with missing readiness logic

- [ ] **Step 3: Implement readiness state polling**

```go
func (s *Sandbox) WaitUntilReady(ctx context.Context, opts *WaitUntilReadyOptions) error
```

- [ ] **Step 4: Implement `IsHealthy`, `GetEndpoint`, and `GetEndpointURL` helpers**

```go
func (s *Sandbox) GetEndpointURL(ctx context.Context, port int) (string, error)
```

- [ ] **Step 5: Run targeted tests**

Run: `cd sdks/sandbox/go && go test ./tests/unit -run 'TestWaitUntilReady'`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add sdks/sandbox/go
git commit -m "feat(go-sdk): add readiness and endpoint helpers"
```

### Task 16: Add package docs and runnable examples

**Files:**
- Create: `sdks/sandbox/go/README.md`
- Create: `sdks/sandbox/go/examples/quickstart/main.go`
- Create: `sdks/sandbox/go/examples/streaming/main.go`
- Create: `sdks/sandbox/go/examples/manager/main.go`

- [ ] **Step 1: Write the quickstart example first**

```go
ctx := context.Background()
sbx, err := sandbox.Create(ctx, sandbox.SandboxCreateOptions{
    Image: models.ImageSpec{URI: "ubuntu"},
})
```

- [ ] **Step 2: Verify examples compile**

Run: `cd sdks/sandbox/go && go test ./...`
Expected: PASS with examples compiling as part of the module

- [ ] **Step 3: Write README matching existing SDK concepts**

```md
# OpenSandbox SDK for Go
```

- [ ] **Step 4: Verify formatting**

Run: `cd sdks/sandbox/go && gofmt -w $(find . -name '*.go' -not -path './sandbox/internal/openapi/*')`
Expected: all handwritten Go files formatted

- [ ] **Step 5: Commit**

```bash
git add sdks/sandbox/go
git commit -m "docs(go-sdk): add examples and readme"
```

### Task 17: Add E2E parity coverage for sandbox lifecycle and manager flows

**Files:**
- Create: `sdks/sandbox/go/tests/e2e/sandbox_e2e_test.go`
- Create: `sdks/sandbox/go/tests/e2e/sandbox_manager_e2e_test.go`
- Create: `sdks/sandbox/go/tests/e2e/code_interpreter_e2e_test.go`

- [ ] **Step 1: Write the first failing E2E test for lifecycle parity**

```go
func TestSandboxLifecycleHealthEndpointMetricsRenewConnect(t *testing.T) {}
```

- [ ] **Step 2: Run the E2E test to verify it fails**

Run: `cd sdks/sandbox/go && go test ./tests/e2e -run TestSandboxLifecycleHealthEndpointMetricsRenewConnect -v`
Expected: FAIL until setup helpers and implementation are complete

- [ ] **Step 3: Port core scenarios from existing JS/Python E2E suites**

```go
// manual cleanup
// network policy
// server proxy
// filesystem
// streaming
// manager flows
```

- [ ] **Step 4: Run full E2E suite**

Run: `cd sdks/sandbox/go && go test ./tests/e2e -v`
Expected: PASS in an environment with OpenSandbox test prerequisites configured

- [ ] **Step 5: Commit**

```bash
git add sdks/sandbox/go
git commit -m "test(go-sdk): add e2e parity coverage"
```

### Task 18: Final verification and repository integration checks

**Files:**
- Modify: `sdks/sandbox/go/Makefile`
- Modify: `sdks/sandbox/go/README.md`

- [ ] **Step 1: Run unit tests**

Run: `cd sdks/sandbox/go && go test ./tests/unit/...`
Expected: PASS

- [ ] **Step 2: Run integration tests**

Run: `cd sdks/sandbox/go && go test ./tests/integration/...`
Expected: PASS

- [ ] **Step 3: Run all package tests**

Run: `cd sdks/sandbox/go && go test ./...`
Expected: PASS or E2E skipped when env is missing

- [ ] **Step 4: Re-run generation to confirm clean working tree**

Run: `cd sdks/sandbox/go && make generate && git diff --exit-code`
Expected: no diff after regeneration

- [ ] **Step 5: Commit final polish**

```bash
git add sdks/sandbox/go
git commit -m "chore(go-sdk): finalize verification and developer workflow"
```

## Implementation Notes

- Prefer `oapi-codegen` or another Go-native generator that emits standard `net/http` compatible clients; do not handwrite the base transport layer.
- Keep generated code out of the public package namespace. Public consumers should work through handwritten models and high-level entry points.
- Use `context.Context` on every network-facing public method.
- Keep `Sandbox.Close()` scoped to local SDK resources only. Remote teardown must remain `Kill(ctx)`.
- `Resume(ctx, ...)` should return a fresh `*Sandbox`, matching the current JavaScript and C# semantics.
- `Timeout == nil` in `SandboxCreateOptions` should map to manual cleanup mode.
- Preserve `UseServerProxy` behavior when resolving execd and egress endpoints.
- `RunStream` should expose ordered events and a terminal error channel; `Run` should internally aggregate the same stream into a final `CommandExecution`.
- Minimize adapter logic duplication by centralizing wire-to-domain conversion in `sandbox/internal/convert`.

## Testing Notes

- Unit tests should use `httptest.Server` and narrowly validate request paths, headers, query params, and model conversion.
- SSE parser tests should cover:
  - multiple `data:` lines
  - blank line event termination
  - non-200 HTTP responses
  - malformed JSON payloads
  - context cancellation
- E2E tests should mirror behavior already asserted in:
  - `tests/javascript/tests/test_sandbox_e2e.test.ts`
  - `tests/python/tests/test_sandbox_e2e.py`
  - `tests/javascript/tests/test_sandbox_manager_e2e.test.ts`
  - `tests/python/tests/test_sandbox_manager_e2e.py`

## Risks and Mitigations

- **Risk:** Generated Go clients may not match the ergonomics of the existing SDKs.
  - **Mitigation:** Keep generated code internal and enforce adapter boundaries early.
- **Risk:** SSE handling may diverge from execd behavior.
  - **Mitigation:** Build parser tests from real event fixtures and aggregate through one execution path only.
- **Risk:** Manual cleanup and server-proxy behaviors are easy to regress.
  - **Mitigation:** Add dedicated unit and E2E tests before polishing examples.
- **Risk:** E2E setup may become expensive or flaky.
  - **Mitigation:** Keep unit coverage broad and mark E2E tests to skip when credentials or runtime prerequisites are absent.

## Definition of Done

- `sdks/sandbox/go` builds cleanly with generated and handwritten code separated.
- Public API includes `ConnectionConfig`, `Sandbox`, `SandboxManager`, stable models, and public errors.
- Commands, filesystem, health, metrics, lifecycle, and egress flows all work through high-level SDK entry points.
- Readiness and SSE command streaming semantics match existing SDK behavior closely enough for documentation and examples to be parallel.
- Unit and integration tests pass locally.
- E2E tests pass in a configured OpenSandbox environment or skip cleanly with explicit prerequisite messaging.
