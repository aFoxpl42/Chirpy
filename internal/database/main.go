package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type DB struct {
	path string
	mux *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users map[int]User `json:"user"`
}

type User struct {
	ID int `json:"id"`
	Email string `json:"email"`
}

type Chirp struct {
	ID int `json:"id"`
	Body string `json:""body`
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux: &sync.RWMutex{}, 
	}
	err := db.EnsureDB()
	return db, err
}

func (db *DB) CreateUser(email string) (User, error) {
	dbStructure, err := db.LoadDB()
	if err != nil {
		return User{}, err
	}

	id := len(dbStructure.Users) + 1
	user := User{
		ID: id,
		Email: email,
	}
	dbStructure.Users[id] = user

	err = db.WriteDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetUsers() ([]User, error) {
	dbStructure, err := db.LoadDB()
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(dbStructure.Users))
	for _, user := range dbStructure.Users {
		users = append(users, user)
	}

	return users, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.LoadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID: id,
		Body: body,
	}
	dbStructure.Chirps[id] = chirp

	err = db.WriteDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.LoadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStructure, err := db.LoadDB()
	if err != nil {
		return Chirp{}, err
	}
	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, os.ErrNotExist
	}

	return chirp, nil
}

func (db *DB) CreateDB() error {
	dbStructure := DBStructure{
		Chirps: map[int]Chirp{},
	}
	return db.WriteDB(dbStructure)
}

func (db *DB) EnsureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.CreateDB()
	}
	return err
}

func (db *DB) LoadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStructure := DBStructure{}
	dat, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}
	err = json.Unmarshal(dat, &dbStructure)
	if err != nil {
		return dbStructure, err
	}
	return dbStructure, nil
}

func (db *DB) WriteDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, dat, 0600)
	if err != nil {
		return err
	}
	return nil
}