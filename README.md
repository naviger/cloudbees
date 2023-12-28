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
a. From, To, User , Price Paid.
i. User should include first name, last name, email address
b. The user is allocated a seat in the train as a result of the purchase. Assume the train has only 2 sections, section A and section B and each section has 10 seats.
3. An authenticated API that shows the details of the receipt for the user
4. An authenticated API that lets an admin view all the users and seats they are allocated by the requested section
5. An authenticated API to allow an admin or the user to remove the user from the train
6. An authenticated API to allow an admin or the user to modify the user’s seat



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

# Create React App (with Typescript)

Execute the following:

```
npx create-react-app cloudbees --template typescript
```


