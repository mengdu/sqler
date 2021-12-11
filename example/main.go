package main

import (
	"fmt"
	"time"

	"github.com/mengdu/sqler"
)

func main() {
	s := sqler.New()

	s.Select("select id, username, age, status from")
	s.Table("users as a")
	s.Join("left join user_token as b on b.userId = a.id")
	s.JoinWithOn("inner join posts as b", func(on *sqler.Condition) {
		on.And("b.userId = a.id")
	})

	s.Where(func(where *sqler.Condition) {
		where.And("username = ?", "admin")
		where.And("(status = 1 or status = 2)")
		where.And("age >= ?", 18)
		where.And("createdAt >= ?", time.Now())
		where.And("find_in_set(?, tags)", "demo")
		where.Or(func(or *sqler.Or) {
			or.Add("id = ?", 1)
			or.Add("name like ?", "%test")
		})

		where.Or(func(or *sqler.Or) {
			or.Add("test = 1")
			or.Add("test = 2")
		})
	})

	s.Having(func(having *sqler.Condition) {
		having.And("a = 1")
	})

	s.Group("id, age")
	s.Order(func(order *sqler.Order) {
		order.Add("id", sqler.DESC)
		order.Add("age", sqler.ASC)
	})
	s.Limit(0, 10)

	fmt.Println(s.Do())
	fmt.Println(s.DoCount())

	s1 := sqler.New()

	s1.Select("select id, username, nickname, age, status, sex, createdAt from")
	s1.Table("users")
	s1.Where(func(where *sqler.Condition) {
		where.And("age >= ?", 18)
		where.And("sex = ?", 1)
		where.Or(func(or *sqler.Or) {
			or.Add("status = ?", 1)
			or.Add("status = ?", 2)
		})
	})
	s1.Order(func(order *sqler.Order) {
		order.Add("id", sqler.DESC)
		order.Add("age", sqler.ASC)
	})
	s1.Limit(0, 5)
	// select id, username, nickname, age, status, sex, createdAt from users where age >= ? and sex = ? and (status = ? or status = ?) order by id desc, age asc limit ?, ? [18 1 1 2 0 5]
	fmt.Println(s1.Do())
	// select count(1) as count from users where age >= ? and sex = ? and (status = ? or status = ?) [18 1 1 2]
	fmt.Println(s1.DoCount())

	// Condition builder
	w := sqler.NewCondition("where")

	w.And("field1 = ?", 1)
	w.And("field2 = ?", 2)
	w.Or(func(or *sqler.Or) {
		or.Add("field3 = ?", 3)
		or.And(func(and *sqler.Condition) {
			and.And("field4 = ?", 4)
			and.And("field5 = ?", 5)
		})
	})

	fmt.Println(w.Do()) // where field1 = ? and field2 = ? and (field3 = ? or (field4 = ? and field5 = ?)) [1 2 3 4 5]
}
