# Cloudbees - Technical Assessment
Cloudbees test implementation
Staff Software Engineer  - Team Lead
 
## Requirements:
1. Code must be published in Github with a link we can access (use public repo).
2. Code must compile with some effort on unit tests, doesn’t have to be 100%, but it shouldn’t be 0%.
3. Please code this with Golang and gRPC
4. No persistence layer is required, just store the data in the current session/in memory.
5. The results can be in the console output from your grpc-server and grpc-client
6. Depending on the level of authentication, take different actions

## App to be coded
__Note:__ All APIs referenced are gRPC APIs, not REST ones.
I want to board a train from London to France. The train ticket will cost $20, regardless of section or seat.
1. Authenticated APIs should be able to parse a JWT, formatted as if from an OAuth2 server, from the metadata to authenticate a request. No signature validation is required.
2. Create a public API where you can submit a purchase for a ticket. Details included in the receipt are:
	<ol type="a"><li>From, To, User , Price Paid.
		<ol type="i"><li>User should include first name, last name, email address</li></ol>
	</li>
	<li>The user is allocated a seat in the train as a result of the purchase. Assume the train has only 2 sections, section A and section B and each section has 10 seats.</li>
	</ol>
3. An authenticated API that shows the details of the receipt for the user
4. An authenticated API that lets an admin view all the users and seats they are allocated by the requested section
5. An authenticated API to allow an admin or the user to remove the user from the train
6. An authenticated API to allow an admin or the user to modify the user’s seat

# Solution Overview
The solution is composed of A GoLang-based service that provides the GRPC/Proto service. This is fronted by a golang API that exposes the GRPC/Proto service as a rest service via Go/Gin. On the front end is a ReactJS application that consumes the Rest API. It uses OIDC for authorization, with client roles of __travel_admin__ and __travel_customer__. The backing database is an in-memory database hosted in the same process as the Go GRPC server. It contains three tables: 

1. Train - a repository for the list of trains. Trains are established as one for each day for the next thirty days from the start of the server. Train names are the date in yyyMMdd format.
2. Seat - The seats that arew available for the given trains. With 20 seats per train, this is a table of 600 records. 
3. Receipt - The receipts for any transactions that are performed. This is an empty table to start.

The UI uses Axios and react-oidc-context, where the react-oidc-context interacts with KeyCloak to receive the JWT. The JWT, when available, is attached to the controller  through wich all calls are routed, ensuring the token is transmitted on all authorized calls.

A seperate __deploy__ project handles the creation of the infrastructure needed.

# Assumptions
1. Authentication is provided by OAuth with OIDC for authorization
1. Unauthenticated user receives details similar to receipt on successfull reservation. Admin or authenticated Customer can pull receipt.
1. Receipt is a receipt, it is the purchase confirmation and not a boarding pass. It is unchanging, and seat may change via future actions (change/Cancel).
2. France is assumed to be Paris, to normalize initial locationa nd destination as cities.
3. UserId/PassengerId is defined as "firstName.lastname"
4. Seats on trains can be reserved up to 30 days in advance. 
5. 

# Design Decisions
1. Utilize Hashicorps go-memdb for in memory database. 
2. Utilize KeyCloak with PostgreSQL for the authentication infrastructure.
3. Authentication and Authorization is provided by KeyCloak, with a JWT token provided with client (not realm) roles.  
4. JWT details are extracted at both the API and the Service layers. Token is passed from API to GRPC Service in context metadata. coreos-oidc library is used for oidc.
5. All layers will use SSL/TLS to encrypt data in transmission.
6. Utilize Fluent UI 9 for React JS UI Components where necessary
7. 

