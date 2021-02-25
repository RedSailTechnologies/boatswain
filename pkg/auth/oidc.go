package auth

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/twitchtv/twirp"

	"github.com/redsailtechnologies/boatswain/pkg/cfg"
	"github.com/redsailtechnologies/boatswain/pkg/logger"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

// Flags initializes the agent's configuration
func Flags() *Config {
	config := Config{}
	flag.StringVar(&config.OIDC, "oidc-url", cfg.EnvOrDefaultString("OIDC_URL", ""), "openid connect configuration url")
	flag.StringVar(&config.Scope, "user-scope", cfg.EnvOrDefaultString("USER_SCOPE", ""), "user scope")
	flag.StringVar(&config.AdminRole, "user-admin-role", cfg.EnvOrDefaultString("USER_ADMIN_ROLE", "Boatswain.Admin"), "user admin role")
	flag.StringVar(&config.EditorRole, "user-editor-role", cfg.EnvOrDefaultString("USER_EDITOR_ROLE", "Boatswain.Editor"), "user editor role")
	flag.StringVar(&config.ReaderRole, "user-reader-role", cfg.EnvOrDefaultString("USER_READER_ROLE", "Boatswain.Reader"), "user reader role")
	return &config
}

// Config is the configuration needed to run an auth service
type Config struct {
	OIDC      string
	Endpoints *oidcConfig

	Scope      string
	AdminRole  string
	EditorRole string
	ReaderRole string
}

type oidcConfig struct {
	JWKS string `json:"jwks_uri"`
}

// OIDCAgent is an auth agent implementation using oidc
type OIDCAgent struct {
	cfg     *Config
	jwtKey  *int
	userKey *int
	jwks    map[string]*jose.JSONWebKey
}

// NewOIDCAgent builds a new agent from the configuration
func NewOIDCAgent(config *Config) *OIDCAgent {
	req, err := http.NewRequest(http.MethodGet, config.OIDC, nil)
	if err != nil {
		logger.Fatal("could not create request for oidc configuration", "error", err)
	}

	var client http.Client
	resp, err := client.Do(req)
	if err != nil {
		logger.Fatal("error getting oidc configuration", "error", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&config.Endpoints)
	if err != nil {
		logger.Fatal("error reading oidc configuration", "error", err)
	}
	return &OIDCAgent{
		cfg:     config,
		jwtKey:  new(int),
		userKey: new(int),
		jwks:    make(map[string]*jose.JSONWebKey),
	}
}

// Authenticate handles validation of auth tokens stored in the context
func (o *OIDCAgent) Authenticate(ctx context.Context) (context.Context, error) {
	tokenVal := ctx.Value(o.jwtKey)
	if tokenVal == nil {
		logger.Warn("jwt token not found")
		return ctx, twirp.NewError(twirp.Unauthenticated, "couldn't get jwt from Authorization header")
	}
	token := tokenVal.(string)

	parsedToken, err := jwt.ParseSigned(token)
	if err != nil {
		logger.Error("couldn't parse signed jwt", "error", err)
		return ctx, twirp.NewError(twirp.Unauthenticated, "couldn't parse signed jwt")
	}

	if len(parsedToken.Headers) < 1 {
		logger.Error("parsed token did not contain any headers", "token", parsedToken)
		return ctx, twirp.NewError(twirp.Unauthenticated, "parsed token did not contain any headers")
	}
	keyID := parsedToken.Headers[0].KeyID

	err = o.getJWK(keyID)
	if err != nil {
		logger.Error("couldn't get key to verify jwt", "error", err)
		return ctx, twirp.NewError(twirp.Unauthenticated, "couldn't get key to verify jwt")
	}
	key := o.jwks[keyID]

	user := user{}
	err = parsedToken.Claims(key.Key, &user)
	if err != nil {
		logger.Error("couldn't parse user claims", "error", err)
		return ctx, twirp.NewError(twirp.Unauthenticated, "couldn't parse user claims")
	}

	if err = user.validateScope(o.cfg.Scope); err != nil {
		logger.Warn("error validating user scope", "error", err)
		return ctx, twirp.NewError(twirp.Unauthenticated, "error validating user scope")
	}

	return context.WithValue(ctx, o.userKey, user), nil
}

// Authorize verifies a context's user has the given access level
func (o *OIDCAgent) Authorize(ctx context.Context, role Role) error {
	user := o.userFromContext(ctx)
	if user == nil {
		logger.Error("user not found in context")
		return NotAuthorizedError{}
	}

	switch role {
	case Reader:
		if user.hasRole(o.cfg.ReaderRole) {
			return nil
		}
		fallthrough
	case Editor:
		if user.hasRole(o.cfg.EditorRole) {
			return nil
		}
		fallthrough
	case Admin:
		if user.hasRole(o.cfg.AdminRole) {
			return nil
		}
	}

	service, _ := twirp.ServiceName(ctx)
	method, _ := twirp.MethodName(ctx)
	logger.Error(fmt.Sprintf("user not authorized for %s.%s", service, method), "user", user)
	return NotAuthorizedError{}
}

// User gets the user from the context
func (o *OIDCAgent) User(ctx context.Context) User {
	u := o.userFromContext(ctx)
	return o.toUser(u)
}

// Wrap wraps an existing http hander to store the auth JWT
func (o *OIDCAgent) Wrap(base http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		split := strings.Split(header, " ")
		if len(split) == 2 {
			jwt := strings.Split(header, " ")[1] // get the token from "Bearer token"
			ctx := r.Context()
			ctx = context.WithValue(ctx, o.jwtKey, jwt)
			r = r.WithContext(ctx)
		}
		base.ServeHTTP(w, r)
	})
}

func (o *OIDCAgent) toUser(u *user) User {
	user := User{
		Name:    u.Name,
		Email:   u.Email,
		Subject: u.Subject,
	}

	roles := make([]Role, len(u.Roles))
	for i, r := range u.Roles {
		switch r {
		case o.cfg.AdminRole:
			roles[i] = Admin
		case o.cfg.EditorRole:
			roles[i] = Editor
		case o.cfg.ReaderRole:
			roles[i] = Reader
		}
	}
	user.Roles = roles

	return user
}

func (o *OIDCAgent) getJWK(key string) error {
	if _, ok := o.jwks[key]; ok {
		return nil
	}
	req, err := http.NewRequest(http.MethodGet, o.cfg.Endpoints.JWKS, nil)
	if err != nil {
		return err
	}

	var client http.Client
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var jwks jose.JSONWebKeySet
	err = json.NewDecoder(resp.Body).Decode(&jwks)
	if err != nil {
		return err
	}

	keys := jwks.Key(key)
	if len(keys) < 1 {
		return errors.New("no jwks keys found")
	} else if len(keys) > 1 {
		logger.Warn("more than one jwks key found", "keys", key)
	}

	o.jwks[key] = &keys[0]
	return nil
}

func (o *OIDCAgent) userFromContext(ctx context.Context) *user {
	userVal := ctx.Value(o.userKey)
	if userVal == nil {
		return nil
	}
	user := userVal.(user)
	return &user
}
