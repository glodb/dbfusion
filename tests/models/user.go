package models

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"

	"github.com/glodb/dbfusion/hooks"
)

type UserTest struct {
	FirstName string `dbfusion:"firstname"`
	Email     string `dbfusion:"email"`
	Username  string `dbfusion:"username"`
	Password  string `dbfusion:"password"`
	CreatedAt int64  `dbfusion:"createdAt"`
	UpdatedAt int64  `dbfusion:"updatedAt"`
}

func (u UserTest) GetEntityName() string {
	return "users"
}

func (ne UserTest) GetCacheIndexes() []string {
	return []string{"email", "email,password", "email,username"}
}

func (u UserTest) PreInsert() hooks.PreInsert {

	//Sample password hashing to show the effect of pre insert hook
	hasher := md5.New()
	io.WriteString(hasher, u.Password)
	u.Password = fmt.Sprintf("%x", hasher.Sum(nil))
	u.CreatedAt = 0
	return u
}

func (u UserTest) PostInsert() hooks.PostInsert {

	//Sample password hashing to show the effect of pre insert hook
	log.Println(u.FirstName, " inserted")
	return u
}

type NonEntityUserTest struct {
	FirstName string `dbfusion:"firstname"`
	Email     string `dbfusion:"email,index:hash"`
	Username  string `dbfusion:"username"`
	Password  string `dbfusion:"password"`
	CreatedAt int64  `dbfusion:"createdAt"`
	UpdatedAt int64  `dbfusion:"updatedAt"`
}

type Address struct {
	City       string `dbfusion:"city"`
	PostalCode string `dbfusion:"postCode"`
	Line1      string `dbfusion:"line1"`
}

type Vehicles struct {
	Vehicles []string `dbfusion:"vehicles"`
}
type UseWithAddress struct {
	FirstName string   `dbfusion:"firstname"`
	Email     string   `dbfusion:"email"`
	Username  string   `dbfusion:"username"`
	Password  string   `dbfusion:"password"`
	Address   Address  `dbfusion:"address"`
	Vehicles  Vehicles `dbfusion:"cars,omitempty"`
	CreatedAt int64    `dbfusion:"createdAt,omitempty"`
	UpdatedAt int64    `dbfusion:"updatedAt"`
}

func (ne NonEntityUserTest) GetCacheIndexes() []string {
	return []string{"email", "email,password", "email,username"}
}
