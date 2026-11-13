package managedplugin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/distribution/reference"
)

func testNamedTagged(t *testing.T, image string) reference.NamedTagged {
	t.Helper()
	ref, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		t.Fatal(err)
	}
	tagged, ok := ref.(reference.NamedTagged)
	if !ok {
		t.Fatalf("%s is not tagged", image)
	}
	return tagged
}

func TestFetchRegistryBearerToken(t *testing.T) {
	var authRequest *http.Request
	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authRequest = r.Clone(context.Background())
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"token":"test-bearer-token"}`))
	}))
	defer authServer.Close()

	registryServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("WWW-Authenticate", `Bearer realm="`+authServer.URL+`/registry/auth",service="cq-docker-registry"`)
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer registryServer.Close()

	ref := testNamedTagged(t, "docker.cloudquery.io/cloudquery/source-typeform:v1.2.3")
	token, err := fetchRegistryBearerToken(context.Background(), registryServer.Client(), registryServer.URL, ref, "test-api-key", "test-team")
	if err != nil {
		t.Fatal(err)
	}
	if token != "test-bearer-token" {
		t.Fatalf("got token %q, want %q", token, "test-bearer-token")
	}

	if authRequest == nil {
		t.Fatal("auth server did not receive a request")
	}
	username, password, ok := authRequest.BasicAuth()
	if !ok {
		t.Fatal("auth request had no basic auth")
	}
	if username != dockerRegistryUsername || password != "test-api-key" {
		t.Fatalf("got basic auth %s/%s", username, password)
	}
	if got := authRequest.Header.Get(teamNameHeader); got != "test-team" {
		t.Fatalf("got %s %q, want %q", teamNameHeader, got, "test-team")
	}
	if got := authRequest.Header.Get(pluginVersionHeader); got != "v1.2.3" {
		t.Fatalf("got %s %q, want %q", pluginVersionHeader, got, "v1.2.3")
	}
	if got := authRequest.URL.Query().Get("scope"); got != "repository:cloudquery/source-typeform:pull" {
		t.Fatalf("got scope %q", got)
	}
	if got := authRequest.URL.Query().Get("service"); got != "cq-docker-registry" {
		t.Fatalf("got service %q", got)
	}
	if got := authRequest.URL.Query().Get("account"); got != dockerRegistryUsername {
		t.Fatalf("got account %q", got)
	}
}

func TestFetchRegistryBearerTokenErrors(t *testing.T) {
	t.Run("no bearer challenge", func(t *testing.T) {
		registryServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer registryServer.Close()

		ref := testNamedTagged(t, "docker.cloudquery.io/cloudquery/source-typeform:v1.2.3")
		_, err := fetchRegistryBearerToken(context.Background(), registryServer.Client(), registryServer.URL, ref, "test-api-key", "test-team")
		if err == nil || !strings.Contains(err.Error(), "bearer auth challenge") {
			t.Fatalf("got error %v", err)
		}
	})

	t.Run("auth endpoint rejects credentials", func(t *testing.T) {
		authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer authServer.Close()

		registryServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("WWW-Authenticate", `Bearer realm="`+authServer.URL+`/registry/auth",service="cq-docker-registry"`)
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer registryServer.Close()

		ref := testNamedTagged(t, "docker.cloudquery.io/cloudquery/source-typeform:v1.2.3")
		_, err := fetchRegistryBearerToken(context.Background(), registryServer.Client(), registryServer.URL, ref, "test-api-key", "test-team")
		if err == nil || !strings.Contains(err.Error(), "401") {
			t.Fatalf("got error %v", err)
		}
	})
}
