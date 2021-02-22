package cluster

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/redsailtechnologies/boatswain/pkg/auth"
	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
	"github.com/redsailtechnologies/boatswain/rpc/agent"
	pb "github.com/redsailtechnologies/boatswain/rpc/cluster"
)

type mockAgentAction struct {
	mock.Mock
}

func (maa *mockAgentAction) Run(context.Context, *agent.Action) (*agent.Result, error) {
	return nil, nil
}

type mockAuth struct {
	mock.Mock
}

func (ma *mockAuth) Authenticate(ctx context.Context) (context.Context, error) {
	args := ma.Called(ctx)
	err := args.Get(1)
	if err != nil {
		return nil, err.(error)
	}
	return args.Get(0).(context.Context), nil
}

func (ma *mockAuth) Authorize(ctx context.Context, role auth.Role) error {
	args := ma.Called(ctx, role)
	err := args.Get(0)
	if err != nil {
		return err.(error)
	}
	return nil
}

func (ma *mockAuth) NewContext(ctx context.Context) (context.Context, error) {
	ma.Called(ctx)
	return context.Background(), nil
}

func (ma *mockAuth) User(ctx context.Context) auth.User {
	ma.Called(ctx)
	return auth.User{}
}

func (ma *mockAuth) Roles(u auth.User) []auth.Role {
	ma.Called(u)
	return []auth.Role{}
}

func (ma *mockAuth) Wrap(h http.Handler) http.Handler {
	ma.Called(h)
	return h
}

type mockStorage struct {
	mock.Mock
}

func (ms *mockStorage) CheckReady() error {
	args := ms.Called(collection)
	err := args.Get(0)
	if err != nil {
		return err.(error)
	}
	return nil
}

func (ms *mockStorage) IDs(collection string) ([]string, error) {
	args := ms.Called(collection)
	err := args.Get(1)
	if err != nil {
		return nil, err.(error)
	}
	return args.Get(0).([]string), nil
}

func (ms *mockStorage) GetEvents(collection, uuid string) ([]*storage.StoredEvent, error) {
	args := ms.Called(collection, uuid)
	err := args.Get(1)
	if err != nil {
		return nil, err.(error)
	}
	return args.Get(0).([]*storage.StoredEvent), nil
}

func (ms *mockStorage) GetVersion(collection, uuid string) int {
	args := ms.Called(collection, uuid)
	return args.Int(0)
}

func (ms *mockStorage) StoreEvent(collection, uuid, eventType, eventData string, version int) error {
	args := ms.Called(collection, uuid, eventType, eventData, version)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(error)
}

func TestNewService(t *testing.T) {
	assert.NotNil(t, NewService(&mockAgentAction{}, &mockAuth{}, &mockStorage{}))
}

