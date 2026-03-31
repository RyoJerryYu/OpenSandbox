package adapters

import (
	"errors"
	"net/http"

	sandboxerrors "github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/errors"
)

// NormalizeHTTPError converts HTTP failures into the public sandbox error shape.
func NormalizeHTTPError(err error, resp *http.Response) error {
	if err == nil {
		return nil
	}

	var sandboxErr *sandboxerrors.SandboxError
	if errors.As(err, &sandboxErr) {
		if resp != nil {
			if sandboxErr.StatusCode == 0 {
				sandboxErr.StatusCode = resp.StatusCode
			}
			if sandboxErr.RequestID == "" {
				sandboxErr.RequestID = resp.Header.Get("X-Request-Id")
			}
		}
		return sandboxErr
	}

	if resp == nil {
		return err
	}

	return &sandboxerrors.SandboxError{
		Message:    err.Error(),
		RequestID:  resp.Header.Get("X-Request-Id"),
		StatusCode: resp.StatusCode,
		Cause:      err,
	}
}
