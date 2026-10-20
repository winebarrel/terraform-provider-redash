package redash

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	redash_go "github.com/winebarrel/redash-go/v2"
)

func TestRemoveFromDefaultGroupOrDelete(t *testing.T) {
	tests := []struct {
		name         string
		deleteStatus int
		wantID       string
		wantErr      []string
	}{
		{
			name:         "data source is deleted",
			deleteStatus: http.StatusOK,
			wantID:       "",
			wantErr:      []string{"remove data source from default group", "was deleted"},
		},
		{
			name:         "cleanup fails too",
			deleteStatus: http.StatusInternalServerError,
			wantID:       "7",
			wantErr:      []string{"remove data source 7 from default group", "deleting it also failed"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				mu    sync.Mutex
				calls []string
			)

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mu.Lock()
				calls = append(calls, r.Method+" "+r.URL.Path)
				mu.Unlock()

				switch r.Method + " " + r.URL.Path {
				case "GET /api/groups":
					_, _ = w.Write([]byte(`[{"id":1,"name":"default","type":"builtin"}]`))
				case "DELETE /api/groups/1/data_sources/7":
					w.WriteHeader(http.StatusInternalServerError)
				case "DELETE /api/data_sources/7":
					w.WriteHeader(tt.deleteStatus)
				default:
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			t.Cleanup(srv.Close)

			client, err := redash_go.NewClient(srv.URL, "test-key")
			require.NoError(t, err)

			d := schema.TestResourceDataRaw(t, resourceDataSource().Schema, map[string]any{
				"name": "ds",
				"type": "pg",
			})
			d.SetId("7")

			diags := removeFromDefaultGroupOrDelete(context.Background(), d, client, 7)

			require.True(t, diags.HasError())
			require.Equal(t, tt.wantID, d.Id())

			for _, s := range tt.wantErr {
				require.Contains(t, diags[0].Summary, s)
			}

			require.Contains(t, calls, "DELETE /api/data_sources/7")
		})
	}
}
