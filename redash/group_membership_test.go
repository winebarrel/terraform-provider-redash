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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	redashgo "github.com/winebarrel/redash-go/v2"
)

func TestCreateGroupMemberSkipsAnExistingMember(t *testing.T) {
	store := newGroupStore(1, []int{10, 10})
	d := schema.TestResourceDataRaw(t, resourceGroupMember().Schema, map[string]any{
		"group_id": 10,
		"user_id":  1,
	})

	diags := createGroupUser(context.Background(), d, &providerMeta{client: store.client(t)})
	require.False(t, diags.HasError(), fmt.Sprint(diags))
	require.Equal(t, "10/1", d.Id())
	require.Equal(t, []int{10, 10}, store.groups(1))
	require.Equal(t, 0, store.posts)
}

func TestCreateGroupMemberAddsAMissingMember(t *testing.T) {
	store := newGroupStore(1, []int{1})
	d := schema.TestResourceDataRaw(t, resourceGroupMember().Schema, map[string]any{
		"group_id": 10,
		"user_id":  1,
	})

	diags := createGroupUser(context.Background(), d, &providerMeta{client: store.client(t)})
	require.False(t, diags.HasError(), fmt.Sprint(diags))
	require.Equal(t, "10/1", d.Id())
	require.Equal(t, []int{1, 10}, store.groups(1))
	require.Equal(t, 1, store.posts)
}

func TestRemoveGroupMemberDeletesEveryCopy(t *testing.T) {
	store := newGroupStore(1, []int{1, 10, 10})

	require.NoError(t, removeGroupMember(context.Background(), store.client(t), 1, 10))
	require.Equal(t, []int{1}, store.groups(1))
	require.Equal(t, 2, store.deletes)
}

func TestRemoveGroupMemberDeletesMoreThanFiveCopies(t *testing.T) {
	store := newGroupStore(1, []int{10, 10, 10, 10, 10, 10})

	require.NoError(t, removeGroupMember(context.Background(), store.client(t), 1, 10))
	require.Empty(t, store.groups(1))
	require.Equal(t, 6, store.deletes)
}

func TestRemoveGroupMemberErrorsWhenACopyRemains(t *testing.T) {
	store := newGroupStore(1, []int{10})
	store.keep = true

	err := removeGroupMember(context.Background(), store.client(t), 1, 10)
	require.EqualError(t, err, "user 1: group 10 still has 1 copy")
	require.Equal(t, []int{10}, store.groups(1))
	require.Equal(t, 1, store.deletes)
}

func TestReadGroupMemberDropsAMissingMember(t *testing.T) {
	store := newGroupStore(1, []int{1})
	d := schema.TestResourceDataRaw(t, resourceGroupMember().Schema, map[string]any{
		"group_id": 10,
		"user_id":  1,
	})
	d.SetId("10/1")

	diags := readGroupMember(context.Background(), d, &providerMeta{client: store.client(t)})
	require.False(t, diags.HasError(), fmt.Sprint(diags))
	require.Empty(t, d.Id())
}

type groupStore struct {
	mu      sync.Mutex
	users   map[int][]int
	posts   int
	deletes int
	keep    bool
}

func newGroupStore(userID int, groups []int) *groupStore {
	return &groupStore{users: map[int][]int{userID: append([]int(nil), groups...)}}
}

func (s *groupStore) groups(userID int) []int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return append([]int(nil), s.users[userID]...)
}

func (s *groupStore) memberList(groupID int) []byte {
	s.mu.Lock()
	defer s.mu.Unlock()

	var members []json.RawMessage
	for userID, groups := range s.users {
		if slices.Contains(groups, groupID) {
			members = append(members, userBody(userID, groups))
		}
	}

	body, err := json.Marshal(members)
	if err != nil {
		panic(err)
	}

	return body
}

func (s *groupStore) client(t *testing.T) *redashgo.Client {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/users/"):
			userID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/users/"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			_, _ = w.Write(userBody(userID, s.groups(userID)))
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/members"):
			groupID, err := groupIDFromPath(r.URL.Path)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			_, _ = w.Write(s.memberList(groupID))
		case r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/members"):
			groupID, err := groupIDFromPath(r.URL.Path)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			var body struct {
				UserID int `json:"user_id"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			s.mu.Lock()
			s.users[body.UserID] = append(s.users[body.UserID], groupID)
			s.posts++
			groups := append([]int(nil), s.users[body.UserID]...)
			s.mu.Unlock()
			_, _ = w.Write(userBody(body.UserID, groups))
		case r.Method == http.MethodDelete && strings.Contains(r.URL.Path, "/members/"):
			groupID, userID, err := memberIDsFromPath(r.URL.Path)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			s.mu.Lock()
			groups := s.users[userID]
			removed := false
			if !s.keep {
				for i, id := range groups {
					if id == groupID {
						s.users[userID] = append(groups[:i], groups[i+1:]...)
						removed = true
						break
					}
				}
			}
			s.deletes++
			s.mu.Unlock()
			if !s.keep && !removed {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)

	client, err := redashgo.NewClient(srv.URL, "test")
	require.NoError(t, err)

	return client
}

func groupIDFromPath(path string) (int, error) {
	path = strings.TrimPrefix(path, "/api/groups/")
	path = strings.TrimSuffix(path, "/members")

	return strconv.Atoi(path)
}

func memberIDsFromPath(path string) (int, int, error) {
	path = strings.TrimPrefix(path, "/api/groups/")
	group, user, ok := strings.Cut(path, "/members/")
	if !ok {
		return 0, 0, fmt.Errorf("bad path %s", path)
	}

	groupID, err := strconv.Atoi(group)
	if err != nil {
		return 0, 0, err
	}

	userID, err := strconv.Atoi(user)
	if err != nil {
		return 0, 0, err
	}

	return groupID, userID, nil
}

func userBody(id int, groups []int) []byte {
	body, err := json.Marshal(map[string]any{
		"id":                    id,
		"name":                  "member",
		"email":                 "member@example.com",
		"groups":                groups,
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
