# gocloser

```
go get github.com/artarts36/gocloser
```

**gocloser** - go package for closing resources (db, storage, etc)

```go
package main

import (
	"database/sql"
	"time"

	"github.com/artarts36/gocloser"
)

func main() {
	db, _ := sql.Open("postgres", "dsn")

	gocloser.Add("db", db.Close)

	// graceful shutdown

	gocloser.Subscribe(30 * time.Second)
}
```
