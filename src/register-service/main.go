package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/dhamith93/sample-microservice-app/src/register-service/pkg/user"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type Msg struct {
	Success bool
	Msg     string
}

func main() {
	var port string

	flag.StringVar(&port, "port", "", "Port number :port")
	flag.Parse()

	db, err := Connect("root", "1234", "localhost", "dashboard_user")

	if err != nil {
		log.Fatalf("error with connecting to dashboard_user db %v\n", err)
	}

	db.AutoMigrate(&user.User{})

	r := mux.NewRouter()
	r.HandleFunc("/register", RegisterHandler).Methods("POST")

	err = http.ListenAndServe(port, r)
	if err != nil {
		log.Fatalf("error with register service %v\n", err)
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user user.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("cannot decode the user struct given")
		json.NewEncoder(w).Encode(&Msg{Success: false, Msg: "cannot decode user"})
		return
	}
	fmt.Println(user)
	err = ValidateUser(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("validation error: %v \n", err)
		json.NewEncoder(w).Encode(&Msg{Success: false, Msg: err.Error()})
		return
	}
	user.Password, err = HashPassword(user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		log.Printf("failed to hash user password: %v \n", err)
		json.NewEncoder(w).Encode(&Msg{Success: false, Msg: "failed to add user"})
		return
	}
	err = AddUser(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("failed to add user: %v \n", err)
		json.NewEncoder(w).Encode(&Msg{Success: false, Msg: "failed to add user"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&Msg{Success: true, Msg: "user added"})
}

func ValidateUser(user user.User) error {
	if len(user.Username) < 5 {
		return fmt.Errorf("username length is less than 5")
	}
	if len(user.Username) > 32 {
		return fmt.Errorf("username length is grater than 32")
	}
	if len(user.Password) < 8 {
		return fmt.Errorf("password length is less than 8")
	}
	return nil
}

func AddUser(user user.User) error {
	db, err := Connect("root", "1234", "localhost", "dashboard_user")
	if err != nil {
		log.Fatalf("error with connecting to dashboard_user db %v\n", err)
		return err
	}
	tx := db.Create(&user)
	return tx.Error
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
