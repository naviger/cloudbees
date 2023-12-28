package interceptor

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	oidc "github.com/coreos/go-oidc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

func JwtInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	strippedAccessToken := ""
	isAuth := false
	isAdmin := false
	isCustomer := false
	username := ""

	md, ok := metadata.FromIncomingContext(ctx)

	newMD := md.Copy()
	if !ok {

	} else {

		rawAccessToken := md.Get("authorization")
		if len(rawAccessToken) > 0 && len(rawAccessToken[0]) > 5 {

			strippedAccessToken = strings.ReplaceAll(rawAccessToken[0], "bearer ", "")

			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			client := &http.Client{
				Timeout:   time.Duration(6000) * time.Second,
				Transport: tr,
			}
			octx := oidc.ClientContext(context.Background(), client)
			provider, err := oidc.NewProvider(octx, RealmConfigURL)

			if err != nil {
				authorisationFailed("authorisation failed while getting the provider: "+err.Error(), &newMD)
				log.Fatalf("AUTH Failed while getting provider: %v", err)
			} else {
				user_access_roles := make([]string, 0)
				oidcConfig := &oidc.Config{
					ClientID: clientID,
				}
				verifier := provider.Verifier(oidcConfig)
				idToken, err := verifier.Verify(octx, strippedAccessToken)

				//var IDTokenClaims Claims
				if err != nil {
					authorisationFailed("authorisation failed while verifying the token: "+err.Error(), &newMD)
					log.Fatalf("AUTH Failed while verifying token: %v", err)
				} else {
					var IDTokenClaims Claims // ID Token payload is just JSON.
					if err := idToken.Claims(&IDTokenClaims); err != nil {

						authorisationFailed("claims : "+err.Error(), &newMD)
						log.Fatalf("AUTH Failed while processing claims: %v", err)
					}
					user_access_roles = IDTokenClaims.ResourceAccess.ServiceClient.Roles
					username = IDTokenClaims.Username
				}

				authRoles := strings.Split("travel_admin travel_customer", " ")
				adminRoles := strings.Split("travel_admin", " ")
				customerRoles := strings.Split("travel_customer", " ")

				isAuth = AtLeastOneRole(user_access_roles, authRoles)
				isAdmin = AtLeastOneRole(user_access_roles, adminRoles)
				isCustomer = AtLeastOneRole(user_access_roles, customerRoles)
			}

		}
	}

	if isAuth {
		newMD.Append("isAuth", strconv.FormatBool(isAuth))
		newMD.Append("isAdmin", strconv.FormatBool(isAdmin))
		newMD.Append("isCustomer", strconv.FormatBool(isCustomer))
		newMD.Append("user", username)
		newMD.Append("message", "authenticated")
	} else {
		newMD.Append("isAuth", "false")
		newMD.Append("isAdmin", "false")
		newMD.Append("isCustomer", "false")
		newMD.Append("user", "")
		newMD.Append("message", "not authenticated")
	}

	newCtx := metadata.NewIncomingContext(ctx, newMD)
	return handler(newCtx, req)
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

func authorisationFailed(message string, newMD *metadata.MD) {
	newMD.Append("isAuth", "false")
	newMD.Append("isAdmin", "false")
	newMD.Append("isCustomer", "false")
	newMD.Append("user", "")
	newMD.Append("message", message)
}