func TestCreateAuth(t *testing.T) {
	a := &mockAuth{}
	a.On("Authorize", mock.Anything, mock.Anything).Return(auth.NotAuthorizedError{})
	store := &mockStorage{}
	sut := NewService(&mockAgentAction{}, a, store)

	res, err := sut.Create(context.TODO(), &pb.CreateCluster{
		Name: "name",
	})

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestCreateValid(t *testing.T) {
	auth := &mockAuth{}
	auth.On("Authorize", mock.Anything, mock.Anything).Return(nil)
	store := &mockStorage{}
	store.On("GetVersion", mock.Anything, mock.Anything).Return(0)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sut := NewService(&mockAgentAction{}, auth, store)

	res, err := sut.Create(context.TODO(), &pb.CreateCluster{
		Name: "name",
	})

	assert.NotEqual(t, "", res.Uuid)
	assert.Nil(t, err)
}

func TestCreateInvalid(t *testing.T) {
	auth := &mockAuth{}
	auth.On("Authorize", mock.Anything, mock.Anything).Return(nil)
	store := &mockStorage{}
	store.On("GetVersion", mock.Anything, mock.Anything).Return(0)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sut := NewService(&mockAgentAction{}, auth, store)

	res, err := sut.Create(context.TODO(), &pb.CreateCluster{
		Name: "",
	})

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestCreateSaveError(t *testing.T) {
	auth := &mockAuth{}
	auth.On("Authorize", mock.Anything, mock.Anything).Return(nil)
	store := &mockStorage{}
	store.On("GetVersion", mock.Anything, mock.Anything).Return(0)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New(""))
	sut := NewService(&mockAgentAction{}, auth, store)

	res, err := sut.Create(context.TODO(), &pb.CreateCluster{
		Name: "name",
	})

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestUpdateAuth(t *testing.T) {
	a := &mockAuth{}
	a.On("Authorize", mock.Anything, mock.Anything).Return(auth.NotAuthorizedError{})
	store := &mockStorage{}
	store.On("GetEvents", mock.Anything, "a").Return([]*storage.StoredEvent{
		{
			UUID:      "a",
			Version:   1,
			Timestamp: ddd.NewTimestamp(),
			Type:      "ClusterCreated",
			Data:      `{"Timestamp":0,"UUID":"a","Name":"name","Endpoint":"endpoint","Token":"token","Cert":"cert"}`,
		},
	}, nil)
	store.On("GetVersion", mock.Anything, mock.Anything).Return(1)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sut := NewService(&mockAgentAction{}, a, store)

	res, err := sut.Update(context.TODO(), &pb.UpdateCluster{
		Uuid: "a",
		Name: "NEWname",
	})

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestUpdateValid(t *testing.T) {
	auth := &mockAuth{}
	auth.On("Authorize", mock.Anything, mock.Anything).Return(nil)
	store := &mockStorage{}
	store.On("GetEvents", mock.Anything, "a").Return([]*storage.StoredEvent{
		{
			UUID:      "a",
			Version:   1,
			Timestamp: ddd.NewTimestamp(),
			Type:      "ClusterCreated",
			Data:      `{"Timestamp":0,"UUID":"a","Name":"name","Endpoint":"endpoint","Token":"token","Cert":"cert"}`,
		},
	}, nil)
	store.On("GetVersion", mock.Anything, mock.Anything).Return(1)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	sut := NewService(&mockAgentAction{}, auth, store)

	res, err := sut.Update(context.TODO(), &pb.UpdateCluster{
		Uuid: "a",
		Name: "NEWname",
	})

	assert.Equal(t, &pb.ClusterUpdated{}, res)
	assert.Nil(t, err)
}

func TestUpdateValidMultiple(t *testing.T) {
	auth := &mockAuth{}
	auth.On("Authorize", mock.Anything, mock.Anything).Return(nil)
	store := &mockStorage{}
	store.On("GetEvents", mock.Anything, "a").Return([]*storage.StoredEvent{
		{
			UUID:      "a",
			Version:   1,
			Timestamp: ddd.NewTimestamp(),
			Type:      "ClusterCreated",
			Data:      `{"Timestamp":0,"UUID":"a","Name":"name","Endpoint":"endpoint","Token":"token","Cert":"cert"}`,
		},
		{
			UUID:      "a",
			Version:   2,
			Timestamp: ddd.NewTimestamp(),
			Type:      "ClusterUpdated",
			Data:      `{"Timestamp":0,"Name":"NEWname","Endpoint":"NEWendpoint","Token":"NEWtoken","Cert":"NEWcert"}`,
		},
	}, nil)
	store.On("GetVersion", mock.Anything, mock.Anything).Return(2)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	sut := NewService(&mockAgentAction{}, auth, store)

	res, err := sut.Update(context.TODO(), &pb.UpdateCluster{
		Uuid: "a",
		Name: "NEWERname",
	})

	assert.Equal(t, &pb.ClusterUpdated{}, res)
	assert.Nil(t, err)
}

func TestUpdateLoadError(t *testing.T) {
	auth := &mockAuth{}
	auth.On("Authorize", mock.Anything, mock.Anything).Return(nil)
	store := &mockStorage{}
	store.On("GetEvents", mock.Anything, "a").Return(nil, errors.New(""))
	store.On("GetVersion", mock.Anything, mock.Anything).Return(1)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	sut := NewService(&mockAgentAction{}, auth, store)

	res, err := sut.Update(context.TODO(), &pb.UpdateCluster{
		Uuid: "a",
		Name: "NEWname",
	})

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestUpdateInvalid(t *testing.T) {
	auth := &mockAuth{}
	auth.On("Authorize", mock.Anything, mock.Anything).Return(nil)
	store := &mockStorage{}
	store.On("GetEvents", mock.Anything, "a").Return([]*storage.StoredEvent{
		{
			UUID:      "a",
			Version:   1,
			Timestamp: ddd.NewTimestamp(),
			Type:      "ClusterCreated",
			Data:      `{"Timestamp":0,"UUID":"a","Name":"name","Endpoint":"endpoint","Token":"token","Cert":"cert"}`,
		},
	}, nil)
	store.On("GetVersion", mock.Anything, mock.Anything).Return(1)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	sut := NewService(&mockAgentAction{}, auth, store)

	res, err := sut.Update(context.TODO(), &pb.UpdateCluster{
		Uuid: "a",
		Name: "",
	})

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestUpdateStoreEventError(t *testing.T) {
	auth := &mockAuth{}
	auth.On("Authorize", mock.Anything, mock.Anything).Return(nil)
	store := &mockStorage{}
	store.On("GetEvents", mock.Anything, "a").Return([]*storage.StoredEvent{
		{
			UUID:      "a",
			Version:   1,
			Timestamp: ddd.NewTimestamp(),
			Type:      "ClusterCreated",
			Data:      `{"Timestamp":0,"UUID":"a","Name":"name","Endpoint":"endpoint","Token":"token","Cert":"cert"}`,
		},
	}, nil)
	store.On("GetVersion", mock.Anything, mock.Anything).Return(1)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New(""))

	sut := NewService(&mockAgentAction{}, auth, store)

	res, err := sut.Update(context.TODO(), &pb.UpdateCluster{
		Uuid: "a",
		Name: "NEWname",
	})

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestDestroyAuth(t *testing.T) {
	a := &mockAuth{}
	a.On("Authorize", mock.Anything, mock.Anything).Return(auth.NotAuthorizedError{})
	store := &mockStorage{}
	store.On("GetEvents", mock.Anything, "a").Return([]*storage.StoredEvent{
		{
			UUID:      "a",
			Version:   1,
			Timestamp: ddd.NewTimestamp(),
			Type:      "ClusterCreated",
			Data:      `{"Timestamp":0,"UUID":"a","Name":"name","Endpoint":"endpoint","Token":"token","Cert":"cert"}`,
		},
	}, nil)
	store.On("GetVersion", mock.Anything, mock.Anything).Return(1)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	sut := NewService(&mockAgentAction{}, a, store)

	res, err := sut.Destroy(context.TODO(), &pb.DestroyCluster{
		Uuid: "a",
	})

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestDestroyValid(t *testing.T) {
	auth := &mockAuth{}
	auth.On("Authorize", mock.Anything, mock.Anything).Return(nil)
	store := &mockStorage{}
	store.On("GetEvents", mock.Anything, "a").Return([]*storage.StoredEvent{
		{
			UUID:      "a",
			Version:   1,
			Timestamp: ddd.NewTimestamp(),
			Type:      "ClusterCreated",
			Data:      `{"Timestamp":0,"UUID":"a","Name":"name","Endpoint":"endpoint","Token":"token","Cert":"cert"}`,
		},
	}, nil)
	store.On("GetVersion", mock.Anything, mock.Anything).Return(1)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	sut := NewService(&mockAgentAction{}, auth, store)

	res, err := sut.Destroy(context.TODO(), &pb.DestroyCluster{
		Uuid: "a",
	})

	assert.Equal(t, &pb.ClusterDestroyed{}, res)
	assert.Nil(t, err)
}

func TestDestroyLoadError(t *testing.T) {
	auth := &mockAuth{}
	auth.On("Authorize", mock.Anything, mock.Anything).Return(nil)
	store := &mockStorage{}
	store.On("GetEvents", mock.Anything, "a").Return(nil, errors.New(""))
	store.On("GetVersion", mock.Anything, mock.Anything).Return(1)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	sut := NewService(&mockAgentAction{}, auth, store)

	res, err := sut.Destroy(context.TODO(), &pb.DestroyCluster{
		Uuid: "a",
	})

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestDestroyAlreadyDestroyed(t *testing.T) {
	auth := &mockAuth{}
	auth.On("Authorize", mock.Anything, mock.Anything).Return(nil)
	store := &mockStorage{}
	store.On("GetEvents", mock.Anything, "a").Return([]*storage.StoredEvent{
		{
			UUID:      "a",
			Version:   1,
			Timestamp: ddd.NewTimestamp(),
			Type:      "ClusterCreated",
			Data:      `{"Timestamp":0,"UUID":"a","Name":"name","Endpoint":"endpoint","Token":"token","Cert":"cert"}`,
		},
		{
			UUID:      "a",
			Version:   2,
			Timestamp: ddd.NewTimestamp(),
			Type:      "ClusterDestroyed",
			Data:      `{"Timestamp":0}`,
		},
	}, nil)
	store.On("GetVersion", mock.Anything, mock.Anything).Return(2)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	sut := NewService(&mockAgentAction{}, auth, store)

	res, err := sut.Destroy(context.TODO(), &pb.DestroyCluster{
		Uuid: "a",
	})

	assert.IsType(t, &pb.ClusterDestroyed{}, res)
	assert.Nil(t, err)
}

func TestDestroyStoreEventError(t *testing.T) {
	auth := &mockAuth{}
	auth.On("Authorize", mock.Anything, mock.Anything).Return(nil)
	store := &mockStorage{}
	store.On("GetEvents", mock.Anything, "a").Return([]*storage.StoredEvent{
		{
			UUID:      "a",
			Version:   1,
			Timestamp: ddd.NewTimestamp(),
			Type:      "ClusterCreated",
			Data:      `{"Timestamp":0,"UUID":"a","Name":"name","Endpoint":"endpoint","Token":"token","Cert":"cert"}`,
		},
	}, nil)
	store.On("GetVersion", mock.Anything, mock.Anything).Return(1)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New(""))

	sut := NewService(&mockAgentAction{}, auth, store)

	res, err := sut.Destroy(context.TODO(), &pb.DestroyCluster{
		Uuid: "a",
	})

	assert.Nil(t, res)
	assert.Error(t, err)
}
