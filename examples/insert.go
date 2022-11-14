package main

import (
	"context"
	"fmt"
	"github.com/uole/jsondb"
	"os"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) PrimaryKey() string {
	return "id"
}

func main() {
	var (
		err error
	)
	db := jsondb.New()
	if err = db.Open(context.Background(), os.TempDir()); err != nil {
		fmt.Println(err)
	}
	if err = db.Insert(context.Background(), &User{
		ID:   1,
		Name: "Hello",
	}); err != nil {
		fmt.Println(err)
	}
	if err = db.ReplaceInto(context.Background(), &User{
		ID:   2,
		Name: "Test",
	}); err != nil {
		fmt.Println(err)
	}
	user := &User{ID: 2}
	if err = db.Find(context.Background(), user); err != nil {
		fmt.Println(err)
	}
}
