package middleware

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	oidc "github.com/coreos/go-oidc"
	gin "github.com/gin-gonic/gin"
)

type Res401Struct struct {
	Status   string `json:"status" example:"FAILED"`
	HTTPCode int    `json:"httpCode" example:"401"`
	Message  string `json:"message" example:"authorisation failed"`
}

type client struct {
	ServiceClient clientRoles `json:"cloudbees-client,omitempty"`
}

type Claims struct {
	ResourceAccess client `json:"resource_access"`
	JTI            string `json:"jti"`
	Email          string `json:"email"`
	Username       string `json:"preferred_username"`
}

type clientRoles struct {
	Roles []string `json:"roles,omitempty"`
}

var RealmConfigURL string = "https://cloudbees.dev:8443/realms/cloudbees"
var clientID string = "cloudbees-client"

func SetAuthorizedJWT(required bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		isAuth := false
		isAdmin := false
		isCustomer := false
		username := ""
		strippedAccessToken := ""

		accessTokens := c.Request.Header["Authorization"]
		if len(accessTokens) > 0 {
			strippedAccessToken = strings.ReplaceAll(accessTokens[0], "Bearer ", "")
		}

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		client := &http.Client{
			Timeout:   time.Duration(6000) * time.Second,
			Transport: tr,
		}

		ctx := oidc.ClientContext(context.Background(), client)
		if len(accessTokens) > 0 {
			provider, err := oidc.NewProvider(ctx, RealmConfigURL)
			if err != nil {
				authorisationFailed("authorisation failed while getting the provider: "+err.Error(), c.Writer, c.Request)
				log.Fatalf("AUTH Failed while getting provider: %v", err)
				return
			}

			oidcConfig := &oidc.Config{
				ClientID: clientID,
			}
			verifier := provider.Verifier(oidcConfig)
			idToken, err := verifier.Verify(ctx, strippedAccessToken)
			if err != nil {
				authorisationFailed("authorisation failed while verifying the token: "+err.Error(), c.Writer, c.Request)
				log.Fatalf("AUTH Failed while verifying token: %v", err)
				return
			}

			var IDTokenClaims Claims // ID Token payload is just JSON.
			if err := idToken.Claims(&IDTokenClaims); err != nil {
				authorisationFailed("claims : "+err.Error(), c.Writer, c.Request)
				log.Fatalf("AUTH Failed while processing claims: %v", err)
				return
			}

			user_access_roles := IDTokenClaims.ResourceAccess.ServiceClient.Roles

			authRoles := strings.Split("travel_admin travel_customer", " ")
			adminRoles := strings.Split("travel_admin", " ")
			customerRoles := strings.Split("travel_customer", " ")

			isAuth = AtLeastOneRole(user_access_roles, authRoles)
			isAdmin = AtLeastOneRole(user_access_roles, adminRoles)
			isCustomer = AtLeastOneRole(user_access_roles, customerRoles)
		}

		if (required && isAuth) || !required {
			c.Set("isAuth", isAuth)
			c.Set("isAdmin", isAdmin)
			c.Set("isCustomer", isCustomer)
			c.Set("username", username)
			c.Next()
		} else {
			c.Set("isAuth", false)
			c.Set("isAdmin", false)
			c.Set("isCustomer", false)
			c.Set("username", "")
			c.JSON(http.StatusUnauthorized, gin.H{"data": nil})
		}
	}
}

func AtLeastOneRole(userRoles []string, requiredRoles []string) bool {
	var found bool = false
	for i := 0; i < len(userRoles); i++ {
		for j := 0; j < len(requiredRoles); j++ {
			if userRoles[i] == requiredRoles[j] {
				found = true
				break
			}
			if found == true {
				break
			}
		}
	}
	return found
}

func authorisationFailed(message string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	data := Res401Struct{
		Status:   "FAILED",
		HTTPCode: http.StatusUnauthorized,
		Message:  message,
	}
	res, _ := json.Marshal(data)
	w.Write(res)
}
