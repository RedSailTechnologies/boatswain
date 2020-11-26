package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/twitchtv/twirp"
	"golang.org/x/oauth2"
)

var jwtKey = new(int)

// FIXME
var issuerURI = "http://localhost:4011"
var oidcConfig = &oidc.Config{
	// ClientID: "some-app",
	SkipClientIDCheck: true,
}

// WithJWT wraps an existing http hander to verify the JWTs
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

// JWTHook is the twirp server hook to validate JWTs
func JWTHook() *twirp.ServerHooks {
	hooks := &twirp.ServerHooks{}
	hooks.RequestRouted = func(ctx context.Context) (context.Context, error) {
		jwt := ctx.Value(jwtKey)
		if jwt == nil {
			return ctx, twirp.NewError(twirp.Unauthenticated, "no jwt provided")
		}

		keys := oidc.NewRemoteKeySet(ctx, issuerURI+"/.well-known/openid-configuration/jwks")
		payload, err := keys.VerifySignature(ctx, jwt.(string))
		if payload == nil {

		}

		ver := oidc.NewVerifier(issuerURI, keys, oidcConfig)
		idd, err := ver.Verify(ctx, jwt.(string))
		if idd != nil {

		}

		provider, err := oidc.NewProvider(ctx, issuerURI)
		if err != nil {
			return ctx, twirp.NewError(twirp.Unauthenticated, "could not create oidc provider")
		}

		var claims struct {
			ScopesSupported []string `json:"scopes_supported"`
			ClaimsSupported []string `json:"claims_supported"`
		}

		if err := provider.Claims(&claims); err != nil {
			// handle unmarshaling error
		}

		info, err := provider.UserInfo(ctx, &customtoken{
			token: jwt.(string),
		})
		if err != nil {

		}
		if info != nil {

		}

		verifier := provider.Verifier(oidcConfig)
		id, err := verifier.Verify(ctx, jwt.(string))
		if id != nil {

		}

		return ctx, nil
	}
	return hooks
}

type customtoken struct {
	token string
}

func (j customtoken) Token() (*oauth2.Token, error) {
	return &oauth2.Token{
		AccessToken: j.token,
	}, nil
}
