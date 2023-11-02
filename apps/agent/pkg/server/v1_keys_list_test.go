package server_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	authenticationv1 "github.com/unkeyed/unkey/apps/agent/gen/proto/authentication/v1"
	"github.com/unkeyed/unkey/apps/agent/pkg/hash"
	"github.com/unkeyed/unkey/apps/agent/pkg/server"
	"github.com/unkeyed/unkey/apps/agent/pkg/testutil"
	"github.com/unkeyed/unkey/apps/agent/pkg/uid"
	"github.com/unkeyed/unkey/apps/agent/pkg/util"
)

func TestV1ListKeys_Simple(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	resources := testutil.SetupResources(t)

	srv := testutil.NewServer(t, resources)

	createdKeyIds := make([]string, 10)
	for i := range createdKeyIds {
		key := &authenticationv1.Key{
			KeyId:       uid.Key(),
			KeyAuthId:   resources.UserKeyAuth.KeyAuthId,
			Name:        util.Pointer(fmt.Sprintf("test-%d", i)),
			WorkspaceId: resources.UserWorkspace.WorkspaceId,
			Hash:        hash.Sha256(uid.New(16, "test")),
			CreatedAt:   time.Now().UnixMilli(),
		}
		err := resources.Database.InsertKey(ctx, key)
		require.NoError(t, err)
		createdKeyIds[i] = key.KeyId

	}

	successResponse := testutil.Json[server.ListKeysResponse](t, srv.App, testutil.JsonRequest{

		Path:       fmt.Sprintf("/v1/apis.listKeys?apiId=%s", resources.UserApi.ApiId),
		Bearer:     resources.UserRootKey,
		StatusCode: 200,
	})

	require.GreaterOrEqual(t, successResponse.Total, int64(len(createdKeyIds)))
	require.GreaterOrEqual(t, len(successResponse.Keys), len(createdKeyIds))
	require.LessOrEqual(t, len(successResponse.Keys), 100) //  default page size

}

func TestV1ListKeys_FilterOwnerId(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	resources := testutil.SetupResources(t)

	srv := testutil.NewServer(t, resources)

	createdKeyIds := make([]string, 10)
	for i := range createdKeyIds {
		key := &authenticationv1.Key{
			KeyId:       uid.Key(),
			KeyAuthId:   resources.UserKeyAuth.KeyAuthId,
			WorkspaceId: resources.UserWorkspace.WorkspaceId,
			Hash:        hash.Sha256(uid.New(16, "test")),
			CreatedAt:   time.Now().UnixMilli(),
		}
		// just add an ownerId to half of them
		if i%2 == 0 {
			key.OwnerId = util.Pointer("chronark")
		}
		err := resources.Database.InsertKey(ctx, key)
		require.NoError(t, err)
		createdKeyIds[i] = key.KeyId

	}

	successResponse := testutil.Json[server.ListKeysResponse](t, srv.App, testutil.JsonRequest{

		Path:       fmt.Sprintf("/v1/apis.listKeys?apiId=%s", resources.UserApi.ApiId),
		Bearer:     resources.UserRootKey,
		StatusCode: 200,
	})

	require.GreaterOrEqual(t, successResponse.Total, int64(len(createdKeyIds)))
	require.Equal(t, 5, len(successResponse.Keys))
	require.LessOrEqual(t, len(successResponse.Keys), 100) //  default page size

	for _, key := range successResponse.Keys {
		require.Equal(t, "chronark", key.OwnerId)
	}

}

func TestV1ListKeys_WithLimit(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	resources := testutil.SetupResources(t)
	srv := testutil.NewServer(t, resources)

	createdKeyIds := make([]string, 10)
	for i := range createdKeyIds {
		key := &authenticationv1.Key{
			KeyId:       uid.Key(),
			KeyAuthId:   resources.UserKeyAuth.KeyAuthId,
			WorkspaceId: resources.UserWorkspace.WorkspaceId,
			Hash:        hash.Sha256(uid.New(16, "test")),
			CreatedAt:   time.Now().UnixMilli(),
		}
		err := resources.Database.InsertKey(ctx, key)
		require.NoError(t, err)
		createdKeyIds[i] = key.KeyId

	}

	successResponse := testutil.Get[server.ListKeysResponse](t, srv.App, testutil.GetRequest{
		Path:       fmt.Sprintf("/v1/apis.listKeys?apiId=%s", resources.UserApi.ApiId),
		Bearer:     resources.UserRootKey,
		StatusCode: 200,
	})

	require.GreaterOrEqual(t, successResponse.Total, int64(len(createdKeyIds)))
	require.Equal(t, 2, len(successResponse.Keys))

}

func TestV1ListKeys_WithOffset(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	resources := testutil.SetupResources(t)

	srv := testutil.NewServer(t, resources)

	createdKeyIds := make([]string, 10)
	for i := range createdKeyIds {
		key := &authenticationv1.Key{
			KeyId:       uid.Key(),
			KeyAuthId:   resources.UserKeyAuth.KeyAuthId,
			WorkspaceId: resources.UserWorkspace.WorkspaceId,
			Hash:        hash.Sha256(uid.New(16, "test")),
			CreatedAt:   time.Now().UnixMilli(),
		}
		err := resources.Database.InsertKey(ctx, key)
		require.NoError(t, err)
		createdKeyIds[i] = key.KeyId

	}

	res1 := testutil.Get[server.ListKeysResponse](t, srv.App, testutil.GetRequest{
		Path:       fmt.Sprintf("/v1/apis.listKeys?apiId=%s", resources.UserApi.ApiId),
		Bearer:     resources.UserRootKey,
		StatusCode: 200,
	})

	require.GreaterOrEqual(t, res1.Total, int64(len(createdKeyIds)))
	require.GreaterOrEqual(t, 10, len(res1.Keys))

	res2 := testutil.Get[server.ListKeysResponse](t, srv.App, testutil.GetRequest{

		Path:       fmt.Sprintf("/v1/apis.listKeys?apiId=%s", resources.UserApi.ApiId),
		Bearer:     resources.UserRootKey,
		StatusCode: 200,
	})

	require.GreaterOrEqual(t, res2.Total, int64(len(createdKeyIds)))
	require.GreaterOrEqual(t, len(res2.Keys), len(createdKeyIds)-2)

	require.Equal(t, res1.Keys[1].Id, res2.Keys[0].Id)

}