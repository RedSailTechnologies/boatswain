package repo

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
)

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

func TestNewRepository(t *testing.T) {
	assert.NotNil(t, NewRepository("collection", &mockStorage{}))
}

func TestAllWithValidEvents(t *testing.T) {
	store := &mockStorage{}
	store.On("IDs", mock.Anything).Return([]string{"a", "b", "c"}, nil)
	store.On("GetEvents", mock.Anything, "a").Return([]*storage.StoredEvent{
		&storage.StoredEvent{
			UUID:      "a",
			Version:   1,
			Timestamp: ddd.NewTimestamp(),
			Type:      "RepoCreated",
			Data:      `{"Timestamp":0,"UUID":"a","Name":"name","Endpoint":"http://endpoint"}`,
		},
	}, nil)
	store.On("GetEvents", mock.Anything, "b").Return([]*storage.StoredEvent{
		&storage.StoredEvent{
			UUID:      "b",
			Version:   1,
			Timestamp: ddd.NewTimestamp(),
			Type:      "RepoCreated",
			Data:      `{"Timestamp":0,"UUID":"b","Name":"name","Endpoint":"http://endpoint"}`,
		},
		&storage.StoredEvent{
			UUID:      "b",
			Version:   2,
			Timestamp: ddd.NewTimestamp(),
			Type:      "RepoUpdated",
			Data:      `{"Timestamp":0,"Name":"NEWname","Endpoint":"NEWendpoint","Token":"NEWtoken","Cert":"NEWcert"}`,
		},
	}, nil)
	store.On("GetEvents", mock.Anything, "c").Return([]*storage.StoredEvent{
		&storage.StoredEvent{
			UUID:      "c",
			Version:   1,
			Timestamp: ddd.NewTimestamp(),
			Type:      "RepoCreated",
			Data:      `{"Timestamp":0,"UUID":"c","Name":"name","Endpoint":"http://endpoint"}`,
		},
		&storage.StoredEvent{
			UUID:      "c",
			Version:   2,
			Timestamp: ddd.NewTimestamp(),
			Type:      "RepoDestroyed",
			Data:      `{"Timestamp":0}`,
		},
	}, nil)

	sut := NewRepository("collection", store)
	result, err := sut.All()

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)

	first := result[0]
	assert.Equal(t, "a", first.UUID())
	assert.Equal(t, "name", first.Name())
	assert.Equal(t, "http://endpoint", first.Endpoint())
	assert.Equal(t, 1, first.Version())

	second := result[1]
	assert.Equal(t, "b", second.UUID())
	assert.Equal(t, "NEWname", second.Name())
	assert.Equal(t, "NEWendpoint", second.Endpoint())
	assert.Equal(t, 2, second.Version())
}

func TestAllUUIDError(t *testing.T) {
	e := errors.New("test error")
	store := &mockStorage{}
	store.On("IDs", mock.Anything).Return(nil, e)

	sut := NewRepository("collection", store)
	result, err := sut.All()

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, e, err)
}

func TestAllGetEventsError(t *testing.T) {
	e := errors.New("test error")
	store := &mockStorage{}
	store.On("IDs", mock.Anything).Return([]string{"a"}, nil)
	store.On("GetEvents", mock.Anything, "a").Return(nil, e)

	sut := NewRepository("collection", store)
	result, err := sut.All()

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, e, err)
}

func TestAllGetEventsMalformedJson(t *testing.T) {
	store := &mockStorage{}
	store.On("IDs", mock.Anything).Return([]string{"a", "b", "c"}, nil)
	store.On("GetEvents", mock.Anything, "a").Return([]*storage.StoredEvent{
		&storage.StoredEvent{
			UUID:      "a",
			Version:   1,
			Timestamp: ddd.NewTimestamp(),
			Type:      "RepoCreated",
			Data:      "malformed!",
		},
	}, nil)

	sut := NewRepository("collection", store)
	result, err := sut.All()

	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestAllGetEventsUnsupportedEvent(t *testing.T) {
	e := ddd.UnsupportedEventError{
		EventType: "NotSupported",
		Type:      "Repo",
	}
	store := &mockStorage{}
	store.On("IDs", mock.Anything).Return([]string{"a", "b", "c"}, nil)
	store.On("GetEvents", mock.Anything, "a").Return([]*storage.StoredEvent{
		&storage.StoredEvent{
			UUID:      "a",
			Version:   1,
			Timestamp: ddd.NewTimestamp(),
			Type:      "NotSupported",
			Data:      "malformed!",
		},
	}, nil)

	sut := NewRepository("collection", store)
	result, err := sut.All()

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, e, err)
}

