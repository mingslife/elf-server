// see: github.com/appleboy/gin-jwt
package component

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mingslife/bone"

	"elf-server/pkg/conf"
)

type MapClaims map[string]any

var (
	// ErrMissingSecretKey indicates Secret key is required
	ErrMissingSecretKey = errors.New("secret key is required")

	// ErrForbidden when HTTP status 403 is given
	ErrForbidden = errors.New("you don't have permission to access this resource")

	// ErrMissingAuthenticatorFunc indicates Authenticator is required
	ErrMissingAuthenticatorFunc = errors.New("ginJWTMiddleware.Authenticator func is undefined")

	// ErrMissingLoginValues indicates a user tried to authenticate without username or password
	ErrMissingLoginValues = errors.New("missing Username or Password")

	// ErrFailedAuthentication indicates authentication failed, could be faulty username or password
	ErrFailedAuthentication = errors.New("incorrect Username or Password")

	// ErrFailedTokenCreation indicates JWT Token failed to create, reason unknown
	ErrFailedTokenCreation = errors.New("failed to create JWT Token")

	// ErrExpiredToken indicates JWT token has expired. Can't refresh.
	ErrExpiredToken = errors.New("token is expired")

	// ErrEmptyAuthHeader can be thrown if authing with a HTTP header, the Auth header needs to be set
	ErrEmptyAuthHeader = errors.New("auth header is empty")

	// ErrMissingExpField missing exp field in token
	ErrMissingExpField = errors.New("missing exp field")

	// ErrWrongFormatOfExp field must be float64 format
	ErrWrongFormatOfExp = errors.New("exp must be float64 format")

	// ErrInvalidAuthHeader indicates auth header is invalid, could for example have the wrong Realm name
	ErrInvalidAuthHeader = errors.New("auth header is invalid")

	// ErrEmptyQueryToken can be thrown if authing with URL Query, the query token variable is empty
	ErrEmptyQueryToken = errors.New("query token is empty")

	// ErrEmptyCookieToken can be thrown if authing with a cookie, the token cookie is empty
	ErrEmptyCookieToken = errors.New("cookie token is empty")

	// ErrEmptyParamToken can be thrown if authing with parameter in path, the parameter in path is empty
	ErrEmptyParamToken = errors.New("parameter token is empty")

	// ErrInvalidSigningAlgorithm indicates signing algorithm is invalid, needs to be HS256, HS384, HS512, RS256, RS384 or RS512
	ErrInvalidSigningAlgorithm = errors.New("invalid signing algorithm")

	// ErrNoPrivKeyFile indicates that the given private key is unreadable
	ErrNoPrivKeyFile = errors.New("private key file unreadable")

	// ErrNoPubKeyFile indicates that the given public key is unreadable
	ErrNoPubKeyFile = errors.New("public key file unreadable")

	// ErrInvalidPrivKey indicates that the given private key is invalid
	ErrInvalidPrivKey = errors.New("private key invalid")

	// ErrInvalidPubKey indicates the the given public key is invalid
	ErrInvalidPubKey = errors.New("public key invalid")

	// IdentityKey default identity key
	IdentityKey = "identity"
)

