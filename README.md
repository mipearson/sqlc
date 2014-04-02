# SQL Composer

[![wercker status](https://app.wercker.com/status/afe1fcd1c7eb84818df88cd4fe732bad/s/ "wercker status")](https://app.wercker.com/project/bykey/afe1fcd1c7eb84818df88cd4fe732bad) [![GoDoc](https://godoc.org/github.com/mipearson/sqlc?status.png)](https://godoc.org/github.com/mipearson/sqlc)

SQL Composer (`sqlc`) makes it easier to join together bits of SQL programatically.

It is very, very similar to [squirrel](https://github.com/lann/squirrel) but has less features.

Most SQL is:

``` sql
SELECT * FROM Users WHERE Name=?
```

SQL Composer isn't for that. You should keep using string literals for that.

SQL Composer is for when you are putting together many different bits of a query together programatically, and where you'd usually use some kind of intelligent string replacement. In our use case, we often do this for searches. SQL Composer lets you do this slightly differently:

``` go
s := sqlc.Statement{}
s.Select("u.*").From("Users u").Where("u.Name = ?", name)

if search.Surname != "" {
  s = s.Where("u.Surname = ?", search.Surname)
}
if search.Role != "" {
  s = s.Join("JOIN Roles r ON u.role_id = r.id")
  s = s.Where("r.Name = ?", search.Role)
}

db.Exec(s.SQL(), s.Args()...)
```

Assuming that `Surname` and `Role` are supplied, calling `s.SQL()` gives you the sql:

``` sql
SELECT u.*
FROM Users u
JOIN Roles r ON u.role_id = r.id
WHERE (u.Name = ?) AND (u.Surname = ?) AND (r.Name = ?)
```

And calling `s.Args()` gives you `name`, `search.Surname` and `search.Role`.

## Features

### Postgres Positional Arguments

PostgreSQL, unlike MySQL and sqlite, uses `$1, $2, $3` in favour of `?, ?, ?` for its positional arguments. SQL Composer will manage this for you:

``` go
s := sqlc.Statement{PostgreSQL: true}
s.Where("Foo = ?", foo).Where("Bar = ?", bar)

```

gives

``` sql
WHERE (Foo = $1) AND (Bar = $2)
```

### Statement Re-Use

Pass-by-value and method chaining allows you to use one base statement in many roles without modifying the original.

For example:

``` go
s = sqlc.Statement{}
s.Select("*").From("Users")

topFiveRows := s.Limit(5).SQL() // SELECT * FROM Users LIMIT 5
allRows := s.SQL()              // SELECT * FROM Users
```
