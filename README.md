# JsonDB

A simple json db

# Usage

## Insert record

```go
package main

import (
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
		os.Exit(1)
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
}


```

## Update record

```go

package main

import (
	"github.com/uole/jsondb"
	"os"
)

func main() {
	db := jsondb.New()
	if err = db.Open(context.Background(), os.TempDir()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err = db.Update(context.Background(), &User{
		ID:   1,
		Name: "World",
	}); err != nil {
		fmt.Println(err)
	}
}

```

## Delete record

```go

package main

import (
	"github.com/uole/jsondb"
	"os"
)

func main() {
	db := jsondb.New()
	if err = db.Open(context.Background(), os.TempDir()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err = db.Delete(context.Background(), &User{
		ID: 1,
	}); err != nil {
		fmt.Println(err)
	}
}
```

## Search

```go

package main

import (
	"github.com/uole/jsondb"
	"os"
)

func main() {
	db := jsondb.New()
	if err = db.Open(context.Background(), os.TempDir()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	models := make([]*User, 0)
	if err = db.Query(context.Background(), &User{}).Where("name", jsondb.ExprEqual, "hello").Where("id", jsondb.ExprEqual, 1).Find(&models); err != nil {
		fmt.Println(err)
	}
}
```