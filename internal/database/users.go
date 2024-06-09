package database

import (
	"errors"
)

var ErrAlreadyExists = errors.New("already exists")

type User struct {
	ID int `json:"id"`
	Email string `json:"email"`
	HashedPassword string `json:"password"`
}

func (db *DB) CreateUser(email, hashedPassword string) (User, error) {
	dbStructure, err := db.LoadDB()
	if err != nil {
		return User{}, err
	}

	if _, err := db.GetUserByEmail(email); !errors.Is(err, ErrNotExist) {
		return User{}, ErrAlreadyExists
	}

	id := len(dbStructure.Users) + 1
	user := User{
		ID: id,
		Email: email,
		HashedPassword: hashedPassword,
	}
	dbStructure.Users[id] = user

	err = db.WriteDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) UpdateUser(ID int, email, hashedPassword string) (User, error) {
	dbStructure, err := db.LoadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[ID]
	if !ok {
		return User{}, errors.New("user with this ID doesn't exist")
	}
	user.Email = email
	user.HashedPassword = hashedPassword
	dbStructure.Users[ID] = user

	err = db.WriteDB(dbStructure)
	if err != nil {
		return User{}, err
	}
	return user, nil
}


func (db *DB) GetUserByID(ID int) (User, error) {
	dbStructure, err := db.LoadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[ID]
	if !ok {
		return User{}, ErrNotExist
	}

	return user, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.LoadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return User{}, ErrNotExist
}