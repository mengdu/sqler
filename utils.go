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
	str := strings.Trim(sql, " ")

	if str == "" {
		b.sqls = []string{}
		b.args = []interface{}{}
	} else {
		b.sqls = []string{sql}
		b.args = args
	}
}

func (b *Block) Join(sep string) (string, []interface{}, error) {
	if len(b.sqls) == 0 {
		return "", []interface{}{}, nil
	}

	return In(strings.Join(b.sqls, sep), b.args...)
}

// Condition builder, use for `where`, `on`, `having`
type Condition struct {
	b    Block
	name string
}

func (c *Condition) And(sql string, args ...interface{}) {
	c.b.Add(sql, args...)
}

func (c *Condition) Or(orFn func(or *Or)) {
	o := &Or{}
	orFn(o)
	sql, args, err := o.Do()

	if err != nil {
		panic(err)
	}

	if sql != "" {
		c.b.Add(fmt.Sprintf("(%s)", sql), args...)
	}
}

func (c *Condition) Do() (string, []interface{}, error) {
	sql, args, err := c.b.Join(" and ")

	if err != nil {
		return "", []interface{}{}, err
	}

	if sql == "" {
		return "", []interface{}{}, nil
	}

	if c.name == "" {
		return In(sql, args...)
	}

	return In(c.name+" "+sql, args...)
}

func NewCondition(name string) *Condition {
	return &Condition{name: name}
}

type Or struct {
	b Block
}

func (o *Or) Add(sql string, args ...interface{}) {
	o.b.Add(sql, args...)
}

func (o *Or) And(andFn func(and *Condition)) {
	b := NewCondition("")
	andFn(b)
	sql, args, err := b.Do()

	if err != nil {
		panic(err)
	}

	if sql != "" {
		o.b.Add(fmt.Sprintf("(%s)", sql), args...)
	}
}

func (o *Or) Do() (string, []interface{}, error) {
	sql, args, err := o.b.Join(" or ")

	if err != nil {
		return "", []interface{}{}, err
	}

	if sql == "" {
		return "", []interface{}{}, nil
	}

	return In(sql, args...)
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

	sql, _, _ := o.b.Join(", ")
	return "order by " + sql
}

type Group struct {
	b Block
}

func (g *Group) Add(field string) {
	g.b.Add(field)
}

func (g *Group) String() string {
	if len(g.b.sqls) == 0 {
		return ""
	}

	sql, _, _ := g.b.Join(", ")
	return "group by " + sql
}