type Jwt struct {
	Router *bone.Router `inject:"application.router"`

	// GinJWTMiddleware part:

	// Realm name to display to the user. Required.
	Realm string

	// signing algorithm - possible values are HS256, HS384, HS512, RS256, RS384 or RS512
	// Optional, default is HS256.
	SigningAlgorithm string

	// Secret key used for signing. Required.
	Key []byte

	// Duration that a jwt token is valid. Optional, defaults to one hour.
	Timeout time.Duration

	// This field allows clients to refresh their token until MaxRefresh has passed.
	// Note that clients can refresh their token in the last moment of MaxRefresh.
	// This means that the maximum validity timespan for a token is TokenTime + MaxRefresh.
	// Optional, defaults to 0 meaning not refreshable.
	MaxRefresh time.Duration

	// Callback function that should perform the authentication of the user based on login info.
	// Must return user data as user identifier, it will be stored in Claim Array. Required.
	// Check error (e) to determine the appropriate error message.
	Authenticator func(w http.ResponseWriter) (any, error)

	// Callback function that should perform the authorization of the authenticated user. Called
	// only after an authentication success. Must return true on success, false on failure.
	// Optional, default to success.
	Authorizator func(data any, w http.ResponseWriter) bool

	// Callback function that will be called during login.
	// Using this function it is possible to add additional payload data to the webtoken.
	// The data is then made available during requests via c.Get("JWT_PAYLOAD").
	// Note that the payload is not encrypted.
	// The attributes mentioned on jwt.io can't be used as keys for the map.
	// Optional, by default no additional data will be set.
	PayloadFunc func(data any) MapClaims

	// User can define own Unauthorized func.
	Unauthorized func(http.ResponseWriter, int, string)

	// User can define own LoginResponse func.
	LoginResponse func(http.ResponseWriter, int, string, time.Time)

	// User can define own LogoutResponse func.
	LogoutResponse func(http.ResponseWriter, int)

	// User can define own RefreshResponse func.
	RefreshResponse func(http.ResponseWriter, int, string, time.Time)

	// Set the identity handler function
	IdentityHandler func(*http.Request, http.ResponseWriter) any

	// Set the identity key
	IdentityKey string

	// TokenLookup is a string in the form of "<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>"
	// - "query:<name>"
	// - "cookie:<name>"
	TokenLookup string

	// TokenHeadName is a string in the header. Default value is "Bearer"
	TokenHeadName string

	// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
	TimeFunc func() time.Time

	// HTTP Status messages for when something in the JWT middleware fails.
	// Check error (e) to determine the appropriate error message.
	HTTPStatusMessageFunc func(e error) string

	// Private key file for asymmetric algorithms
	PrivKeyFile string

	// Public key file for asymmetric algorithms
	PubKeyFile string

	// Private key
	privKey *rsa.PrivateKey

	// Public key
	pubKey *rsa.PublicKey

	// Optionally return the token as a cookie
	SendCookie bool

	// Duration that a cookie is valid. Optional, by default equals to Timeout value.
	CookieMaxAge time.Duration

	// Allow insecure cookies for development over http
	SecureCookie bool

	// Allow cookies to be accessed client side for development
	CookieHTTPOnly bool

	// Allow cookie domain change for development
	CookieDomain string

	// Disable abort() of context.
	DisabledAbort bool

	// CookieName allow cookie name change for development
	CookieName string
}

func (*Jwt) Name() string {
	return "component.jwt"
}

func (mw *Jwt) Init() error {
	return nil
}

func (mw *Jwt) Register() error {
	cfg := conf.GetConfig()
	mw.Realm = cfg.Name
	mw.Key = []byte(cfg.JwtKey)
	mw.Timeout = time.Duration(cfg.JwtTimeout) * time.Hour
	mw.MaxRefresh = time.Duration(cfg.JwtMaxRefresh) * time.Hour
	mw.PayloadFunc = func(data any) MapClaims {
		if dataMap, ok := data.(map[string]any); ok {
			return dataMap
		}
		return map[string]any{}
	}
	mw.MiddlewareInit()
	mw.Router.Use(mw.Middleware)
	return nil
}

func (*Jwt) Unregister() error {
	return nil
}

func (mw *Jwt) Middleware(next http.Handler) http.Handler {
	excludePaths := [...]string{
		"/api/v1/auth/login",
		"/api/v1/auth/register",
		"/api/v1/auth/settings",
		// "/api/v1/upload/file",
		// "/api/v1/upload/image",
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.RequestURI
		if !strings.HasPrefix(path, "/api/") {
			next.ServeHTTP(w, r)
			return
		} else {
			for _, excludePath := range excludePaths {
				if path == excludePath {
					next.ServeHTTP(w, r)
					return
				}
			}
		}

		if ok := mw.middlewareImpl(r, w); ok {
			next.ServeHTTP(w, r)
		}
	})
}

var _ bone.Component = (*Jwt)(nil)

// Methods from *GinJWTMiddleware

func (mw *Jwt) readKeys() error {
	err := mw.privateKey()
	if err != nil {
		return err
	}
	err = mw.publicKey()
	if err != nil {
		return err
	}
	return nil
}

func (mw *Jwt) privateKey() error {
	keyData, err := ioutil.ReadFile(mw.PrivKeyFile)
	if err != nil {
		return ErrNoPrivKeyFile
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		return ErrInvalidPrivKey
	}
	mw.privKey = key
	return nil
}

func (mw *Jwt) publicKey() error {
	keyData, err := ioutil.ReadFile(mw.PubKeyFile)
	if err != nil {
		return ErrNoPubKeyFile
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return ErrInvalidPubKey
	}
	mw.pubKey = key
	return nil
}

func (mw *Jwt) usingPublicKeyAlgo() bool {
	switch mw.SigningAlgorithm {
	case "RS256", "RS512", "RS384":
		return true
	}
	return false
}

