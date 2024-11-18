package domain

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             string
	Name           string
	hashedPassword string
}

func NewUser(name, password string) (*User, error) {
	user := User{uuid.New().String(), name, ""}
	err := user.SetPassword(password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	u.hashedPassword = string(hashedPassword)
	return nil
}

func (u *User) PasswordIsValid(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.hashedPassword), []byte(password))
	return err == nil
}

type UserRepository interface {
	Add(*User) bool
	GetAll() []*User
	FindByID(string) (user *User, userExists bool)
	FindByName(string) (user *User, userExists bool)
	Delete(string) (user *User, userExists bool)
}

type InMemoryUserRepository struct {
	users map[string]*User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{users: make(map[string]*User)}
}

func (i *InMemoryUserRepository) Add(user *User) bool {
	if _, userExists := i.users[user.Name]; userExists {
		return false
	}
	if _, userExists := i.FindByID(user.ID); userExists {
		return false
	}
	i.users[user.Name] = user
	return true
}

func (i *InMemoryUserRepository) GetAll() []*User {
	users := make([]*User, 0)
	for _, user := range i.users {
		users = append(users, user)
	}
	return users
}

func (i *InMemoryUserRepository) FindByID(userID string) (*User, bool) {
	for _, user := range i.users {
		if user.ID == userID {
			return user, true
		}
	}
	return nil, false
}

func (i *InMemoryUserRepository) FindByName(name string) (user *User, userExists bool) {
	user, userExists = i.users[name]
	return
}

func (i *InMemoryUserRepository) Delete(name string) (user *User, userExists bool) {
	if user, userExists = i.users[name]; !userExists {
		return
	}
	delete(i.users, name)
	return
}
