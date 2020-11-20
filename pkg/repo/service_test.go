package repo

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/redsailtechnologies/boatswain/pkg/helm"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
	pb "github.com/redsailtechnologies/boatswain/rpc/repo"
)

func TestNewService(t *testing.T) {
	assert.NotNil(t, NewService(helm.DefaultAgent{}, &mockStorage{}))
}

func TestCreateValid(t *testing.T) {
	store := &mockStorage{}
	store.On("GetVersion", mock.Anything, mock.Anything).Return(0)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sut := NewService(helm.DefaultAgent{}, store)

	res, err := sut.Create(context.TODO(), &pb.CreateRepo{
		Name:     "name",
		Endpoint: "http://endpoint",
	})

	assert.Equal(t, &pb.RepoCreated{}, res)
	assert.Nil(t, err)
}

func TestCreateInvalid(t *testing.T) {
	store := &mockStorage{}
	store.On("GetVersion", mock.Anything, mock.Anything).Return(0)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sut := NewService(helm.DefaultAgent{}, store)

	res, err := sut.Create(context.TODO(), &pb.CreateRepo{
		Name:     "",
		Endpoint: "http://endpoint",
	})

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestCreateSaveError(t *testing.T) {
	store := &mockStorage{}
	store.On("GetVersion", mock.Anything, mock.Anything).Return(0)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New(""))
	sut := NewService(helm.DefaultAgent{}, store)

	res, err := sut.Create(context.TODO(), &pb.CreateRepo{
		Name:     "name",
		Endpoint: "http://endpoint",
	})

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestUpdateValid(t *testing.T) {
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

	sut := NewService(helm.DefaultAgent{}, store)

	res, err := sut.Update(context.TODO(), &pb.UpdateRepo{
		Uuid:     "a",
		Name:     "NEWname",
		Endpoint: "http://NEWendpoint",
	})

	assert.Equal(t, &pb.RepoUpdated{}, res)
	assert.Nil(t, err)
}

func TestUpdateValidMultiple(t *testing.T) {
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

	sut := NewService(helm.DefaultAgent{}, store)

	res, err := sut.Update(context.TODO(), &pb.UpdateRepo{
		Uuid:     "a",
		Name:     "NEWERname",
		Endpoint: "http://NEWERendpoint",
	})

	assert.Equal(t, &pb.RepoUpdated{}, res)
	assert.Nil(t, err)
}

func TestUpdateLoadError(t *testing.T) {
	store := &mockStorage{}
	store.On("GetEvents", mock.Anything, "a").Return(nil, errors.New(""))
	store.On("GetVersion", mock.Anything, mock.Anything).Return(1)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	sut := NewService(helm.DefaultAgent{}, store)

	res, err := sut.Update(context.TODO(), &pb.UpdateRepo{
		Uuid:     "a",
		Name:     "NEWname",
		Endpoint: "http://NEWendpoint",
	})

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestUpdateInvalid(t *testing.T) {
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

	sut := NewService(helm.DefaultAgent{}, store)

	res, err := sut.Update(context.TODO(), &pb.UpdateRepo{
		Uuid:     "a",
		Name:     "",
		Endpoint: "http://NEWendpoint",
	})

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestUpdateStoreEventError(t *testing.T) {
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

	sut := NewService(helm.DefaultAgent{}, store)

	res, err := sut.Update(context.TODO(), &pb.UpdateRepo{
		Uuid:     "a",
		Name:     "NEWname",
		Endpoint: "http://NEWendpoint",
	})

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestDestroyValid(t *testing.T) {
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

	sut := NewService(helm.DefaultAgent{}, store)

	res, err := sut.Destroy(context.TODO(), &pb.DestroyRepo{
		Uuid: "a",
	})

	assert.Equal(t, &pb.RepoDestroyed{}, res)
	assert.Nil(t, err)
}

func TestDestroyLoadError(t *testing.T) {
	store := &mockStorage{}
	store.On("GetEvents", mock.Anything, "a").Return(nil, errors.New(""))
	store.On("GetVersion", mock.Anything, mock.Anything).Return(1)
	store.On("StoreEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	sut := NewService(helm.DefaultAgent{}, store)

	res, err := sut.Destroy(context.TODO(), &pb.DestroyRepo{
		Uuid: "a",
	})

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestDestroyAlreadyDestroyed(t *testing.T) {
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

	sut := NewService(helm.DefaultAgent{}, store)

	res, err := sut.Destroy(context.TODO(), &pb.DestroyRepo{
		Uuid: "a",
	})

	assert.Equal(t, &pb.RepoDestroyed{}, res)
	assert.Nil(t, err)
}

func TestDestroyStoreEventError(t *testing.T) {
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

	sut := NewService(helm.DefaultAgent{}, store)

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
	sut, _ := Create(ddd.NewUUID(), "name", "https://endpoint", ddd.NewTimestamp())
	cr, err := sut.toChartRepo()
	assert.Nil(t, err)
	assert.NotNil(t, cr)
}
