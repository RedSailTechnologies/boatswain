package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthorizationNoUser(t *testing.T) {
	ctx := context.TODO()
	cfg := &Config{
		AdminRole:  "admin",
		EditorRole: "editor",
		ReaderRole: "reader",
	}
	sut := &OIDCAgent{cfg: cfg}

	adminErr := sut.Authorize(ctx, Admin)
	editorErr := sut.Authorize(ctx, Editor)
	readerErr := sut.Authorize(ctx, Reader)

	assert.Error(t, adminErr)
	assert.Error(t, editorErr)
	assert.Error(t, readerErr)
}

func TestAuthorizationValidAdmin(t *testing.T) {
	ctx := context.TODO()
	user := user{
		Name: "Some User",
		Roles: []string{
			"admin",
		},
	}
	cfg := &Config{
		AdminRole:  "admin",
		EditorRole: "editor",
		ReaderRole: "reader",
	}
	sut := &OIDCAgent{cfg: cfg}
	ctx = context.WithValue(ctx, sut.userKey, user)

	adminErr := sut.Authorize(ctx, Admin)
	editorErr := sut.Authorize(ctx, Editor)
	readerErr := sut.Authorize(ctx, Reader)

	assert.Nil(t, adminErr)
	assert.Nil(t, editorErr)
	assert.Nil(t, readerErr)
}

func TestAuthorizationInvalidAdmin(t *testing.T) {
	ctx := context.TODO()
	user := user{
		Name: "Some User",
		Roles: []string{
			"administrator!",
		},
	}
	cfg := &Config{
		AdminRole:  "admin",
		EditorRole: "editor",
		ReaderRole: "reader",
	}
	sut := &OIDCAgent{cfg: cfg}
	ctx = context.WithValue(ctx, sut.userKey, user)

	adminErr := sut.Authorize(ctx, Admin)
	editorErr := sut.Authorize(ctx, Editor)
	readerErr := sut.Authorize(ctx, Reader)

	assert.Error(t, adminErr)
	assert.Error(t, editorErr)
	assert.Error(t, readerErr)
}

func TestAuthorizationValidEditor(t *testing.T) {
	ctx := context.TODO()
	user := user{
		Name: "Some User",
		Roles: []string{
			"editor",
		},
	}
	cfg := &Config{
		AdminRole:  "admin",
		EditorRole: "editor",
		ReaderRole: "reader",
	}
	sut := &OIDCAgent{cfg: cfg}
	ctx = context.WithValue(ctx, sut.userKey, user)

	adminErr := sut.Authorize(ctx, Admin)
	editorErr := sut.Authorize(ctx, Editor)
	readerErr := sut.Authorize(ctx, Reader)

	assert.Error(t, adminErr)
	assert.Nil(t, editorErr)
	assert.Nil(t, readerErr)
}

func TestAuthorizationInvalidEditor(t *testing.T) {
	ctx := context.TODO()
	user := user{
		Name: "Some User",
		Roles: []string{
			"edit!",
		},
	}
	cfg := &Config{
		AdminRole:  "admin",
		EditorRole: "editor",
		ReaderRole: "reader",
	}
	sut := &OIDCAgent{cfg: cfg}
	ctx = context.WithValue(ctx, sut.userKey, user)

	adminErr := sut.Authorize(ctx, Admin)
	editorErr := sut.Authorize(ctx, Editor)
	readerErr := sut.Authorize(ctx, Reader)

	assert.Error(t, adminErr)
	assert.Error(t, editorErr)
	assert.Error(t, readerErr)
}

func TestAuthorizationValidReader(t *testing.T) {
	ctx := context.TODO()
	user := user{
		Name: "Some User",
		Roles: []string{
			"reader",
		},
	}
	cfg := &Config{
		AdminRole:  "admin",
		EditorRole: "editor",
		ReaderRole: "reader",
	}
	sut := &OIDCAgent{cfg: cfg}
	ctx = context.WithValue(ctx, sut.userKey, user)

	adminErr := sut.Authorize(ctx, Admin)
	editorErr := sut.Authorize(ctx, Editor)
	readerErr := sut.Authorize(ctx, Reader)

	assert.Error(t, adminErr)
	assert.Error(t, editorErr)
	assert.Nil(t, readerErr)
}

func TestAuthorizationInvalidReader(t *testing.T) {
	ctx := context.TODO()
	user := user{
		Name: "Some User",
		Roles: []string{
			"README.md",
		},
	}
	cfg := &Config{
		AdminRole:  "admin",
		EditorRole: "editor",
		ReaderRole: "reader",
	}
	sut := &OIDCAgent{cfg: cfg}
	ctx = context.WithValue(ctx, sut.userKey, user)

	adminErr := sut.Authorize(ctx, Admin)
	editorErr := sut.Authorize(ctx, Editor)
	readerErr := sut.Authorize(ctx, Reader)

	assert.Error(t, adminErr)
	assert.Error(t, editorErr)
	assert.Error(t, readerErr)
}

func TestFlagsReturnsConfig(t *testing.T) {
	var cfg *Config
	cfg = Flags()
	assert.NotNil(t, cfg)
}