func TestLoadValid(t *testing.T) {
	store := &mockStorage{}
	store.On("GetEvents", mock.Anything, "a").Return([]*storage.StoredEvent{
		&storage.StoredEvent{
			UUID:      "a",
			Version:   1,
			Timestamp: ddd.NewTimestamp(),
			Type:      "RepoCreated",
			Data:      `{"Timestamp":0,"UUID":"a","Name":"name","Endpoint":"http://endpoint"}`,
		},
	}, nil)

	sut := NewRepository("collection", store)
	result, err := sut.Load("a")

	assert.Nil(t, err)
	assert.Equal(t, "a", result.UUID())
	assert.Equal(t, "name", result.Name())
	assert.Equal(t, "http://endpoint", result.Endpoint())
	assert.Equal(t, 1, result.Version())
}

func TestLoadStoreError(t *testing.T) {
	e := errors.New("test error")
	store := &mockStorage{}
	store.On("GetEvents", mock.Anything, "a").Return(nil, e)

	sut := NewRepository("collection", store)
	result, err := sut.Load("a")

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, e, err)
}

func TestLoadMalformedJson(t *testing.T) {
	store := &mockStorage{}
	store.On("GetEvents", mock.Anything, "a").Return([]*storage.StoredEvent{
		&storage.StoredEvent{
			UUID:      "a",
			Version:   1,
			Timestamp: ddd.NewTimestamp(),
			Type:      "RepoCreated",
			Data:      "malformed!",
		},
	}, nil)

	sut := NewRepository("collection", store)
	result, err := sut.Load("a")

	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestLoadEmptyEvents(t *testing.T) {
	store := &mockStorage{}
	store.On("GetEvents", mock.Anything, "a").Return([]*storage.StoredEvent{}, nil)

	sut := NewRepository("collection", store)
	result, err := sut.Load("a")

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, ddd.NotFoundError{Entity: "Repo"}, err)
}

func TestLoadDestroyed(t *testing.T) {
	store := &mockStorage{}
	store.On("GetEvents", mock.Anything, "a").Return([]*storage.StoredEvent{
		&storage.StoredEvent{
			UUID:      "a",
			Version:   1,
			Timestamp: ddd.NewTimestamp(),
			Type:      "RepoCreated",
			Data:      `{"Timestamp":0,"UUID":"a","Name":"name","Endpoint":"http://endpoint"}`,
		},
		&storage.StoredEvent{
			UUID:      "a",
			Version:   2,
			Timestamp: ddd.NewTimestamp(),
			Type:      "RepoDestroyed",
			Data:      `{"Timestamp":0}`,
		},
	}, nil)

	sut := NewRepository("collection", store)
	result, err := sut.Load("a")

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, ddd.DestroyedError{Entity: "Repo"}, err)
}

func TestSaveValid(t *testing.T) {
	c, _ := Create("a", "name", "http://endpoint", ddd.NewTimestamp())
	store := &mockStorage{}
	store.On("GetVersion", mock.Anything, "a").Return(0)
	store.On("StoreEvent", mock.Anything, "a", mock.Anything, mock.Anything, 1).Return(nil)
	sut := NewRepository("collection", store)
	err := sut.Save(c)
	assert.Nil(t, err)
}

func TestSaveStorageError(t *testing.T) {
	e := errors.New("test error")
	c, _ := Create("a", "name", "http://endpoint", ddd.NewTimestamp())
	store := &mockStorage{}
	store.On("GetVersion", mock.Anything, "a").Return(0)
	store.On("StoreEvent", mock.Anything, "a", mock.Anything, mock.Anything, 1).Return(e)
	sut := NewRepository("collection", store)
	err := sut.Save(c)
	assert.Error(t, err)
	assert.Equal(t, e, err)
}
