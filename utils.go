package sqler

import (
	"fmt"
	"strings"
)

type Block struct {
	sqls []string
	args []interface{}
}

func (b *Block) Add(sql string, args ...interface{}) {
	str := strings.Trim(sql, " ")

	if str == "" {
		return
	}

	b.sqls = append(b.sqls, sql)
	b.args = append(b.args, args...)
}

func (b *Block) Set(sql string, args ...interface{}) {
	b.sqls = []string{sql}
	b.args = args
}

func (b *Block) Join(sep string) (sql string, args []interface{}) {
	if len(b.sqls) == 0 {
		return "", []interface{}{}
	}

	return strings.Join(b.sqls, sep), b.args
}

type Condition struct {
	b    Block
	name string
}

func (c *Condition) And(sql string, args ...interface{}) {
	c.b.Add(sql, args...)
}

func (c *Condition) Or(orFn func(or *Block)) {
	o := &Block{}
	orFn(o)
	sql, args := o.Join(" or ")
	c.b.Add(fmt.Sprintf("(%s)", sql), args...)
}

func (c *Condition) Do() (sql string, args []interface{}) {
	sql, args = c.b.Join(" and ")

	if sql == "" {
		return "", []interface{}{}
	}

	if c.name == "" {
		return sql, args
	}

	return c.name + " " + sql, args
}

func NewCondition(name string) *Condition {
	return &Condition{name: name}
}

const (
	DESC = "desc"
	ASC  = "asc"
)

type Order struct {
	b Block
}

func (o *Order) Add(field string, sort string) {
	o.b.Add(fmt.Sprintf("%s %s", field, sort))
}

func (o *Order) String() string {
	if len(o.b.sqls) == 0 {
		return ""
	}

	sql, _ := o.b.Join(", ")
	return "order by " + sql
}
