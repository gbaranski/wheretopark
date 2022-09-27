package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/gin-gonic/gin"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type config struct {
	ApplicationClientID     string `env:"APPLICATION_CLIENT_ID"`
	ApplicationClientSecret string `env:"APPLICATION_CLIENT_SECRET"`
	ProviderClientID        string `env:"PROVIDER_CLIENT_ID"`
	ProviderClientSecret    string `env:"PROVIDER_CLIENT_SECRET"`
	JwtSecret               string `env:"JWT_SECRET"`
	Port                    uint16 `env:"PORT" envDefault:"8080"`
}

var (
	ApplicationAccessTypes = []string{
		"read:metadata",
		"read:state",
		"read:status",
	}
	ProviderAccessTypes = []string{
		"read:metadata",
		"read:state",
		"read:status",
		"write:metadata",
		"write:state",
	}
)

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func main() {
	config := config{}
	if err := env.Parse(&config); err != nil {
		fmt.Printf("%+v\n", err)
	}
	log.Printf("config=%+v\n", config)
	log.Printf("port=%d\n", config.Port)
	manager := manage.NewDefaultManager()
	// TODO: remove the token store
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	manager.MapAccessGenerate(NewJWTAccessGenerate([]byte(config.JwtSecret), jwt.SigningMethodHS512))

	// client store
	clientStore := store.NewClientStore()
	clients := [][2]string{
		{config.ApplicationClientID, config.ApplicationClientSecret},
		{config.ProviderClientID, config.ProviderClientSecret},
	}
	for _, client := range clients {
		id, secret := client[0], client[1]
		fmt.Printf("adding client with ID: %s and secret: %s\n", id, secret)
		err := clientStore.Set(id, &models.Client{
			ID:     id,
			Secret: secret,
		})
		if err != nil {
			log.Fatalf("set client error: %v\n", err)
			return
		}
	}
	manager.MapClientStorage(clientStore)
	srv := server.NewDefaultServer(manager)
	srv.SetAllowedGrantType(oauth2.ClientCredentials)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})
	srv.SetClientScopeHandler(func(tgr *oauth2.TokenGenerateRequest) (bool, error) {
		var allowedAccessTypes []string
		if tgr.ClientID == config.ApplicationClientID {
			allowedAccessTypes = ApplicationAccessTypes
		} else if tgr.ClientID == config.ProviderClientID {
			allowedAccessTypes = ProviderAccessTypes
		} else {
			return false, errors.New("using not configured Client ID")
		}

		if tgr.Scope == "" {
			return true, nil
		}

		accessTypes := strings.Split(tgr.Scope, " ")
		for _, accessType := range accessTypes {
			if !contains(allowedAccessTypes, accessType) {
				return false, errors.New(fmt.Sprintf("using access type %s is not allowed for application", accessType))
			}
		}
		return true, nil
	})

	r := gin.Default()

	oauth := r.Group("/oauth")
	{
		oauth.Any("/token", func(c *gin.Context) {
			err := srv.HandleTokenRequest(c.Writer, c.Request)
			if err != nil {
				fmt.Printf("token error: %v\n", err)
				c.String(http.StatusBadRequest, err.Error())
			}
		})
	}

	r.Any("/health-check", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World")
	})

	err := r.Run(":" + strconv.Itoa(int(config.Port)))
	if err != nil {
		panic(err)
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
func (a *JWTAccessGenerate) Token(_ context.Context, data *oauth2.GenerateBasic, isGenRefresh bool) (string, string, error) {
	log.Printf("scope=%s\n", data.TokenInfo.GetScope())

	claims := JWTAccessClaims{
		Scope: data.TokenInfo.GetScope(),
		RegisteredClaims: jwt.RegisteredClaims{
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
