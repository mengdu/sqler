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

## Condition

```go
func main() {
  // Condition builder
	w := sqler.NewCondition("where")

	w.And("field1 = ?", 1)
	w.And("field2 in(?)", []int{21, 22, 23, 24})
	w.Or(func(or *sqler.Or) {
		or.Add("field3 = ?", 3)
		or.And(func(and *sqler.Condition) {
			and.And("field4 = ?", 4)
			and.And("field5 = ?", 5)
		})
	})

	fmt.Println(w.Do()) // where field1 = ? and field2 in(?, ?, ?, ?) and (field3 = ? or (field4 = ? and field5 = ?)) [1 21 22 23 24 3 4 5] <nil>
}
```
