package awx

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	awx "github.com/denouche/goawx/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// newTestAWX starts a fake AWX that answers GET /api/v2/jobs/42/ with status
// and /api/v2/ping/ (which NewAWX calls) with 200.
func newTestAWX(t *testing.T, status int) *awx.AWX {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v2/ping/":
			fmt.Fprint(w, `{"version":"24.6.1"}`)
		case "/api/v2/jobs/42/":
			w.WriteHeader(status)
			if status == http.StatusOK {
				fmt.Fprint(w, `{"id":42,"status":"successful"}`)
			} else {
				fmt.Fprint(w, `{"detail":"Not found."}`)
			}
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)

	client, err := awx.NewAWX(srv.URL, "user", "password", srv.Client())
	if err != nil {
		t.Fatalf("NewAWX: %v", err)
	}
	return client
}

func testLaunchResourceData(t *testing.T) *schema.ResourceData {
	t.Helper()
	d := schema.TestResourceDataRaw(t, resourceJobTemplateLaunch().Schema, map[string]interface{}{
		"job_template_id": 7,
	})
	d.SetId("42")
	return d
}

func TestResourceJobDelete(t *testing.T) {
	cases := []struct {
		name      string
		status    int
		wantError bool
		wantID    string
	}{
		{"job exists", http.StatusOK, false, ""},
		{"job purged by AWX cleanup", http.StatusNotFound, false, ""},
		{"AWX error", http.StatusInternalServerError, true, "42"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := testLaunchResourceData(t)
			diags := resourceJobDelete(context.Background(), d, newTestAWX(t, tc.status))

			if diags.HasError() != tc.wantError {
				t.Fatalf("HasError() = %v, want %v: %v", diags.HasError(), tc.wantError, diags)
			}
			if d.Id() != tc.wantID {
				t.Errorf("Id() = %q, want %q", d.Id(), tc.wantID)
			}
		})
	}
}

func TestIsNotFoundError(t *testing.T) {
	cases := []struct {
		err  error
		want bool
	}{
		{nil, false},
		{fmt.Errorf("responsed with 404, resp: &{...}"), true},
		{fmt.Errorf("responsed with 500, resp: &{...}"), false},
		{fmt.Errorf("Errors:\n- detail: [404]"), false},
	}
	for _, tc := range cases {
		if got := isNotFoundError(tc.err); got != tc.want {
			t.Errorf("isNotFoundError(%v) = %v, want %v", tc.err, got, tc.want)
		}
	}
}
