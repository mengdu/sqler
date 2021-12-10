# Sqler

A SQL builder for Golang.

```go
s := sqler.New()

s.Select("select id, username, nickname, age, status, sex, createdAt from")
s.Table("users")
s.Where(func(where *sqler.Condition) {
    where.And("age >= ?", 18)
    where.And("sex = ?", 1)
    where.Or(func(or *sqler.Block) {
        or.Add("status = ?", 1)
        or.Add("status = ?", 2)
    })
})
s.Order(func(order *sqler.Order) {
    order.Add("id", sqler.DESC)
    order.Add("age", sqler.ASC)
})
s.Limit(0, 5)
// select id, username, nickname, age, status, sex, createdAt from users where age >= ? and sex = ? and (status = ? or status = ?) order by id desc, age asc limit ?, ? [18 1 1 2 0 5]
fmt.Println(s.Do())
// select count(1) as count from users where age >= ? and sex = ? and (status = ? or status = ?) [18 1 1 2]
fmt.Println(s.DoCount())
```
