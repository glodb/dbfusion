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

type UserCreateTable struct {
	Id        int    `dbfusion:"id,INT,AUTO_INCREMENT,PRIMARY KEY"`
	Email     string `dbfusion:"email,VARCHAR(255),NOT NULL,UNIQUE" json:"email"`
	Phone     string `dbfusion:"phone,VARCHAR(255),NOT NULL" json:"phone"`
	Password  string `dbfusion:"password,VARCHAR(50),NOT NULL" json:"password,omitempty"`
	FirstName string `dbfusion:"firstName,VARCHAR(50)" json:"firstName"`
	LastName  string `dbfusion:"lastName,VARCHAR(50)" json:"lastName"`
	CreatedAt int    `dbfusion:"createdAt,INT" json:"createdAt"`
	UpdatedAt int    `dbfusion:"updatedAt,INT" json:"updatedAt"`
}

func (ne UserCreateTable) GetNormalIndexes() []string {
	return []string{"id:1,email:1", "id:-1,phone:1"}
}

func (ne UserCreateTable) GetUniqueIndexes() []string {
	return []string{"id:1,phone:1"}
}

func (ne UserCreateTable) GetTextIndex() string {
	return "email"
}

func (ne UserCreateTable) Get2DIndexes() []string {
	return []string{}
}

func (ne UserCreateTable) Get2DSpatialIndexes() []string {
	return []string{}
}

func (ne UserCreateTable) GetHashedIndexes() []string {
	return []string{"id"}
}

// func (ne UserCreateTable) GetSparseIndexes() []string {
// 	return []string{"email"}
// }
