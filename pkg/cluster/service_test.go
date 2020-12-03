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
	"github.com/redsailtechnologies/boatswain/pkg/kube"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
	pb "github.com/redsailtechnologies/boatswain/rpc/cluster"
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

func (ma *mockAuth) Wrap(h http.Handler) http.Handler {
	ma.Called(h)
	return h
}

func TestNewService(t *testing.T) {
	assert.NotNil(t, NewService(&mockAuth{}, kube.DefaultAgent{}, &mockStorage{}))
}

func TestCreateAuth(t *testing.T) {
	a := &mockAuth{}
	a.On("Authorize", mock.Anything, mock.Anything).Return(auth.NotAuthorizedError{})
	store := &mockStorage{}
	sut := NewService(a, kube.DefaultAgent{}, store)

	res, err := sut.Create(context.TODO(), &pb.CreateCluster{
		Name:     "name",
		Endpoint: "endpoint",
		Token:    "token",
		Cert:     "cert",
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
	sut := NewService(auth, kube.DefaultAgent{}, store)

	res, err := sut.Create(context.TODO(), &pb.CreateCluster{
		Name:     "name",
		Endpoint: "endpoint",
		Token:    "token",
		Cert:     "cert",
	})

	assert.Equal(t, &pb.ClusterCreated{}, res)
	assert.Nil(t, err)
}

func TestCreateInvalid(t *testing.T) {
	auth := &mockAuth{}
	auth.On("Authorize", mock.Anything, mock.Anything).Return(nil)
	store := &mockStorage{}
	store.On("GetVersion", mock.Anything, mock.Anything).Return(0)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sut := NewService(auth, kube.DefaultAgent{}, store)

	res, err := sut.Create(context.TODO(), &pb.CreateCluster{
		Name:     "",
		Endpoint: "endpoint",
		Token:    "token",
		Cert:     "cert",
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
	sut := NewService(auth, kube.DefaultAgent{}, store)

	res, err := sut.Create(context.TODO(), &pb.CreateCluster{
		Name:     "name",
		Endpoint: "endpoint",
		Token:    "token",
		Cert:     "cert",
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
			Type:      "ClusterCreated",
			Data:      `{"Timestamp":0,"UUID":"a","Name":"name","Endpoint":"endpoint","Token":"token","Cert":"cert"}`,
		},
	}, nil)
	store.On("GetVersion", mock.Anything, mock.Anything).Return(1)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sut := NewService(a, kube.DefaultAgent{}, store)

	res, err := sut.Update(context.TODO(), &pb.UpdateCluster{
		Uuid:     "a",
		Name:     "NEWname",
		Endpoint: "NEWendpoint",
		Token:    "NEWtoken",
		Cert:     "NEWcert",
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
			Type:      "ClusterCreated",
			Data:      `{"Timestamp":0,"UUID":"a","Name":"name","Endpoint":"endpoint","Token":"token","Cert":"cert"}`,
		},
	}, nil)
	store.On("GetVersion", mock.Anything, mock.Anything).Return(1)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	sut := NewService(auth, kube.DefaultAgent{}, store)

	res, err := sut.Update(context.TODO(), &pb.UpdateCluster{
		Uuid:     "a",
		Name:     "NEWname",
		Endpoint: "NEWendpoint",
		Token:    "NEWtoken",
		Cert:     "NEWcert",
	})

	assert.Equal(t, &pb.ClusterUpdated{}, res)
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
			Type:      "ClusterCreated",
			Data:      `{"Timestamp":0,"UUID":"a","Name":"name","Endpoint":"endpoint","Token":"token","Cert":"cert"}`,
		},
		&storage.StoredEvent{
			UUID:      "a",
			Version:   2,
			Timestamp: ddd.NewTimestamp(),
			Type:      "ClusterUpdated",
			Data:      `{"Timestamp":0,"Name":"NEWname","Endpoint":"NEWendpoint","Token":"NEWtoken","Cert":"NEWcert"}`,
		},
	}, nil)
	store.On("GetVersion", mock.Anything, mock.Anything).Return(2)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	sut := NewService(auth, kube.DefaultAgent{}, store)

	res, err := sut.Update(context.TODO(), &pb.UpdateCluster{
		Uuid:     "a",
		Name:     "NEWERname",
		Endpoint: "NEWERendpoint",
		Token:    "NEWERtoken",
		Cert:     "NEWERcert",
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

	sut := NewService(auth, kube.DefaultAgent{}, store)

	res, err := sut.Update(context.TODO(), &pb.UpdateCluster{
		Uuid:     "a",
		Name:     "NEWname",
		Endpoint: "NEWendpoint",
		Token:    "NEWtoken",
		Cert:     "NEWcert",
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
			Type:      "ClusterCreated",
			Data:      `{"Timestamp":0,"UUID":"a","Name":"name","Endpoint":"endpoint","Token":"token","Cert":"cert"}`,
		},
	}, nil)
	store.On("GetVersion", mock.Anything, mock.Anything).Return(1)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	sut := NewService(auth, kube.DefaultAgent{}, store)

	res, err := sut.Update(context.TODO(), &pb.UpdateCluster{
		Uuid:     "a",
		Name:     "",
		Endpoint: "NEWendpoint",
		Token:    "NEWtoken",
		Cert:     "NEWcert",
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
			Type:      "ClusterCreated",
			Data:      `{"Timestamp":0,"UUID":"a","Name":"name","Endpoint":"endpoint","Token":"token","Cert":"cert"}`,
		},
	}, nil)
	store.On("GetVersion", mock.Anything, mock.Anything).Return(1)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New(""))

	sut := NewService(auth, kube.DefaultAgent{}, store)

	res, err := sut.Update(context.TODO(), &pb.UpdateCluster{
		Uuid:     "a",
		Name:     "NEWname",
		Endpoint: "NEWendpoint",
		Token:    "NEWtoken",
		Cert:     "NEWcert",
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
			Type:      "ClusterCreated",
			Data:      `{"Timestamp":0,"UUID":"a","Name":"name","Endpoint":"endpoint","Token":"token","Cert":"cert"}`,
		},
	}, nil)
	store.On("GetVersion", mock.Anything, mock.Anything).Return(1)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	sut := NewService(a, kube.DefaultAgent{}, store)

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
		&storage.StoredEvent{
			UUID:      "a",
			Version:   1,
			Timestamp: ddd.NewTimestamp(),
			Type:      "ClusterCreated",
			Data:      `{"Timestamp":0,"UUID":"a","Name":"name","Endpoint":"endpoint","Token":"token","Cert":"cert"}`,
		},
	}, nil)
	store.On("GetVersion", mock.Anything, mock.Anything).Return(1)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	sut := NewService(auth, kube.DefaultAgent{}, store)

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

	sut := NewService(auth, kube.DefaultAgent{}, store)

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
		&storage.StoredEvent{
			UUID:      "a",
			Version:   1,
			Timestamp: ddd.NewTimestamp(),
			Type:      "ClusterCreated",
			Data:      `{"Timestamp":0,"UUID":"a","Name":"name","Endpoint":"endpoint","Token":"token","Cert":"cert"}`,
		},
		&storage.StoredEvent{
			UUID:      "a",
			Version:   2,
			Timestamp: ddd.NewTimestamp(),
			Type:      "ClusterDestroyed",
			Data:      `{"Timestamp":0}`,
		},
	}, nil)
	store.On("GetVersion", mock.Anything, mock.Anything).Return(2)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	sut := NewService(auth, kube.DefaultAgent{}, store)

	res, err := sut.Destroy(context.TODO(), &pb.DestroyCluster{
		Uuid: "a",
	})

	assert.Equal(t, &pb.ClusterDestroyed{}, res)
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
			Type:      "ClusterCreated",
			Data:      `{"Timestamp":0,"UUID":"a","Name":"name","Endpoint":"endpoint","Token":"token","Cert":"cert"}`,
		},
	}, nil)
	store.On("GetVersion", mock.Anything, mock.Anything).Return(1)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New(""))

	sut := NewService(auth, kube.DefaultAgent{}, store)

	res, err := sut.Destroy(context.TODO(), &pb.DestroyCluster{
		Uuid: "a",
	})

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestToClientset(t *testing.T) {
	sut, _ := Create(ddd.NewUUID(), "name", "endpoint", "token", "cert", ddd.NewTimestamp())
	cs, err := sut.toClientset()
	assert.Nil(t, err)
	assert.NotNil(t, cs)
}
