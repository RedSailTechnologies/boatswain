package repo

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/redsailtechnologies/boatswain/pkg/auth"
	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/redsailtechnologies/boatswain/pkg/git"
	"github.com/redsailtechnologies/boatswain/pkg/helm"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
	pb "github.com/redsailtechnologies/boatswain/rpc/repo"
)

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
	assert.NotNil(t, NewService(&mockAuth{}, git.DefaultAgent{}, helm.DefaultAgent{}, &mockStorage{}))
}

func TestCreateAuth(t *testing.T) {
	a := &mockAuth{}
	a.On("Authorize", mock.Anything, mock.Anything).Return(auth.NotAuthorizedError{})
	store := &mockStorage{}
	sut := NewService(a, git.DefaultAgent{}, helm.DefaultAgent{}, store)

	res, err := sut.Create(context.TODO(), &pb.CreateRepo{
		Name:     "name",
		Endpoint: "http://endpoint",
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
	sut := NewService(auth, git.DefaultAgent{}, helm.DefaultAgent{}, store)

	res, err := sut.Create(context.TODO(), &pb.CreateRepo{
		Name:     "name",
		Endpoint: "http://endpoint",
	})

	assert.Equal(t, &pb.RepoCreated{}, res)
	assert.Nil(t, err)
}

func TestCreateInvalid(t *testing.T) {
	auth := &mockAuth{}
	auth.On("Authorize", mock.Anything, mock.Anything).Return(nil)
	store := &mockStorage{}
	store.On("GetVersion", mock.Anything, mock.Anything).Return(0)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sut := NewService(auth, git.DefaultAgent{}, helm.DefaultAgent{}, store)

	res, err := sut.Create(context.TODO(), &pb.CreateRepo{
		Name:     "",
		Endpoint: "http://endpoint",
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
	sut := NewService(auth, git.DefaultAgent{}, helm.DefaultAgent{}, store)

	res, err := sut.Create(context.TODO(), &pb.CreateRepo{
		Name:     "name",
		Endpoint: "http://endpoint",
	})

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestUpdateAuth(t *testing.T) {
	a := &mockAuth{}
	a.On("Authorize", mock.Anything, mock.Anything).Return(auth.NotAuthorizedError{})
	store := &mockStorage{}
	store.On("GetEvents", mock.Anything, "a").Return([]*storage.StoredEvent{
		&storage.StoredEvent{
			UUID:      "a",
			Version:   1,
			Timestamp: ddd.NewTimestamp(),
			Type:      "RepoCreated",
			Data:      `{"Timestamp":0,"UUID":"a","Name":"name","Endpoint":"http://endpoint","Token":"token","Cert":"cert"}`,
		},
	}, nil)
	store.On("GetVersion", mock.Anything, mock.Anything).Return(1)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	sut := NewService(a, git.DefaultAgent{}, helm.DefaultAgent{}, store)

	res, err := sut.Update(context.TODO(), &pb.UpdateRepo{
		Uuid:     "a",
		Name:     "NEWname",
		Endpoint: "http://NEWendpoint",
	})

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestUpdateValid(t *testing.T) {
	auth := &mockAuth{}
	auth.On("Authorize", mock.Anything, mock.Anything).Return(nil)
	store := &mockStorage{}
	store.On("GetEvents", mock.Anything, "a").Return([]*storage.StoredEvent{
		&storage.StoredEvent{
			UUID:      "a",
			Version:   1,
			Timestamp: ddd.NewTimestamp(),
			Type:      "RepoCreated",
			Data:      `{"Timestamp":0,"UUID":"a","Name":"name","Endpoint":"http://endpoint","Token":"token","Cert":"cert"}`,
		},
	}, nil)
	store.On("GetVersion", mock.Anything, mock.Anything).Return(1)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	sut := NewService(auth, git.DefaultAgent{}, helm.DefaultAgent{}, store)

	res, err := sut.Update(context.TODO(), &pb.UpdateRepo{
		Uuid:     "a",
		Name:     "NEWname",
		Endpoint: "http://NEWendpoint",
	})

	assert.Equal(t, &pb.RepoUpdated{}, res)
	assert.Nil(t, err)
}

func TestUpdateValidMultiple(t *testing.T) {
	auth := &mockAuth{}
	auth.On("Authorize", mock.Anything, mock.Anything).Return(nil)
	store := &mockStorage{}
	store.On("GetEvents", mock.Anything, "a").Return([]*storage.StoredEvent{
		&storage.StoredEvent{
			UUID:      "a",
			Version:   1,
			Timestamp: ddd.NewTimestamp(),
			Type:      "RepoCreated",
			Data:      `{"Timestamp":0,"UUID":"a","Name":"name","Endpoint":"http://endpoint","Token":"token","Cert":"cert"}`,
		},
		&storage.StoredEvent{
			UUID:      "a",
			Version:   2,
			Timestamp: ddd.NewTimestamp(),
			Type:      "RepoUpdated",
			Data:      `{"Timestamp":0,"Name":"NEWname","Endpoint":"http://NEWendpoint","Token":"NEWtoken","Cert":"NEWcert"}`,
		},
	}, nil)
	store.On("GetVersion", mock.Anything, mock.Anything).Return(2)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	sut := NewService(auth, git.DefaultAgent{}, helm.DefaultAgent{}, store)

	res, err := sut.Update(context.TODO(), &pb.UpdateRepo{
		Uuid:     "a",
		Name:     "NEWERname",
		Endpoint: "http://NEWERendpoint",
	})

	assert.Equal(t, &pb.RepoUpdated{}, res)
	assert.Nil(t, err)
}

func TestUpdateLoadError(t *testing.T) {
	auth := &mockAuth{}
	auth.On("Authorize", mock.Anything, mock.Anything).Return(nil)
	store := &mockStorage{}
	store.On("GetEvents", mock.Anything, "a").Return(nil, errors.New(""))
	store.On("GetVersion", mock.Anything, mock.Anything).Return(1)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	sut := NewService(auth, git.DefaultAgent{}, helm.DefaultAgent{}, store)

	res, err := sut.Update(context.TODO(), &pb.UpdateRepo{
		Uuid:     "a",
		Name:     "NEWname",
		Endpoint: "http://NEWendpoint",
	})

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestUpdateInvalid(t *testing.T) {
	auth := &mockAuth{}
	auth.On("Authorize", mock.Anything, mock.Anything).Return(nil)
	store := &mockStorage{}
	store.On("GetEvents", mock.Anything, "a").Return([]*storage.StoredEvent{
		&storage.StoredEvent{
			UUID:      "a",
			Version:   1,
			Timestamp: ddd.NewTimestamp(),
			Type:      "RepoCreated",
			Data:      `{"Timestamp":0,"UUID":"a","Name":"name","Endpoint":"http://endpoint","Token":"token","Cert":"cert"}`,
		},
	}, nil)
	store.On("GetVersion", mock.Anything, mock.Anything).Return(1)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	sut := NewService(auth, git.DefaultAgent{}, helm.DefaultAgent{}, store)

	res, err := sut.Update(context.TODO(), &pb.UpdateRepo{
		Uuid:     "a",
		Name:     "",
		Endpoint: "http://NEWendpoint",
	})

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestUpdateStoreEventError(t *testing.T) {
	auth := &mockAuth{}
	auth.On("Authorize", mock.Anything, mock.Anything).Return(nil)
	store := &mockStorage{}
	store.On("GetEvents", mock.Anything, "a").Return([]*storage.StoredEvent{
		&storage.StoredEvent{
			UUID:      "a",
			Version:   1,
			Timestamp: ddd.NewTimestamp(),
			Type:      "RepoCreated",
			Data:      `{"Timestamp":0,"UUID":"a","Name":"name","Endpoint":"http://endpoint","Token":"token","Cert":"cert"}`,
		},
	}, nil)
	store.On("GetVersion", mock.Anything, mock.Anything).Return(1)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New(""))

	sut := NewService(auth, git.DefaultAgent{}, helm.DefaultAgent{}, store)

	res, err := sut.Update(context.TODO(), &pb.UpdateRepo{
		Uuid:     "a",
		Name:     "NEWname",
		Endpoint: "http://NEWendpoint",
	})

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestDestroyAuth(t *testing.T) {
	a := &mockAuth{}
	a.On("Authorize", mock.Anything, mock.Anything).Return(auth.NotAuthorizedError{})
	store := &mockStorage{}
	store.On("GetEvents", mock.Anything, "a").Return([]*storage.StoredEvent{
		&storage.StoredEvent{
			UUID:      "a",
			Version:   1,
			Timestamp: ddd.NewTimestamp(),
			Type:      "RepoCreated",
			Data:      `{"Timestamp":0,"UUID":"a","Name":"name","Endpoint":"http://endpoint","Token":"token","Cert":"cert"}`,
		},
	}, nil)
	store.On("GetVersion", mock.Anything, mock.Anything).Return(1)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	sut := NewService(a, git.DefaultAgent{}, helm.DefaultAgent{}, store)

	res, err := sut.Destroy(context.TODO(), &pb.DestroyRepo{
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
		&storage.StoredEvent{
			UUID:      "a",
			Version:   1,
			Timestamp: ddd.NewTimestamp(),
			Type:      "RepoCreated",
			Data:      `{"Timestamp":0,"UUID":"a","Name":"name","Endpoint":"http://endpoint","Token":"token","Cert":"cert"}`,
		},
	}, nil)
	store.On("GetVersion", mock.Anything, mock.Anything).Return(1)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	sut := NewService(auth, git.DefaultAgent{}, helm.DefaultAgent{}, store)

	res, err := sut.Destroy(context.TODO(), &pb.DestroyRepo{
		Uuid: "a",
	})

	assert.Equal(t, &pb.RepoDestroyed{}, res)
	assert.Nil(t, err)
}

func TestDestroyLoadError(t *testing.T) {
	auth := &mockAuth{}
	auth.On("Authorize", mock.Anything, mock.Anything).Return(nil)
	store := &mockStorage{}
	store.On("GetEvents", mock.Anything, "a").Return(nil, errors.New(""))
	store.On("GetVersion", mock.Anything, mock.Anything).Return(1)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	sut := NewService(auth, git.DefaultAgent{}, helm.DefaultAgent{}, store)

	res, err := sut.Destroy(context.TODO(), &pb.DestroyRepo{
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
		&storage.StoredEvent{
			UUID:      "a",
			Version:   1,
			Timestamp: ddd.NewTimestamp(),
			Type:      "RepoCreated",
			Data:      `{"Timestamp":0,"UUID":"a","Name":"name","Endpoint":"http://endpoint","Token":"token","Cert":"cert"}`,
		},
		&storage.StoredEvent{
			UUID:      "a",
			Version:   2,
			Timestamp: ddd.NewTimestamp(),
			Type:      "RepoDestroyed",
			Data:      `{"Timestamp":0}`,
		},
	}, nil)
	store.On("GetVersion", mock.Anything, mock.Anything).Return(2)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	sut := NewService(auth, git.DefaultAgent{}, helm.DefaultAgent{}, store)

	res, err := sut.Destroy(context.TODO(), &pb.DestroyRepo{
		Uuid: "a",
	})

	assert.IsType(t, &pb.RepoDestroyed{}, res)
	assert.Nil(t, err)
}

func TestDestroyStoreEventError(t *testing.T) {
	auth := &mockAuth{}
	auth.On("Authorize", mock.Anything, mock.Anything).Return(nil)
	store := &mockStorage{}
	store.On("GetEvents", mock.Anything, "a").Return([]*storage.StoredEvent{
		&storage.StoredEvent{
			UUID:      "a",
			Version:   1,
			Timestamp: ddd.NewTimestamp(),
			Type:      "RepoCreated",
			Data:      `{"Timestamp":0,"UUID":"a","Name":"name","Endpoint":"http://endpoint","Token":"token","Cert":"cert"}`,
		},
	}, nil)
	store.On("GetVersion", mock.Anything, mock.Anything).Return(1)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New(""))

	sut := NewService(auth, git.DefaultAgent{}, helm.DefaultAgent{}, store)

	res, err := sut.Destroy(context.TODO(), &pb.DestroyRepo{
		Uuid: "a",
	})

	assert.Nil(t, res)
	assert.Error(t, err)
}

// buildChartUrl
func TestBuildChartURL(t *testing.T) {
	assert.Equal(t, "http://repo.com/chart", buildChartURL("http://repo.com", "chart"))
	assert.Equal(t, "http://repo.com/chart", buildChartURL("http://repo.com/", "chart"))
}

// toChartRepo
func TestToChartRepo(t *testing.T) {
	sut, _ := Create(ddd.NewUUID(), "name", "https://endpoint", HELM, ddd.NewTimestamp())
	cr, err := sut.toChartRepo()
	assert.Nil(t, err)
	assert.NotNil(t, cr)
}