// MiddlewareInit initialize jwt configs.
func (mw *Jwt) MiddlewareInit() error {
	if mw.TokenLookup == "" {
		mw.TokenLookup = "header:Authorization"
	}

	if mw.SigningAlgorithm == "" {
		mw.SigningAlgorithm = "HS256"
	}

	if mw.Timeout == 0 {
		mw.Timeout = time.Hour
	}

	if mw.TimeFunc == nil {
		mw.TimeFunc = time.Now
	}

	mw.TokenHeadName = strings.TrimSpace(mw.TokenHeadName)
	if len(mw.TokenHeadName) == 0 {
		mw.TokenHeadName = "Bearer"
	}

	if mw.Authorizator == nil {
		mw.Authorizator = func(data any, w http.ResponseWriter) bool {
			return true
		}
	}

	if mw.Unauthorized == nil {
		mw.Unauthorized = func(w http.ResponseWriter, code int, message string) {
			mw.writeJSON(w, code, map[string]any{
				"code":    code,
				"message": message,
			})
		}
	}

	if mw.LoginResponse == nil {
		mw.LoginResponse = func(w http.ResponseWriter, code int, token string, expire time.Time) {
			mw.writeJSON(w, http.StatusOK, map[string]any{
				"code":   http.StatusOK,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		}
	}

	if mw.LogoutResponse == nil {
		mw.LogoutResponse = func(w http.ResponseWriter, code int) {
			mw.writeJSON(w, http.StatusOK, map[string]any{
				"code": http.StatusOK,
			})
		}
	}

	if mw.RefreshResponse == nil {
		mw.RefreshResponse = func(w http.ResponseWriter, code int, token string, expire time.Time) {
			mw.writeJSON(w, http.StatusOK, map[string]any{
				"code":   http.StatusOK,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		}
	}

	if mw.IdentityKey == "" {
		mw.IdentityKey = IdentityKey
	}

	if mw.IdentityHandler == nil {
		mw.IdentityHandler = func(r *http.Request, w http.ResponseWriter) any {
			claims := ExtractClaims(r)
			return claims[mw.IdentityKey]
		}
	}

	if mw.HTTPStatusMessageFunc == nil {
		mw.HTTPStatusMessageFunc = func(e error) string {
			return e.Error()
		}
	}

	if mw.Realm == "" {
		mw.Realm = "jwt"
	}

	if mw.CookieMaxAge == 0 {
		mw.CookieMaxAge = mw.Timeout
	}

	if mw.CookieName == "" {
		mw.CookieName = "jwt"
	}

	if mw.usingPublicKeyAlgo() {
		return mw.readKeys()
	}

	if mw.Key == nil {
		return ErrMissingSecretKey
	}

	return nil
}

func (mw *Jwt) middlewareImpl(r *http.Request, w http.ResponseWriter) bool {
	claims, err := mw.GetClaimsFromJWT(r)
	if err != nil {
		return mw.unauthorized(w, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(err))
	}

	if claims["exp"] == nil {
		return mw.unauthorized(w, http.StatusBadRequest, mw.HTTPStatusMessageFunc(ErrMissingExpField))
	}

	if _, ok := claims["exp"].(float64); !ok {
		return mw.unauthorized(w, http.StatusBadRequest, mw.HTTPStatusMessageFunc(ErrWrongFormatOfExp))
	}

	if int64(claims["exp"].(float64)) < mw.TimeFunc().Unix() {
		return mw.unauthorized(w, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(ErrExpiredToken))
	}

	ctx := context.WithValue(r.Context(), "JWT_PAYLOAD", claims)
	*r = *r.Clone(ctx)
	identity := mw.IdentityHandler(r, w)

	if identity != nil {
		ctx := context.WithValue(r.Context(), mw.IdentityKey, identity)
		*r = *r.Clone(ctx)
	}

	if !mw.Authorizator(identity, w) {
		return mw.unauthorized(w, http.StatusForbidden, mw.HTTPStatusMessageFunc(ErrForbidden))
	}

	return true
}

// GetClaimsFromJWT get claims from JWT token
func (mw *Jwt) GetClaimsFromJWT(r *http.Request) (MapClaims, error) {
	token, err := mw.ParseToken(r)

	if err != nil {
		return nil, err
	}

	claims := MapClaims{}
	for key, value := range token.Claims.(jwt.MapClaims) {
		claims[key] = value
	}

	return claims, nil
}

// LoginHandler can be used by clients to get a jwt token.
// Payload needs to be json in the form of {"username": "USERNAME", "password": "PASSWORD"}.
// Reply will be of the form {"token": "TOKEN"}.
func (mw *Jwt) LoginHandler(r *http.Request, w http.ResponseWriter) {
	if mw.Authenticator == nil {
		mw.unauthorized(w, http.StatusInternalServerError, mw.HTTPStatusMessageFunc(ErrMissingAuthenticatorFunc))
		return
	}

	data, err := mw.Authenticator(w)

	if err != nil {
		mw.unauthorized(w, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(err))
		return
	}

	// Create the token
	token := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)

	if mw.PayloadFunc != nil {
		for key, value := range mw.PayloadFunc(data) {
			claims[key] = value
		}
	}

	expire := mw.TimeFunc().Add(mw.Timeout)
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = mw.TimeFunc().Unix()
	tokenString, err := mw.signedString(token)

	if err != nil {
		mw.unauthorized(w, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(ErrFailedTokenCreation))
		return
	}

	// set cookie
	if mw.SendCookie {
		expireCookie := mw.TimeFunc().Add(mw.CookieMaxAge)
		maxage := int(expireCookie.Unix() - mw.TimeFunc().Unix())

		mw.setCookie(
			w,
			mw.CookieName,
			tokenString,
			maxage,
			"/",
			mw.CookieDomain,
			mw.SecureCookie,
			mw.CookieHTTPOnly,
		)
	}

	mw.LoginResponse(w, http.StatusOK, tokenString, expire)
}

// LogoutHandler can be used by clients to remove the jwt cookie (if set)
func (mw *Jwt) LogoutHandler(r *http.Request, w http.ResponseWriter) {
	// delete auth cookie
	mw.setCookie(
		w,
		mw.CookieName,
		"",
		-1,
		"/",
		mw.CookieDomain,
		mw.SecureCookie,
		mw.CookieHTTPOnly,
	)

	mw.LogoutResponse(w, http.StatusOK)
}

func (mw *Jwt) signedString(token *jwt.Token) (string, error) {
	var tokenString string
	var err error
	if mw.usingPublicKeyAlgo() {
		tokenString, err = token.SignedString(mw.privKey)
	} else {
		tokenString, err = token.SignedString(mw.Key)
	}
	return tokenString, err
}

// RefreshHandler can be used to refresh a token. The token still needs to be valid on refresh.
// Shall be put under an endpoint that is using the GinJWTMiddleware.
// Reply will be of the form {"token": "TOKEN"}.
func (mw *Jwt) RefreshHandler(r *http.Request, w http.ResponseWriter) {
	tokenString, expire, err := mw.RefreshToken(r, w)
	if err != nil {
		mw.unauthorized(w, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(err))
		return
	}

	mw.RefreshResponse(w, http.StatusOK, tokenString, expire)
}

// RefreshToken refresh token and check if token is expired
func (mw *Jwt) RefreshToken(r *http.Request, w http.ResponseWriter) (string, time.Time, error) {
	claims, err := mw.CheckIfTokenExpire(r)
	if err != nil {
		return "", time.Now(), err
	}

	// Create the token
	newToken := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	newClaims := newToken.Claims.(jwt.MapClaims)

	for key := range claims {
		newClaims[key] = claims[key]
	}

	expire := mw.TimeFunc().Add(mw.Timeout)
	newClaims["exp"] = expire.Unix()
	newClaims["orig_iat"] = mw.TimeFunc().Unix()
	tokenString, err := mw.signedString(newToken)

	if err != nil {
		return "", time.Now(), err
	}

	// set cookie
	if mw.SendCookie {
		expireCookie := mw.TimeFunc().Add(mw.CookieMaxAge)
		maxage := int(expireCookie.Unix() - time.Now().Unix())

		mw.setCookie(
			w,
			mw.CookieName,
			tokenString,
			maxage,
			"/",
			mw.CookieDomain,
			mw.SecureCookie,
			mw.CookieHTTPOnly,
		)
	}

	return tokenString, expire, nil
}

// CheckIfTokenExpire check if token expire
func (mw *Jwt) CheckIfTokenExpire(r *http.Request) (jwt.MapClaims, error) {
	token, err := mw.ParseToken(r)

	if err != nil {
		// If we receive an error, and the error is anything other than a single
		// ValidationErrorExpired, we want to return the error.
		// If the error is just ValidationErrorExpired, we want to continue, as we can still
		// refresh the token if it's within the MaxRefresh time.
		// (see https://github.com/appleboy/gin-jwt/issues/176)
		validationErr, ok := err.(*jwt.ValidationError)
		if !ok || validationErr.Errors != jwt.ValidationErrorExpired {
			return nil, err
		}
	}

	claims := token.Claims.(jwt.MapClaims)

	origIat := int64(claims["orig_iat"].(float64))

	if origIat < mw.TimeFunc().Add(-mw.MaxRefresh).Unix() {
		return nil, ErrExpiredToken
	}

	return claims, nil
}

// TokenGenerator method that clients can use to get a jwt token.
func (mw *Jwt) TokenGenerator(data any) (string, time.Time, error) {
	token := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)

	if mw.PayloadFunc != nil {
		for key, value := range mw.PayloadFunc(data) {
			claims[key] = value
		}
	}

	expire := mw.TimeFunc().UTC().Add(mw.Timeout)
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = mw.TimeFunc().Unix()
	tokenString, err := mw.signedString(token)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expire, nil
}

func (mw *Jwt) jwtFromHeader(r *http.Request, key string) (string, error) {
	authHeader := r.Header.Get(key)

	if authHeader == "" {
		return "", ErrEmptyAuthHeader
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == mw.TokenHeadName) {
		return "", ErrInvalidAuthHeader
	}

	return parts[1], nil
}

func (mw *Jwt) jwtFromQuery(r *http.Request, key string) (string, error) {
	token := r.URL.Query().Get(key)

	if token == "" {
		return "", ErrEmptyQueryToken
	}

	return token, nil
}

func (mw *Jwt) jwtFromCookie(r *http.Request, key string) (string, error) {
	cookie, err := r.Cookie(key)

	if err != nil {
		return "", ErrEmptyCookieToken
	}

	return cookie.Value, nil
}

// ParseToken parse jwt token from gin context
func (mw *Jwt) ParseToken(r *http.Request) (*jwt.Token, error) {
	var token string
	var err error

	methods := strings.Split(mw.TokenLookup, ",")
	for _, method := range methods {
		if len(token) > 0 {
			break
		}
		parts := strings.Split(strings.TrimSpace(method), ":")
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		switch k {
		case "header":
			token, err = mw.jwtFromHeader(r, v)
		case "query":
			token, err = mw.jwtFromQuery(r, v)
		case "cookie":
			token, err = mw.jwtFromCookie(r, v)
		}
	}

	if err != nil {
		return nil, err
	}

	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if jwt.GetSigningMethod(mw.SigningAlgorithm) != t.Method {
			return nil, ErrInvalidSigningAlgorithm
		}
		if mw.usingPublicKeyAlgo() {
			return mw.pubKey, nil
		}

		// save token string if vaild
		ctx := context.WithValue(r.Context(), "JWT_TOKEN", token)
		*r = *r.Clone(ctx)

		return mw.Key, nil
	})
}

