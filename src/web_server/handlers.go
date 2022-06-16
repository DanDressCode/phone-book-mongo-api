package web_server

import (
	"encoding/json"
	_ "fmt"
	"github.com/vertex/phoneBook/mongo_db"
	_ "golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
)

type userAuth struct {
	UserAd   string `bson:"userad"`
	Password string `bson:"password"`
}
type userPasswords struct {
	username string
}

var (
	username = "Admin"
	password = "Admin"
)

/*
func AddBasicAuth(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	if ok && verifyUserPass(user, pass) {
		fmt.Fprintf(w, "You get to see the secret\n")
	} else {
		w.Header().Set("WWW-Authenticate", `Basic realm="api"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func verifyUserPass(username, password string) bool {

	wantPass, hasUser := userPasswords[username]
	if !hasUser {
		return false
	}
	if cmperr := bcrypt.CompareHashAndPassword(wantPass, []byte(password)); cmperr == nil {
		return true
	}
	return false
}*/

/*func AddAuth(w http.ResponseWriter, r *http.Request) string {

	trace := userAuth{
		"Admin",
		"Admin",
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Метод запрещен!", 405)
		return ""
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return ""
	}
	var userAu userAuth
	err = json.Unmarshal(body, &userAu)
	if err != nil {
		http.Error(w, "can't decode user", http.StatusBadRequest)
		return ""
	}

	if userAu == trace {
		fmt.Fprintf(w, "Authenticate")
		return ""
	} else {
		w.Header().Set("WWW-Authenticate", `Basic realm="api"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
	return ""
}*/

func GetAllUsers(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Метод запрещен!", 405)
		return
	}

	users, err := mongo_db.GetAll()
	if err != nil {
		http.Error(w, error.Error(err), 500)
		return
	}

	jsonResp, err := json.Marshal(users)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Метод запрещен!", 405)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	var user mongo_db.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "can't decode user", http.StatusBadRequest)
		return
	}
	//
	err = mongo_db.CreateUser(&user)
	if err != nil {
		http.Error(w, "mongo db is down!", 500)
		return
	}

}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Метод запрещен!", 405)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	var user mongo_db.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "can't decode user", http.StatusBadRequest)
		return
	}

	err = mongo_db.DeleteUser(&user)
	if err != nil {
		http.Error(w, "mongo db is down!", 500)
		return
	}

}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPut {
		w.Header().Set("Allow", http.MethodPut)
		http.Error(w, "Метод запрещен!", 405)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	var user mongo_db.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "can't decode user", http.StatusBadRequest)
		return
	}

	err = mongo_db.UpdateUser(&user)
	if err != nil {
		http.Error(w, "mongo db is down!", 500)
		return
	}

}

func findUser(w http.ResponseWriter, r *http.Request) (bool, error) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Метод запрещен!", 405)
		return false, nil

	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return false, nil
	}
	var user mongo_db.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "can't decode user", http.StatusBadRequest)
		return false, nil
	}
	r.SetBasicAuth(username, password)

	return false, nil

}
