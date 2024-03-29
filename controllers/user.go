package controllers

import (
	"books-app/models"
	"books-app/repository/user"
	"books-app/utils"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"fmt"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var users []models.User

type CurrentUser struct {
	JWT  models.JWT
	User models.User
}

func (c Controller) Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var jwt models.JWT
		var error models.Error

		json.NewDecoder(r.Body).Decode(&user)

		if user.Email == "" {
			error.Message = "Email is missing."
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		if user.Password == "" {
			error.Message = "Password is missing."
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		password := user.Password

		userRepo := userRepository.UserRepository{}

		user, err := userRepo.Login(db, user)

		log.Println(err)

		if err != nil {
			if err == sql.ErrNoRows {
				error.Message = "The user does not exist"
				utils.RespondWithError(w, http.StatusBadRequest, error)
				return
			} else {
				log.Fatal(err)
			}
		}

		hashedPassword := user.Password

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

		if err != nil {
			error.Message = "Invalid Password"
			utils.RespondWithError(w, http.StatusUnauthorized, error)
			return
		}

		token, err := utils.GenerateToken(user)

		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
		jwt.Token = token

		user.Password = ""

		currentUser, err := json.Marshal(user)

		if err != nil {
			panic(err)
		}

		os.Setenv("currentUser", string(currentUser))

		utils.ResponseJSON(w, CurrentUser{
			JWT:  jwt,
			User: user,
		})
	}
}

func (c Controller) Signup(db *sql.DB) http.HandlerFunc {
	fmt.Println("test")
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var error models.Error

		json.NewDecoder(r.Body).Decode(&user)

		if user.Email == "" {
			error.Message = "Email is missing."
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		if user.Password == "" {
			error.Message = "Password is missing."
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

		if err != nil {
			log.Fatal(err)
		}

		user.Password = string(hash)

		userRepo := userRepository.UserRepository{}
		user = userRepo.Signup(db, user)

		if err != nil {
			error.Message = "Server error."
			utils.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}

		user.Password = ""

		w.Header().Set("Content-Type", "application/json")
		utils.ResponseJSON(w, user)
	}
}

func (c Controller) GetUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		users = []models.User{}
		userRepo := userRepository.UserRepository{}

		users = userRepo.GetUsers(db, user, users)

		json.NewEncoder(w).Encode(users)
	}
}

func (c Controller) GetUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		params := mux.Vars(r)

		users = []models.User{}
		userRepo := userRepository.UserRepository{}

		id, err := strconv.Atoi(params["id"])
		logFatal(err)

		user = userRepo.GetUser(db, user, id)

		json.NewEncoder(w).Encode(user)
	}
}

func (c Controller) RemoveUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		userRepo := userRepository.UserRepository{}

		id, err := strconv.Atoi(params["id"])
		logFatal(err)

		rowsDeleted := userRepo.RemoveUser(db, id)

		json.NewEncoder(w).Encode(rowsDeleted)
	}
}
