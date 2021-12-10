package sqler

import "strings"

func mergeBlock(blocks ...Block) (sql string, args []interface{}) {
	sqls := []string{}

	for _, v := range blocks {
		chunk, chunkArgs := v.Join(" ")
		if chunk != "" {
			sqls = append(sqls, chunk)
			args = append(args, chunkArgs...)
		}
	}

	return strings.Join(sqls, " "), args
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

func (s *Sqler) Select(sql string, args ...interface{}) {
	s.query.Add(sql, args...)
}

func (s *Sqler) Table(sql string, args ...interface{}) {
	s.table.Set(sql, args...)
}

func (s *Sqler) Join(sql string) {
	s.join.Add(sql)
}

func (s *Sqler) JoinWithOn(sql string, condition func(on *Condition)) {
	s.join.Add(sql)
	on := NewCondition("on")
	condition(on)
	str, args := on.Do()
	s.join.Add(str, args...)
}

func (s *Sqler) WhereString(sql string, args ...interface{}) {
	s.where.Set(sql, args...)
}

func (s *Sqler) Where(condition func(where *Condition)) {
	where := NewCondition("where")
	condition(where)
	sql, args := where.Do()
	s.where.Set(sql, args...)
}

func (s *Sqler) Group(fields ...string) {
	if len(fields) > 0 {
		s.group.Add("group by")
	}

	for _, v := range fields {
		s.group.Add(v)
	}
}

func (s *Sqler) HavingString(sql string, args ...interface{}) {
	s.having.Set(sql, args...)
}

func (s *Sqler) Having(condition func(having *Condition)) {
	having := NewCondition("having")
	condition(having)
	sql, args := having.Do()
	s.having.Set(sql, args...)
}

func (s *Sqler) OrderString(sql string) {
	s.order.Set(sql)
}

func (s *Sqler) Order(order func(order *Order)) {
	o := &Order{}
	order(o)
	s.order.Set(o.String())
}

func (s *Sqler) Limit(offset uint, limit uint) {
	s.limit.Set("limit ?, ?", offset, limit)
}

func (s *Sqler) DoCount(countSql ...string) (sql string, args []interface{}) {
	b := &Block{}
	if len(countSql) > 0 {
		b.Add(countSql[0])
	} else {
		b.Add("select count(1) as count from")
	}
	sql, args = mergeBlock(*b, s.table, s.join, s.where, s.group, s.having)
	return
}

func (s *Sqler) Do() (sql string, args []interface{}) {
	sql, args = mergeBlock(s.query, s.table, s.join, s.where, s.group, s.having, s.order, s.limit)
	return
}

func New() *Sqler {
	return &Sqler{}
}