// ParseTokenString parse jwt token string
func (mw *Jwt) ParseTokenString(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if jwt.GetSigningMethod(mw.SigningAlgorithm) != t.Method {
			return nil, ErrInvalidSigningAlgorithm
		}
		if mw.usingPublicKeyAlgo() {
			return mw.pubKey, nil
		}

		return mw.Key, nil
	})
}

func (mw *Jwt) unauthorized(w http.ResponseWriter, code int, message string) (shouldContinue bool) {
	w.Header().Set("WWW-Authenticate", "JWT realm="+mw.Realm)
	mw.Unauthorized(w, code, message)
	return mw.DisabledAbort
}

func (mw *Jwt) writeJSON(w http.ResponseWriter, code int, payload map[string]any) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(payload)
}

func (mw *Jwt) setCookie(w http.ResponseWriter, name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: 0,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

// ExtractClaims help to extract the JWT claims
func ExtractClaims(r *http.Request) MapClaims {
	claims := r.Context().Value("JWT_PAYLOAD")
	if claims == nil {
		return make(MapClaims)
	}

	return claims.(MapClaims)
}

// ExtractClaimsFromToken help to extract the JWT claims from token
func ExtractClaimsFromToken(token *jwt.Token) MapClaims {
	if token == nil {
		return make(MapClaims)
	}

	claims := MapClaims{}
	for key, value := range token.Claims.(jwt.MapClaims) {
		claims[key] = value
	}

	return claims
}

// GetToken help to get the JWT token string
func GetToken(r *http.Request) string {
	token := r.Context().Value("JWT_TOKEN")
	if token == nil {
		return ""
	}

	return token.(string)
}
