package main

import (
	"log"
	"testing"
)

func TestGetAdminToken(t *testing.T) {
	token := GetAdminToken()

	if len(token) == 0 {
		t.Fatalf("GetAdminToken failed")
	}
}

func TestGetClientId(t *testing.T) {
	r := GetClientId("", "cloudbees-client")
	if len(r) == 0 {
		t.Fatalf("Get ClientID failed")
	} else {

	}
}

func TestGetUser(t *testing.T) {
	u := GetUser("", "admin")
	if len(u.Id) == 0 {
		t.Fatalf("Get ClientID failed")
	}
}

func TestCreateNewUserAsCustomer(t *testing.T) {
	token := GetAdminToken()
	u := GetUser(token, "test.user")
	if len(u.Id) > 0 {
		log.Println(". . . found user, perform delete")
		d := DeleteUser(token, u.Username)
		if d {
			log.Println(". . . deleted user ", u.Username)
		}
	}
	r := CreateUser(token, "test", "user", "test@customer.dev")
	if len(r) == 0 {
		t.Fatalf("RegisterCustomer failed: User Not Created")
	} else {
		b := AssignUserAsClientRole(token, r, "travel_customer")
		if !b {
			t.Fatalf("RegisterCustomer failed: User Not Assigned Role")
		}
	}
}
