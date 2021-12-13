package main

import (
	"fmt"

	"github.com/mengdu/sqler"
)

func conditionDemo() {
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

	fmt.Println(w.Do()) // where field1 = ? and field2 in(?, ?, ?, ?) and (field3 = ? or (field4 = ? and field5 = ?)) [1 21 22 23 24 3 4 5]
}

func sqlerDemo() {
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
}

func inDemo() {
	sql, args, err := sqler.In("field1 in(?) and field2 = ?", []interface{}{1, 2, 3, 4}, 5)
	fmt.Println(sql, args, err) // field1 in(?, ?, ?, ?) and field2 = ? [1 2 3 4 5] <nil>
}

func main() {
	sqlerDemo()
	conditionDemo()
	inDemo()
}
