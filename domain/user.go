package domain

import "golang.org/x/crypto/bcrypt"

type User struct {
	Name           string
	hashedPassword string
}

func NewUser(name, password string) (*User, error) {
	user := User{name, ""}
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
	if err != nil {
		return false
	}
	return true
}

type UserRepository interface {
	Add(User) bool
	Delete(name string) (*User, bool)
	FindByName(name string) (*User, bool)
}

type InMemoryUserRepository struct {
	users map[string]*User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{users: make(map[string]*User)}
}

func (i *InMemoryUserRepository) Add(user *User) bool {
	if _, ok := i.users[user.Name]; ok {
		return false
	}
	i.users[user.Name] = user
	return true
}

func (i *InMemoryUserRepository) Delete(name string) (user *User, ok bool) {
	if user, ok = i.users[name]; !ok {
		return
	}
	delete(i.users, name)
	return
}

func (i *InMemoryUserRepository) FindByName(name string) (user *User, ok bool) {
	user, ok = i.users[name]
	return
}
