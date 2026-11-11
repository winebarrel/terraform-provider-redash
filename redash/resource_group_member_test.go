package redash

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	redashgo "github.com/winebarrel/redash-go/v2"
)

func TestGroupMemberParallelCreateKeepsEveryGroup(t *testing.T) {
	const userID = 7
	var (
		mu       sync.Mutex
		groups   []int
		inFlight atomic.Int32
		maxSeen  atomic.Int32
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodGet {
			groupID, err := strconv.Atoi(strings.TrimPrefix(strings.TrimSuffix(r.URL.Path, "/members"), "/api/groups/"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			mu.Lock()
			present := slices.Contains(groups, groupID)
			mu.Unlock()
			if !present {
				_, _ = w.Write([]byte("[]"))
				return
			}
			_, _ = w.Write([]byte("[" + string(memberUserJSON(userID)) + "]"))
			return
		}

		groupID, err := strconv.Atoi(strings.TrimPrefix(strings.TrimSuffix(r.URL.Path, "/members"), "/api/groups/"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Same shape as Redash: read the array, pause, write the whole array back.
		n := inFlight.Add(1)
		for {
			old := maxSeen.Load()
			if n <= old || maxSeen.CompareAndSwap(old, n) {
				break
			}
		}
		mu.Lock()
		snapshot := append([]int(nil), groups...)
		mu.Unlock()
		time.Sleep(30 * time.Millisecond)
		inFlight.Add(-1)

		snapshot = append(snapshot, groupID)
		mu.Lock()
		groups = snapshot
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(memberUserJSON(userID))
	}))
	t.Cleanup(srv.Close)

	client, err := redashgo.NewClient(srv.URL, "test")
	require.NoError(t, err)

	res := resourceGroupMember()
	errCh := make(chan error, 3)
	var wg sync.WaitGroup
	for _, groupID := range []int{10, 11, 12} {
		wg.Add(1)
		go func(groupID int) {
			defer wg.Done()
			d := res.Data(nil)
			if err := d.Set("group_id", groupID); err != nil {
				errCh <- err
				return
			}
			if err := d.Set("user_id", userID); err != nil {
				errCh <- err
				return
			}
			if diags := createGroupUser(context.Background(), d, &providerMeta{client: client}); diags.HasError() {
				errCh <- fmt.Errorf("group %d: %v", groupID, diags)
			}
		}(groupID)
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		require.NoError(t, err)
	}

	require.Equal(t, int32(1), maxSeen.Load(), "adds for one user must not run together")
	require.ElementsMatch(t, []int{10, 11, 12}, groups)
}

func TestDeleteGroupUserAbsentMembershipSucceeds(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(memberUserJSON(7))
	}))
	t.Cleanup(srv.Close)

	client, err := redashgo.NewClient(srv.URL, "test")
	require.NoError(t, err)

	d := resourceGroupMember().Data(nil)
	d.SetId("10/7")
	diags := deleteGroupUser(context.Background(), d, &providerMeta{client: client})
	require.False(t, diags.HasError(), fmt.Sprint(diags))
	require.Empty(t, d.Id())
}

func memberUserJSON(id int) []byte {
	body, err := json.Marshal(map[string]any{
		"id":                    id,
		"name":                  "member",
		"email":                 "member@example.com",
		"groups":                []int{},
		"auth_type":             "external",
		"created_at":            "2020-01-01T00:00:00Z",
		"updated_at":            "2020-01-01T00:00:00Z",
		"active_at":             "2020-01-01T00:00:00Z",
		"disabled_at":           "0001-01-01T00:00:00Z",
		"is_disabled":           false,
		"is_email_verified":     true,
		"is_invitation_pending": false,
	})
	if err != nil {
		panic(err)
	}
	return body
}