# ToDos
1. Login from tests to KeyCloak to get JWT - currently JWT is hardcoaded. Due to the timely nature of the JWT, the JWT must be inserted into the test files.
2. API-level testing
3. UI-level testing
4. Currently Proto output is copied from server to API to provide client. look for a solution to share.
5. All certs are currently shared and saved to each component service. Need to create scripting to manage properly with individial certs.
6. Deploy the Server to a container
7. Deploy the API to a container
8. Deploy the React UI to a container
9. ~~Logout on the UI should redirect to the home page~~
10. Ensure user can not double book
11. Modify Change Seat in service to allow retreival of receipt on seat change.
12. Modify receipt client-side to allow Admin to get receipt  

# Certificate Setup
## Create CA
	1. openssl req -x509 -nodes -new -sha256 -days 1024 -newkey rsa:2048 -keyout RootCA.key -out RootCA.pem -subj "/C=US/CN=cloudbees-CA"
	2. openssl x509 -outform pem -in RootCA.pem -out RootCA.crt

## Create Server Certificate

1. First, create a file domains.ext that lists all your local domains:

```
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names
[alt_names]
DNS.1 = localhost
DNS.2 = cloudbees.dev
DNS.3 = cloudbees.dev
```

2. openssl req -new -nodes -newkey rsa:2048 -keyout localhost.key -out localhost.csr -subj "/C=US/ST=FL/L=Miami/O=Example-Certificates/CN=localhost.local"
3. openssl x509 -req -sha256 -days 1024 -in localhost.csr -CA RootCA.pem -CAkey RootCA.key -CAcreateserial -extfile domains.ext -out localhost.crt

## Upload the CA Cert to Keychain
1. Open the Keychain App
2. Select the System from System Keychains in side bar
3. Drag and drop the RootCA.crt file to the System keychains
4. Change the trust to Always Trust

## Configure Server
1. Add localhost.crt
2. Add localhost.key

# Configure Key Cloak
1. Create Application Client
	- Tab 1: Client Type: OpenID Connect
	- Tab 1: Client ID: cloudbees_client
	- Tab 1: Name: CloudbeesClient
	- Tab 2: Client authentication: On
	- Tab 2: Authorization: On~~
	- Tab 2: Authentication flow: Standard flow, Direct Access Grants
	- Tab 3: Root URL: https://cloudbees.dev:3443
	- Tab 3: Home URL: http://cloudbees.dev:3443
	- Tab 3: Valid RedirectURLs:  http://cloudbees.dev:3443, http://cloudbees.dev:3444
	- Tab 3: Web origins: http://cloudbees.dev:3443
4. Add Client Roles
	- travel_admin
	- travel_user
5. Add Groups
	- admins => travel_admin
	- customers => stravel_customer
6. Add Users (admin)
	- Add username (email Address)
	- Add email address
	- add first name and last name
	- Assign Roles
7. Create API Client
	- Tab 1: Client type: OpenID Connect
	- Tab 2: Client ID: cloudbees_service
	- Tab 3: Name: cloudbees Service Account
	- Tab 2: Client authentication: On
	- Tab 2: Authorization: On
	- Tab 3: Root URL: http://cloudbees.dev:3443
	- Tab 3: Home URL: http://cloudbees.dev:3443
	- Tab 3: Valid RedirectURLs:  http://cloudbees.dev:3443, http://cloudbees.dev:3444
	- Tab 3: Web origins:  http://cloudbees.dev:3443
7. Create Client Scope
	- Name: cloudbees-common
	- Add a mapper with email, full name, client roles
	- create protocol mapper, type: Audience, name: cloudbees-common-audience, Include client audience: cloudbees-client, Add to access token: true
8. Add the newly created client scope " cloudbees-common-audience" to "cloudbees-service" client's client roles

# Common Terminal Calls
1. Start Container environments
```
docker compose up
```
1. Proto output:
```
protoc train.proto --go_out=. --go-grpc_out=.
```
2. Start Go Service - from cloudbees-service directory
```
go run main.go server.go helpers.go schema.go seedNames.go keycloak.go
```
3. Run Go Service Tests - from cloudbees-service directory
```
go test -v
```

4. Run UI
```
yarn start
```