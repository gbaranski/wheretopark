package main

import (
	"context"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/gin-server"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"log"
	"os"
	"strconv"
	"strings"
)

func getEnv(name string) string {
	value, exists := os.LookupEnv(name)
	if !exists {
		log.Fatalf("%s is not set\n", name)
	}
	return value
}

func getEnvOr(name string, fallback string) string {
	value, exists := os.LookupEnv(name)
	if exists {
		return value
	} else {
		return fallback
	}
}

func main() {
	clientID := getEnv("CLIENT_ID")
	clientSecret := getEnv("CLIENT_SECRET")
	jwtSecret := getEnv("JWT_SECRET")
	port, err := strconv.Atoi(getEnvOr("PORT", "8080"))
	if err != nil {
		log.Fatalf("Invalid port: %v\n", err)
	}
	manager := manage.NewDefaultManager()

	// token store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	manager.MapAccessGenerate(NewJWTAccessGenerate([]byte(jwtSecret), jwt.SigningMethodHS512))

	// client store
	clientStore := store.NewClientStore()
	err = clientStore.Set(clientID, &models.Client{
		ID:     clientID,
		Secret: clientSecret,
	})
	if err != nil {
		log.Fatalf("set client error: %v\n", err)
		return
	}
	manager.MapClientStorage(clientStore)

	// Initialize the oauth2 service
	ginserver.InitServer(manager)
	ginserver.SetAllowGetAccessRequest(true)
	ginserver.SetClientInfoHandler(server.ClientFormHandler)

	g := gin.Default()

	auth := g.Group("/oauth")
	{
		//auth.GET("/authorize", ginserver.HandleAuthorizeRequest)
		auth.POST("/token", ginserver.HandleTokenRequest)
		auth.GET("/token", ginserver.HandleTokenRequest)
	}
	g.GET("/health-check", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})

	err = g.Run(":" + strconv.Itoa(port))
	if err != nil {
		log.Fatalf("server run error: %v\n", err)
		return
	}
}

// JWTAccessClaims jwt claims
type JWTAccessClaims struct {
	Scope string `json:"scope"`
	jwt.RegisteredClaims
}

// JWTAccessGenerate generate the jwt access token
type JWTAccessGenerate struct {
	SignedKey    []byte
	SignedMethod jwt.SigningMethod
}

// NewJWTAccessGenerate create to generate the jwt access token instance
func NewJWTAccessGenerate(key []byte, method jwt.SigningMethod) *JWTAccessGenerate {
	return &JWTAccessGenerate{
		SignedKey:    key,
		SignedMethod: method,
	}
}

// Token based on the UUID generated token
func (a *JWTAccessGenerate) Token(ctx context.Context, data *oauth2.GenerateBasic, isGenRefresh bool) (string, string, error) {
	log.Printf("scope=%s\n", data.TokenInfo.GetScope())

	claims := JWTAccessClaims{
		Scope: data.TokenInfo.GetScope(),
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  []string{data.Client.GetID()},
			Subject:   data.UserID,
			ExpiresAt: jwt.NewNumericDate(data.TokenInfo.GetAccessCreateAt().Add(data.TokenInfo.GetAccessExpiresIn())),
		},
	}

	token := jwt.NewWithClaims(a.SignedMethod, claims)
	access, err := token.SignedString(a.SignedKey)
	if err != nil {
		return "", "", err
	}
	refresh := ""

	if isGenRefresh {
		t := uuid.NewSHA1(uuid.Must(uuid.NewRandom()), []byte(access)).String()
		refresh = base64.URLEncoding.EncodeToString([]byte(t))
		refresh = strings.ToUpper(strings.TrimRight(refresh, "="))
	}

	return access, refresh, nil
}

func (a *JWTAccessGenerate) isEs() bool {
	return strings.HasPrefix(a.SignedMethod.Alg(), "ES")
}

func (a *JWTAccessGenerate) isRsOrPS() bool {
	isRs := strings.HasPrefix(a.SignedMethod.Alg(), "RS")
	isPs := strings.HasPrefix(a.SignedMethod.Alg(), "PS")
	return isRs || isPs
}

func (a *JWTAccessGenerate) isHs() bool {
	return strings.HasPrefix(a.SignedMethod.Alg(), "HS")
}
