# Sqler

A SQL builder for Golang.

[Demo](example/main.go) | [Test](sqler_test.go)

```go
func main() {
	q := sqler.New()
	q.SelectString("id, username, nickname, type, age, status, createdAt")
	q.From("users")
	q.Where(func(where *sqler.Condition) {
		where.And("status = ?", 1)
		where.Or(func(or *sqler.Or) {
			or.Add("type = ?", 1)
			or.Add("type = ?", 2)
		})
	})
	q.Order(func(order *sqler.Order) {
		order.Add("age", sqler.DESC)
		order.Add("id", sqler.ASC)
	})
	// q.OrderString("age desc, id asc")
	q.Limit(0, 10)

	fmt.Println(q.Do()) // select id, username, nickname, type, age, status, createdAt from users where status = ? and (type = ? or type = ?) order by age desc, id asc limit ?, ? [1 1 2 0 10] <nil>
	fmt.Println(q.DoCount()) // select count(1) as count from users where status = ? and (type = ? or type = ?) [1 1 2] <nil>
}
```

### Condition

**sqler.NewCondition(name string)** Use for `where`, `on` and `having` condition builder.

```go
func main() {
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
	sql, args, err := w.Do()
	fmt.Println("result:", sql, args, err) // "where field1 = ? and field2 in(?, ?, ?, ?) and (field3 = ? or (field4 = ? and field5 = ?))" [1 21 22 23 24 3 4 5] <nil>
}
```

### Order

**sqler.Order{}** Use for `order by` builder.

```go
func main() {
	o := sqler.Order{}
	fmt.Println("sql:", o.String()) // ""
	o.Add("field1", sqler.DESC)
	o.Add("field2", sqler.ASC)
	o.Add("field3", sqler.DESC)
	fmt.Println("sql:", o.String()) // "order by field1 desc, field2 asc, field3 desc"
}
```

### Group

**sqler.Group{}** Use for `group by` builder.

```go
func main() {
	g := sqler.Group{}
	fmt.Println("sql:", g.String()) // ""
	g.Add("field1")
	g.Add("field2")
	g.Add("field3")
	fmt.Println("sql:", g.String()) // "group by field1, field2, field3"
}
```
