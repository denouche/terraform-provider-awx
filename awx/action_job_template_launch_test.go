package awx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	awx "github.com/denouche/goawx/client"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

// TestMuxProviderSchemas fails if the SDKv2 and framework provider schemas
// drift apart, which tf5muxserver rejects at runtime.
func TestMuxProviderSchemas(t *testing.T) {
	ctx := context.Background()
	mux, err := tf5muxserver.NewMuxServer(ctx,
		Provider().GRPCProvider,
		providerserver.NewProtocol5(NewFrameworkProvider()),
	)
	if err != nil {
		t.Fatalf("NewMuxServer: %v", err)
	}

	resp, err := mux.ProviderServer().GetProviderSchema(ctx, &tfprotov5.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatalf("GetProviderSchema: %v", err)
	}
	for _, d := range resp.Diagnostics {
		if d.Severity == tfprotov5.DiagnosticSeverityError {
			t.Fatalf("GetProviderSchema: %s: %s", d.Summary, d.Detail)
		}
	}
	if _, ok := resp.ActionSchemas["awx_job_template_launch"]; !ok {
		t.Errorf("action awx_job_template_launch is not served by the mux server")
	}
	if _, ok := resp.ResourceSchemas["awx_job_template_launch"]; !ok {
		t.Errorf("resource awx_job_template_launch is not served by the mux server")
	}
}

// fakeAWX serves the endpoints the action calls. Job 43 walks through
// jobStatuses, one per poll, and stays on the last one.
type fakeAWX struct {
	mu          sync.Mutex
	jobStatuses []string
	polls       int
	launched    map[string]interface{}
}

func (f *fakeAWX) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	f.mu.Lock()
	defer f.mu.Unlock()

	switch {
	case r.URL.Path == "/api/v2/ping/":
		fmt.Fprint(w, `{"version":"24.6.1"}`)
	case r.URL.Path == "/api/v2/job_templates/7/" && r.Method == http.MethodGet:
		fmt.Fprint(w, `{"id":7,"name":"upgrade"}`)
	case r.URL.Path == "/api/v2/job_templates/7/launch/" && r.Method == http.MethodPost:
		json.NewDecoder(r.Body).Decode(&f.launched)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"id":43,"job":43}`)
	case r.URL.Path == "/api/v2/jobs/43/":
		i := f.polls
		if i >= len(f.jobStatuses) {
			i = len(f.jobStatuses) - 1
		}
		f.polls++
		fmt.Fprintf(w, `{"id":43,"status":%q}`, f.jobStatuses[i])
	default:
		http.NotFound(w, r)
	}
}

func newFakeAWXClient(t *testing.T, f *fakeAWX) *awx.AWX {
	t.Helper()
	srv := httptest.NewServer(f)
	t.Cleanup(srv.Close)
	client, err := awx.NewAWX(srv.URL, "user", "password", srv.Client())
	if err != nil {
		t.Fatalf("NewAWX: %v", err)
	}
	return client
}

func fastJobPolling(t *testing.T) {
	t.Helper()
	delay, minTimeout := jobPollDelay, jobPollMinTimeout
	jobPollDelay, jobPollMinTimeout = 10*time.Millisecond, 10*time.Millisecond
	t.Cleanup(func() { jobPollDelay, jobPollMinTimeout = delay, minTimeout })
}

func invokeLaunchAction(t *testing.T, client *awx.AWX, attrs map[string]tftypes.Value) (*action.InvokeResponse, []string) {
	t.Helper()
	ctx := context.Background()
	a := NewJobTemplateLaunchAction().(*jobTemplateLaunchAction)
	a.client = client

	var schemaResp action.SchemaResponse
	a.Schema(ctx, action.SchemaRequest{}, &schemaResp)
	objType := schemaResp.Schema.Type().TerraformType(ctx).(tftypes.Object)

	values := map[string]tftypes.Value{}
	for name, typ := range objType.AttributeTypes {
		values[name] = tftypes.NewValue(typ, nil)
	}
	for name, v := range attrs {
		values[name] = v
	}

	var progress []string
	resp := &action.InvokeResponse{
		SendProgress: func(e action.InvokeProgressEvent) { progress = append(progress, e.Message) },
	}
	a.Invoke(ctx, action.InvokeRequest{
		Config: tfsdk.Config{Schema: schemaResp.Schema, Raw: tftypes.NewValue(objType, values)},
	}, resp)
	return resp, progress
}

func TestJobTemplateLaunchAction(t *testing.T) {
	fastJobPolling(t)

	cases := []struct {
		name         string
		statuses     []string
		attrs        map[string]tftypes.Value
		wantError    string
		wantProgress []string
	}{
		{
			name:     "waits for a successful job",
			statuses: []string{"pending", "running", "successful"},
			attrs: map[string]tftypes.Value{
				"extra_vars": tftypes.NewValue(tftypes.String, `{"k8s_version":"1.31.9"}`),
			},
			wantProgress: []string{
				"Launched job 43 from job template 7",
				"Job 43 is pending",
				"Job 43 is running",
				"Job 43 is successful",
			},
		},
		{
			name:      "fails when the job fails",
			statuses:  []string{"running", "failed"},
			wantError: `finished with status "failed"`,
		},
		{
			name:     "returns straight after launch without waiting",
			statuses: []string{"running"},
			attrs: map[string]tftypes.Value{
				"wait_for_completion": tftypes.NewValue(tftypes.Bool, false),
			},
			wantProgress: []string{"Launched job 43 from job template 7"},
		},
		{
			name:     "rejects a malformed timeout",
			statuses: []string{"successful"},
			attrs: map[string]tftypes.Value{
				"timeout": tftypes.NewValue(tftypes.String, "ninety minutes"),
			},
			wantError: "Invalid timeout",
		},
		{
			name:     "times out a job that never finishes",
			statuses: []string{"running"},
			attrs: map[string]tftypes.Value{
				"timeout": tftypes.NewValue(tftypes.String, "200ms"),
			},
			wantError: "did not succeed",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := &fakeAWX{jobStatuses: tc.statuses}
			attrs := map[string]tftypes.Value{
				"job_template_id": tftypes.NewValue(tftypes.Number, 7),
			}
			for k, v := range tc.attrs {
				attrs[k] = v
			}

			resp, progress := invokeLaunchAction(t, newFakeAWXClient(t, f), attrs)

			if tc.wantError == "" && resp.Diagnostics.HasError() {
				t.Fatalf("unexpected error: %v", resp.Diagnostics)
			}
			if tc.wantError != "" {
				if !resp.Diagnostics.HasError() {
					t.Fatalf("expected an error containing %q", tc.wantError)
				}
				if got := fmt.Sprint(resp.Diagnostics); !strings.Contains(got, tc.wantError) {
					t.Fatalf("error %q does not contain %q", got, tc.wantError)
				}
			}
			if tc.wantProgress != nil && strings.Join(progress, "|") != strings.Join(tc.wantProgress, "|") {
				t.Errorf("progress = %q, want %q", progress, tc.wantProgress)
			}
			if v, ok := tc.attrs["extra_vars"]; ok {
				var want string
				v.As(&want)
				if got := f.launched["extra_vars"]; got != want {
					t.Errorf("launched extra_vars = %v, want %q", got, want)
				}
			}
		})
	}
}

func TestJobTemplateLaunchActionUnconfigured(t *testing.T) {
	resp, _ := invokeLaunchAction(t, nil, map[string]tftypes.Value{
		"job_template_id": tftypes.NewValue(tftypes.Number, 7),
	})
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error when the provider has no client")
	}
}
