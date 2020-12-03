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
	flag.StringVar(&config.AdminRole, "user-admin-role", cfg.EnvOrDefaultString("USER_ADMIN_ROLE", ""), "user admin role")
	flag.StringVar(&config.EditorRole, "user-editor-role", cfg.EnvOrDefaultString("USER_EDITOR_ROLE", ""), "user editor role")
	flag.StringVar(&config.ReaderRole, "user-reader-role", cfg.EnvOrDefaultString("USER_READER_ROLE", ""), "user reader role")
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
}

// NewOIDCAgent builds a new agent from the configuration
func NewOIDCAgent(config *Config) *OIDCAgent {
	req, err := http.NewRequest(http.MethodGet, config.OIDC, nil)
	if err != nil {
		logger.Panic("could not create request for oidc configuration", "error", err)
	}

	var client http.Client
	resp, err := client.Do(req)
	if err != nil {
		logger.Panic("error getting oidc configuration", "error", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&config.Endpoints)
	if err != nil {
		logger.Panic("error reading oidc configuration", "error", err)
	}
	return &OIDCAgent{
		cfg:     config,
		jwtKey:  new(int),
		userKey: new(int),
	}
}

// Authenticate handles validation of auth tokens stored in the context
func (o *OIDCAgent) Authenticate(ctx context.Context) (context.Context, error) {
	tokenVal := ctx.Value(o.jwtKey)
	if tokenVal == nil {
		logger.Warn("jwt token not found")
		return nil, AuthenticationError{}
	}
	token := tokenVal.(string)

	parsedToken, err := jwt.ParseSigned(token)
	if err != nil {
		logger.Error("couldn't parse signed jwt", "error", err)
		return nil, AuthenticationError{}
	}

	if len(parsedToken.Headers) < 1 {
		logger.Error("parsed token did not contain any headers", "token", parsedToken)
		return nil, AuthenticationError{}
	}
	keyID := parsedToken.Headers[0].KeyID

	key, err := o.getJWK(keyID)
	if err != nil {
		logger.Error("couldn't get key to verify jwt", "error", err)
		return nil, AuthenticationError{}
	}

	user := User{}
	err = parsedToken.Claims(key.Key, &user)
	if err != nil {
		logger.Error("couldn't parse user claims", "error", err)
		return nil, AuthenticationError{}
	}

	if err = user.validateScope(o.cfg.Scope); err != nil {
		logger.Warn("error validating user scope", "error", err)
		return nil, AuthenticationError{}
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

// TODO - cache this on a reasonable refresh timeline (and/or retry on failure?)
func (o *OIDCAgent) getJWK(key string) (*jose.JSONWebKey, error) {
	req, err := http.NewRequest(http.MethodGet, o.cfg.Endpoints.JWKS, nil)
	if err != nil {
		return nil, err
	}

	var client http.Client
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jwks jose.JSONWebKeySet
	err = json.NewDecoder(resp.Body).Decode(&jwks)
	if err != nil {
		return nil, err
	}

	keys := jwks.Key(key)
	if len(keys) < 1 {
		return nil, errors.New("no jwks keys found")
	} else if len(keys) > 1 {
		logger.Warn("more than one jwks key found", "keys", key)
	}

	return &keys[0], nil
}

func (o *OIDCAgent) userFromContext(ctx context.Context) *User {
	userVal := ctx.Value(o.userKey)
	if userVal == nil {
		return nil
	}
	user := userVal.(User)
	return &user
}
