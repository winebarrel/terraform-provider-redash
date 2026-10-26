package redash

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"net/http/httptest"
	"slices"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	redash_go "github.com/winebarrel/redash-go/v2"
)

// fakeGroupGrants serves the grants of group 21 the way Redash does: deleting
// a data source drops its grants, and removing the grant of a data source
// that no longer exists answers 500.
type fakeGroupGrants struct {
	mu          sync.Mutex
	dataSources map[int]struct{}
	grants      map[int]bool // data source id -> view_only
	removeFails bool
	// deleteAfterList is a data source deleted right after the grants are
	// listed for the first time, as when its destroy runs in the same apply.
	deleteAfterList int
}

func newFakeGroupGrants(t *testing.T) (*redash_go.Client, *fakeGroupGrants) {
	t.Helper()

	f := &fakeGroupGrants{
		dataSources: map[int]struct{}{4: {}, 5: {}},
		grants:      map[int]bool{4: false, 5: true},
	}
	srv := httptest.NewServer(http.HandlerFunc(f.serve))
	t.Cleanup(srv.Close)

	client, err := redash_go.NewClient(srv.URL, "test-key")
	require.NoError(t, err)

	return client, f
}

func (f *fakeGroupGrants) deleteDataSource(id int) {
	delete(f.dataSources, id)
	delete(f.grants, id)
}

func (f *fakeGroupGrants) serve(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	var dsId int
	n, _ := fmt.Sscanf(r.URL.Path, "/api/groups/21/data_sources/%d", &dsId)

	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/api/groups/21/data_sources":
		list := make([]map[string]any, 0, len(f.grants))
		for _, id := range slices.Sorted(maps.Keys(f.grants)) {
			list = append(list, map[string]any{"id": id, "name": fmt.Sprintf("ds%d", id), "view_only": f.grants[id]})
		}
		_ = json.NewEncoder(w).Encode(list)

		if f.deleteAfterList != 0 {
			f.deleteDataSource(f.deleteAfterList)
			f.deleteAfterList = 0
		}
	case r.Method == http.MethodDelete && n == 1:
		if _, ok := f.dataSources[dsId]; !ok || f.removeFails {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		delete(f.grants, dsId)
		_, _ = w.Write([]byte("null"))
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func TestRemoveGroupDataSource_dataSourceDeletedInSameApply(t *testing.T) {
	entryPoints := []struct {
		name string
		id   string
		run  func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics
		// listsFirst is true for the paths that list the grants before removing any.
		listsFirst bool
		// wantGrants is the server state once the removal of 4 counts as done.
		wantGrants map[int]bool
		wantId     string
	}{
		{
			name:       "redash_group_data_sources update",
			id:         "21",
			run:        updateGroupDataSources,
			listsFirst: true,
			wantGrants: map[int]bool{5: true},
			wantId:     "21",
		},
		{
			name:       "redash_group_data_sources delete",
			id:         "21",
			run:        deleteGroupDataSources,
			listsFirst: true,
			wantGrants: map[int]bool{},
			wantId:     "",
		},
		{
			name:       "redash_group_data_source delete",
			id:         "21/4",
			run:        deleteGroupDataSource,
			wantGrants: map[int]bool{5: true},
			wantId:     "",
		},
	}

	newData := func(t *testing.T, id string) *schema.ResourceData {
		var d *schema.ResourceData
		if id == "21/4" {
			d = schema.TestResourceDataRaw(t, resourceGroupDataSource().Schema, map[string]any{
				"group_id":       21,
				"data_source_id": 4,
			})
		} else {
			d = schema.TestResourceDataRaw(t, resourceGroupDataSources().Schema, map[string]any{
				"group_id":    21,
				"data_source": []any{map[string]any{"data_source_id": 5, "view_only": true}},
			})
		}
		d.SetId(id)
		return d
	}

	for _, ep := range entryPoints {
		t.Run(ep.name+"/data source 4 deleted meanwhile", func(t *testing.T) {
			client, f := newFakeGroupGrants(t)
			if ep.listsFirst {
				f.deleteAfterList = 4
			} else {
				f.deleteDataSource(4)
			}
			d := newData(t, ep.id)

			diags := ep.run(context.Background(), d, client)

			require.False(t, diags.HasError(), diags)
			require.Equal(t, ep.wantGrants, f.grants)
			require.Equal(t, ep.wantId, d.Id())
		})

		t.Run(ep.name+"/removal fails while the grant remains", func(t *testing.T) {
			client, f := newFakeGroupGrants(t)
			f.removeFails = true
			d := newData(t, ep.id)

			diags := ep.run(context.Background(), d, client)

			require.True(t, diags.HasError())
			require.Equal(t, map[int]bool{4: false, 5: true}, f.grants)
			require.Equal(t, ep.id, d.Id())
		})
	}
}
