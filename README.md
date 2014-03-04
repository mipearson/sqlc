# SQL Composer

[![wercker status](https://app.wercker.com/status/afe1fcd1c7eb84818df88cd4fe732bad/s/ "wercker status")](https://app.wercker.com/project/bykey/afe1fcd1c7eb84818df88cd4fe732bad) [![GoDoc](https://godoc.org/github.com/mipearson/sqlc?status.png)](https://godoc.org/github.com/mipearson/sqlc)

SQL Composer (`sqlc`) makes it easier to join together bits of SQL programatically.

Most SQL is:

``` sql
SELECT * FROM Users WHERE Name=?
```

SQL Composer isn't for that. You should keep using string literals for that.

SQL Composer is for when you are putting together many different bits of a query together programatically, and where you'd usually use some kind of intelligent string replacement. In our use case, we often do this for searches. SQL Composer lets you do this slightly differently:

``` go
s := sqlc.Statement{}
s.Select("*").From("Users").Where("Name = ?", name)

if search.Surname != "" {
  s = s.Where("Surname = ?", search.Surname)
}
if search.Role != "" {
  s = s.From("JOIN Roles ON Users.role_id = Role.id")
  s = s.Where("Roles.Name = ?", search.Role)
}

db.Exec(s.ToSQL())
```

It also supports PostgreSQL-style positional arguments!

``` go
s := sqlc.Statement{PostgreSQL: true}
s.Where("Foo = ?", foo).Where("Bar = ?", bar)

sql, _ := s.ToSql()
// sql == "WHERE (Foo = $1) AND (Bar = $2)"
```

