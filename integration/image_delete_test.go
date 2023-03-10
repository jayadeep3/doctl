package integration

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os/exec"
	"strings"
	"testing"

	"github.com/sclevine/spec"
	"github.com/stretchr/testify/require"
)

var _ = suite("compute/image/delete", func(t *testing.T, when spec.G, it spec.S) {
	var (
		expect *require.Assertions
		server *httptest.Server
	)

	it.Before(func() {
		expect = require.New(t)
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			switch req.URL.Path {
			case "/v2/images/1111":
				auth := req.Header.Get("Authorization")
				if auth != "Bearer some-magic-token" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				if req.Method != http.MethodDelete {
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}

				w.WriteHeader(http.StatusNoContent)
			case "/v2/images/2222":
				auth := req.Header.Get("Authorization")
				if auth != "Bearer some-magic-token" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				if req.Method != http.MethodDelete {
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}

				w.WriteHeader(http.StatusNoContent)
			default:
				dump, err := httputil.DumpRequest(req, true)
				if err != nil {
					t.Fatal("failed to dump request")
				}

				t.Fatalf("received unknown request: %s", dump)
			}
		}))
	})

	when("all the required flags are passed", func() {
		it("deletes the image", func() {
			cmd := exec.Command(builtBinaryPath,
				"-t", "some-magic-token",
				"-u", server.URL,
				"compute",
				"image",
				"delete",
				"1111",
				"--force",
			)

			output, err := cmd.CombinedOutput()
			expect.NoError(err, fmt.Sprintf("received unexpected error: %s", output))
			expect.Empty(output)
		})
	})

	when("multiple images are provided and all the required flags are passed", func() {
		it("deletes the images", func() {
			cmd := exec.Command(builtBinaryPath,
				"-t", "some-magic-token",
				"-u", server.URL,
				"compute",
				"image",
				"delete",
				"1111",
				"2222",
				"--force",
			)

			output, err := cmd.CombinedOutput()
			expect.NoError(err, fmt.Sprintf("received unexpected error: %s", output))
			expect.Empty(output)
		})
	})

	when("deleting one image without force flag", func() {
		it("errors without confirmation", func() {
			cmd := exec.Command(builtBinaryPath,
				"-t", "some-magic-token",
				"-u", server.URL,
				"compute",
				"image",
				"delete",
				"1111",
			)

			output, err := cmd.CombinedOutput()
			expect.Error(err)
			expect.Equal(strings.TrimSpace(confirmNonInteractiveOutput), strings.TrimSpace(string(output)))
		})
	})

	when("deleting two images without force flag", func() {
		it("errors without confirmation", func() {
			cmd := exec.Command(builtBinaryPath,
				"-t", "some-magic-token",
				"-u", server.URL,
				"compute",
				"image",
				"delete",
				"1111",
				"2222",
			)

			output, err := cmd.CombinedOutput()
			expect.Error(err)
			expect.Equal(strings.TrimSpace(confirmNonInteractiveOutput), strings.TrimSpace(string(output)))
		})
	})
})
