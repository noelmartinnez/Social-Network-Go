package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type contextKey string

const claimsKey contextKey = "claims"

var jwtKey []byte
var claveHash string
var aesKey []byte

func loadKeys() error {
	var err error

	jwtKey, err = ioutil.ReadFile("jwtKey.txt")
	if err != nil {
		return fmt.Errorf("error loading jwtKey: %v", err)
	}

	claveHashBytes, err := ioutil.ReadFile("claveHash.txt")
	if err != nil {
		return fmt.Errorf("error loading claveHash: %v", err)
	}
	claveHash = string(claveHashBytes)

	aesKey, err = ioutil.ReadFile("aesKey.txt")
	if err != nil {
		return fmt.Errorf("error loading aesKey: %v", err)
	}

	return nil
}

func Run() {
	if err := loadKeys(); err != nil {
		log.Fatalf("Failed to load keys: %v", err)
	}
	http.HandleFunc("/register", registerUser)
	http.HandleFunc("/login", loginUser)
	http.HandleFunc("/createPost", authenticate(createPost))
	http.HandleFunc("/posts", authenticate(getPosts))
	http.HandleFunc("/message", authenticate(handleMessages))
	http.HandleFunc("/users", authenticate(getUsers))
	http.HandleFunc("/roles", getRoles)
	http.HandleFunc("/user/", authenticate(deleteUser))
	http.HandleFunc("/post/", authenticate(deletePost))
	http.HandleFunc("/all-users", authenticate(getAllUsers))
	http.HandleFunc("/promoteToModerator", authenticate(promoteToModerator))
	http.HandleFunc("/backup", authenticate(adminOnly(backupDatabase)))
	http.HandleFunc("/restore", authenticate(adminOnly(restoreDatabase)))
	http.HandleFunc("/createGroup", authenticate(createGroup))
	http.HandleFunc("/groups", authenticate(getGroups))
	http.HandleFunc("/joinGroup", authenticate(joinGroup))
	http.HandleFunc("/groupPosts", authenticate(getGroupPosts))
	//http.HandleFunc("/group/", authenticate(manageGroup))
	http.HandleFunc("/userGroup/", authenticate(userGroupInfo))
	http.HandleFunc("/leaveGroup", authenticate(leaveGroup))
	http.HandleFunc("/groupMembers", authenticate(getGroupMembers))
	http.HandleFunc("/removeGroupMember", authenticate(removeGroupMember))
	http.HandleFunc("/groupDescription", authenticate(editGroupDescription))
	http.HandleFunc("/deleteGroup", authenticate(deleteGroup))

	fmt.Println("Starting server on port 443")
	err := http.ListenAndServeTLS(":443", "certificado.crt", "llave.key", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}
