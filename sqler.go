package sqler

import "strings"

func mergeBlock(blocks ...Block) (string, []interface{}, error) {
	sqls := []string{}
	args := []interface{}{}
	for _, v := range blocks {
		chunk, chunkArgs, err := v.Join(" ")

		if err != nil {
			return "", []interface{}{}, err
		}

		if chunk != "" {
			sqls = append(sqls, chunk)
			args = append(args, chunkArgs...)
		}
	}

	return strings.Join(sqls, " "), args, nil
}

type Sqler struct {
	query  Block
	table  Block
	join   Block
	where  Block
	group  Block
	having Block
	order  Block
	limit  Block
}

func (s *Sqler) SelectString(sql string, args ...interface{}) {
	s.query.Set("select")
	s.query.Add(sql, args...)
	s.query.Add("from")
}

func (s *Sqler) Select(selectFn func(field *Block)) {
	b := &Block{}
	selectFn(b)
	sql, args, err := b.Join(", ")

	if err != nil {
		panic(err)
	}

	s.query.Set("select")
	s.query.Add(sql, args...)
	s.query.Add("from")
}

func (s *Sqler) From(sql string, args ...interface{}) {
	s.table.Add(sql, args...)
}

// Join table, can be called multiple times
func (s *Sqler) JoinString(sql string, args ...interface{}) {
	s.join.Add(sql, args...)
}

// Join table, can be called multiple times
func (s *Sqler) Join(sql string, condition func(on *Condition)) {
	s.join.Add(sql)
	on := NewCondition("on")
	condition(on)
	str, args, err := on.Do()

	if err != nil {
		panic(err)
	}

	s.join.Add(str, args...)
}

// Where builder, cannot be called more than once
func (s *Sqler) WhereString(sql string, args ...interface{}) {
	s.where.Set("where")
	s.where.Add(sql, args...)
}

// Where builder, cannot be called more than once
func (s *Sqler) Where(condition func(where *Condition)) {
	where := NewCondition("where")
	condition(where)
	sql, args, err := where.Do()

	if err != nil {
		panic(err)
	}

	s.where.Set(sql, args...)
}

func (s *Sqler) GroupString(fields ...string) {
	s.group.Set("")
	if len(fields) > 0 {
		s.group.Add("group by")
	}

	for _, v := range fields {
		s.group.Add(v)
	}
}

func (s *Sqler) Group(groupFn func(group *Group)) {
	g := &Group{}
	groupFn(g)
	s.group.Set(g.String())
}

func (s *Sqler) HavingString(sql string, args ...interface{}) {
	s.having.Set("having")
	s.having.Add(sql, args...)
}

func (s *Sqler) Having(condition func(having *Condition)) {
	having := NewCondition("having")
	condition(having)
	sql, args, err := having.Do()

	if err != nil {
		panic(err)
	}

	s.having.Set(sql, args...)
}

func (s *Sqler) OrderString(sql string) {
	s.order.Set("order by")
	s.order.Add(sql)
}

func (s *Sqler) Order(order func(order *Order)) {
	o := &Order{}
	order(o)
	s.order.Set(o.String())
}

func (s *Sqler) Limit(offset uint, limit uint) {
	s.limit.Set("limit ?, ?", offset, limit)
}

func (s *Sqler) DoCount(countSql ...string) (string, []interface{}, error) {
	b := &Block{}
	if len(countSql) > 0 {
		b.Add("select")
		b.Add(countSql[0])
		b.Add("from")
	} else {
		b.Add("select count(1) as count from")
	}
	sql, args, err := mergeBlock(*b, s.table, s.join, s.where, s.group, s.having)
	return sql, args, err
}

func (s *Sqler) Do() (string, []interface{}, error) {
	sql, args, err := mergeBlock(s.query, s.table, s.join, s.where, s.group, s.having, s.order, s.limit)
	return sql, args, err
}

func New() *Sqler {
	return &Sqler{}
}
