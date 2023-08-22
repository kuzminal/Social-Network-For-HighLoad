package pg

import (
	"SocialNetHL/internal/store"
	"context"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
	"time"
)

var repo store.FriendStore

func init() {
	ctx := context.Background()
	container, _ := postgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:15.2-alpine"),
		postgres.WithDatabase("postgres"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	connString, _ := container.ConnectionString(ctx, "sslmode=disable")
	repo, _ = NewMaster(ctx, connString)
}

func TestPostgres_AddFriendAdd(t *testing.T) {
	ctx := context.Background()
	require.NotNil(t, repo)
	err := repo.AddFriend(ctx, "1", "123")
	require.NoError(t, err)
}

func TestPostgres_FindFriends(t *testing.T) {
	ctx := context.Background()
	require.NotNil(t, repo)
	friends, err := repo.FindFriends(ctx, "123")
	require.NoError(t, err)
	require.Contains(t, friends, "1")
}
