package main

import (
	"fmt"
	"time"

	"github.com/mengdu/sqler"
	"github.com/mengdu/sqler/utils"
)

func main() {
	s := sqler.New()

	s.Select("select id, username, age, status from users as a")
	s.Join("left join user_token as b on b.userId = a.id")
	s.JoinWithOn("inner join posts as b", func(on *utils.Condition) {
		on.And("b.userId = a.id")
	})

	s.Where(func(where *utils.Condition) {
		where.And("username = ?", "admin")
		where.And("(status = 1 or status = 2)")
		where.And("age >= ?", 18)
		where.And("createdAt >= ?", time.Now())
		where.And("find_in_set(?, tags)", "demo")
		where.Or(func(or *utils.Block) {
			or.Add("id = ?", 1)
			or.Add("name like ?", "%test")
		})
	})

	s.Having(func(having *utils.Condition) {
		having.And("a = 1")
	})

	fmt.Println(s.Do())

	where := utils.NewCondition("where")
	where.And("username = ?", "admin")
	where.And("(status = 1 or status = 2)")
	where.And("age >= ?", 18)
	where.And("createdAt >= ?", time.Now())
	where.And("find_in_set(?, tags)", "demo")
	where.Or(func(or *utils.Block) {
		or.Add("id = ?", 1)
		or.Add("username like ?", "test%")
		or.Add("nickname like ?", "%test%")
	})

	fmt.Println(where.Do())
}
