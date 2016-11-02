# sqloader

Work with external SQL files in Go(lang). Use the queries calling it's name.

Example
-------

Have a SQL file called user.sql

```sql
-- /selectUser
SELECT *
FROM tbl_user u
WHERE u.id = $1
-- /

-- /listUser
SELECT id, name FROM tbl_user u LIMIT 1000
-- /
```

Load it in Go

```go
queries, err := sqloader.NewSQLoader("user.sql")
```

Use it in a DB Query

```go
// Select the name of the user with id 3
var name string
err = db.QueryRow(queries.Get("selectUser"), 3).Scan(&name)
```
