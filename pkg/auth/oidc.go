package auth

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"net/http"
	"strings"

	"github.com/redsailtechnologies/boatswain/pkg/cfg"
	"github.com/redsailtechnologies/boatswain/pkg/logger"

	"github.com/twitchtv/twirp"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

var config = &authConfig{}
var jwtKey = new(int)

// Flags initializes the auth package's configuration
func Flags() {
	flag.StringVar(&config.oidc, "oidc-url", cfg.EnvOrDefaultString("OIDC_URL", ""), "openid connect configuration url")
	flag.StringVar(&config.scope, "user-scope", cfg.EnvOrDefaultString("USER_SCOPE", ""), "user scope")
	flag.StringVar(&config.adminRole, "user-admin-role", cfg.EnvOrDefaultString("USER_ADMIN_ROLE", ""), "user admin role")
	flag.StringVar(&config.editorRole, "user-editor-role", cfg.EnvOrDefaultString("USER_EDITOR_ROLE", ""), "user editor role")
	flag.StringVar(&config.readerRole, "user-reader-role", cfg.EnvOrDefaultString("USER_READER_ROLE", ""), "user reader role")
}

// WithJWT wraps an existing http hander to store the auth JWT
func WithJWT(base http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		split := strings.Split(header, " ")
		if len(split) == 2 {
			jwt := strings.Split(header, " ")[1] // get the token from "Bearer token"
			ctx := r.Context()
			ctx = context.WithValue(ctx, jwtKey, jwt)
			r = r.WithContext(ctx)
		}
		base.ServeHTTP(w, r)
	})
}

// ValidateJWT handles validation of auth tokens stored in the context
func ValidateJWT(ctx context.Context) (context.Context, error) {
	tokenVal := ctx.Value(jwtKey)
	if tokenVal == nil {
		return ctx, twirp.NewError(twirp.Unauthenticated, "no jwt provided")
	}
	token := tokenVal.(string)

	parsedToken, err := jwt.ParseSigned(token)
	if err != nil {
		logger.Error("couldn't parse signed jwt", "error", err)
		return ctx, twirp.NewError(twirp.Unauthenticated, "couldn't parse jwt")
	}

	if len(parsedToken.Headers) < 1 {
		logger.Error("parsed token did not contain any headers", "token", parsedToken)
		return ctx, twirp.NewError(twirp.Unauthenticated, "couldn't validate jwt")
	}
	keyID := parsedToken.Headers[0].KeyID

	key, err := getJWK(config.jwks(), keyID)
	if err != nil {
		logger.Error("couldn't get key to verify jwt", "error", err)
		return ctx, twirp.NewError(twirp.Unauthenticated, "couldn't get key to verify jwt")
	}

	user := User{}
	test := make(map[string]interface{})
	err = parsedToken.Claims(key.Key, &user, &test)
	if err != nil {
		logger.Error("couldn't parse user claims", "error", err)
		return ctx, twirp.NewError(twirp.Unauthenticated, "couldn't parse user claims")
	}

	if err = user.ValidateScope(); err != nil {
		logger.Warn("error validating user scope", "error", err)
		return ctx, twirp.NewError(twirp.Unauthenticated, "user doesn't have this application's scope")
	}

	return user.AddToContext(ctx), nil
}

type authConfig struct {
	oidc      string
	endpoints *oidcConfig

	scope      string
	adminRole  string
	editorRole string
	readerRole string
}

type oidcConfig struct {
	Issuer string `json:"issuer"`
	JWKS   string `json:"jwks_uri"`
}

func (c *authConfig) issuer() string {
	if config.endpoints == nil {
		config.endpoints = &oidcConfig{}
		req, err := http.NewRequest(http.MethodGet, c.oidc, nil)
		if err != nil {
			logger.Panic("could not create request for oidc configuration")
		}

		var client http.Client
		resp, err := client.Do(req)
		if err != nil {
			logger.Panic("error getting oidc configuration")
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&config.endpoints)
		if err != nil {
			logger.Panic("error reading oidc configuration")
		}
	}
	return c.endpoints.Issuer
}

func (c *authConfig) jwks() string {
	if config.endpoints == nil {
		config.endpoints = &oidcConfig{}
		req, err := http.NewRequest(http.MethodGet, c.oidc, nil)
		if err != nil {
			logger.Panic("could not create request for oidc configuration")
		}

		var client http.Client
		resp, err := client.Do(req)
		if err != nil {
			logger.Panic("error getting oidc configuration")
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&config.endpoints)
		if err != nil {
			logger.Panic("error reading oidc configuration")
		}
	}
	return c.endpoints.JWKS
}

// TODO - cache this on a reasonable refresh timeline (and/or retry on failure?)
func getJWK(uri, key string) (*jose.JSONWebKey, error) {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
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
