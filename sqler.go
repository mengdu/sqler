package sqler

import (
	"sqler/utils"
)

type Sqler struct{}

func (s *Sqler) Select(sql string, args ...interface{}) {}

func (s *Sqler) Join(sql string) {}

func (s *Sqler) JoinWithOn(sql string, condition func(on *utils.Condition)) {}

func (s *Sqler) WhereString(sql string, args ...interface{}) {}

func (s *Sqler) Where(condition func(where *utils.Condition)) {
	condition(utils.NewCondition("where"))
}

func (s *Sqler) Group(fields ...string) {}

func (s *Sqler) HavingString(sql string, args ...interface{}) {}

func (s *Sqler) Having(condition func(having *utils.Condition)) {
	condition(utils.NewCondition("having"))
}

func (s *Sqler) Order(field string, sort string) {}

func (s *Sqler) Limit(limit uint) {}

func (s *Sqler) Page(offset uint, Limit uint) {}

func (s *Sqler) Count(sql string, args ...interface{}) {}

func (s *Sqler) Do() (sql string, args []interface{}) {
	return "", []interface{}{}
}

func New() *Sqler {
	return &Sqler{}
}
