package sqler

import (
	"reflect"
	"testing"
)

func TestBlock(t *testing.T) {
	b := Block{}
	str, args, err := b.Join("")

	if err != nil {
		t.Fatal(err.Error())
	}

	if str != "" || len(args) != 0 {
		t.Errorf("Expecting: \"\"\nGot: %s", str)
	}

	b.Set("") // reset
	if str != "" || len(args) != 0 {
		t.Errorf("Expecting: \"\"\nGot: %s", str)
	}

	b1 := Block{}
	b1.Add("field1 = ?", 1)
	b1.Add("field2 = ?", 2)
	b1.Add("field3 = ?", 3)
	str1, args1, err := b1.Join(", ")

	if err != nil {
		t.Fatal(err.Error())
	}

	if str1 != "field1 = ?, field2 = ?, field3 = ?" {
		t.Errorf("Expecting: %s\nGot: %s", "field1 = ?, field2 = ?, field3 = ?", str1)
	}

	if len(args1) != 3 {
		t.Errorf("Expecting: %v\nGot: %v", []int{1, 2, 3}, args1)
	}
}

func TestCondition(t *testing.T) {
	c := NewCondition("where")
	str, args, err := c.Do()

	if err != nil {
		t.Errorf(err.Error())
	}

	if str != "" || len(args) != 0 {
		t.Errorf("Expecting: \"\"\nGot: %s", str)
	}

	c1 := NewCondition("where")
	c1.And("field1 = ?", 1)
	c1.And("field2 = ?", 2)
	c1.And("field3 = ? and field4 = ?", 3, 4)
	c1.Or(func(or *Or) {
		or.Add("field5 = ?", 5)
		or.Add("field6 = ?", 6)
	})

	str1, args1, err := c1.Do()

	if err != nil {
		t.Errorf(err.Error())
	}

	expected := "where field1 = ? and field2 = ? and field3 = ? and field4 = ? and (field5 = ? or field6 = ?)"
	if str1 != expected {
		t.Errorf("Expecting: %s\nGot: %s", expected, str1)
	}

	if len(args1) != 6 {
		t.Errorf("Expecting: args len = %d", 6)
	}

	c2 := NewCondition("on")
	c2.And("field in(?)", []int{1, 2, 3, 4, 5})
	str2, _, err := c2.Do()

	if err != nil {
		t.Errorf(err.Error())
	}
	expected2 := "on field in(?, ?, ?, ?, ?)"
	if str2 != expected2 {
		t.Errorf("Expecting: %s\nGot: %s", expected2, str2)
	}
}

func TestOr(t *testing.T) {
	o := Or{}
	sql, args, err := o.Do()

	if err != nil {
		t.Errorf(err.Error())
	}

	if sql != "" || len(args) != 0 {
		t.Errorf("Expecting: \"\"\nGot: %s, %#v", sql, args)
	}

	o1 := Or{}
	o1.Add("field1 = ?", 1)
	o1.Add("field2 = ?", 2)
	o1.And(func(and *Condition) {
		and.And("feild2 = ?", 3)
		and.And("feild3 = ?", 4)
		and.Or(func(or *Or) {
			or.Add("feild4 = ?", 5)
			or.Add("feild4 = ?", 6)
		})
	})

	sql, args, err = o1.Do()

	if err != nil {
		t.Errorf(err.Error())
	}

	expected := "field1 = ? or field2 = ? or (feild2 = ? and feild3 = ? and (feild4 = ? or feild4 = ?))"
	if sql != expected {
		t.Errorf("Expecting: %s\nGot: %s", expected, sql)
	}
	if len(args) != 6 {
		t.Errorf("Expecting: args len = %d, but: %d", 6, len(args))
	}

	o2 := Or{}
	o2.Add("field in(?)", []int{1, 2, 3, 4, 5})
	str2, _, err := o2.Do()

	if err != nil {
		t.Errorf(err.Error())
	}
	expected2 := "field in(?, ?, ?, ?, ?)"
	if str2 != expected2 {
		t.Errorf("Expecting: %s\nGot: %s", expected2, str2)
	}
}

func TestOrder(t *testing.T) {
	o := &Order{}

	if o.String() != "" {
		t.Errorf("Expecting: \"\"\nGot: %s", o.String())
	}

	o.Add("field1", DESC)
	o.Add("field2", ASC)
	o.Add("field3", DESC)

	expected := "order by field1 desc, field2 asc, field3 desc"
	if o.String() != expected {
		t.Errorf("Expecting: %s\nGot: %s", expected, o.String())
	}
}

func TestGroup(t *testing.T) {
	g := &Group{}

	if g.String() != "" {
		t.Errorf("Expecting: %s\nGot: %s", "", g.String())
	}

	g.Add("field1")
	g.Add("field2")
	g.Add("field3")

	expected := "group by field1, field2, field3"
	if g.String() != expected {
		t.Errorf("Expecting: %s\nGot: %s", expected, g.String())
	}
}

func TestIn(t *testing.T) {
	str, _, err := In("?", []int{1, 2, 3, 4, 5})
	if err != nil {
		t.Errorf(err.Error())
	}
	expected := "?, ?, ?, ?, ?"
	if str != expected {
		t.Errorf("Expecting: %s\nGot: %s", expected, str)
	}

	str1, _, err := In("?, ?", []int{1, 2})
	if err == nil {
		t.Errorf("Expecting: %s, Got: %v", err, nil)
	}

	if str1 != "" {
		t.Errorf("Expecting: \"\", Got: %s", str)
	}
}

func TestSqler(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Error(err)
		}
	}()

	// select
	func() {
		q := New()
		q.SelectString("field1, field2, ? as field3, ? as field4", 1, 2)
		expected := "select field1, field2, ? as field3, ? as field4 from"

		sql, args, err := q.Do()
		if err != nil {
			t.Fatal(err)
		}
		if sql != expected {
			t.Errorf("Expecting: %s\nGot: %s", expected, sql)
		}
		if len(args) != 2 {
			t.Errorf("Expecting: %v\nGot: %v", []int{1, 2}, args)
		}

		q.Select(func(field *Block) {
			field.Add("field1")
			field.Add("field2")
			field.Add("? as field3", 1)
			field.Add("? as field4", 2)
		})

		sql, args, err = q.Do()
		if err != nil {
			t.Fatal(err)
		}
		if sql != expected {
			t.Errorf("Expecting: %s\nGot: %s", expected, sql)
		}
		if len(args) != 2 {
			t.Errorf("Expecting: %v\nGot: %v", []int{1, 2}, args)
		}
	}()

	// join
	func() {
		q := New()
		q.JoinString("inner join table1 as b")
		q.JoinString("left join table2 as c")
		q.JoinString("left join table3 as d and d.field = ?", 1)
		q.Join("left jion table4 as e", func(on *Condition) {
			on.And("e.id = d.id")
			on.And("e.status = ?", 2)
		})
		sql, args, err := q.Do()

		if err != nil {
			panic(err)
		}
		expected := "inner join table1 as b left join table2 as c left join table3 as d and d.field = ? left jion table4 as e on e.id = d.id and e.status = ?"
		if sql != expected {
			t.Errorf("Expecting: %s\nGot: %s", expected, sql)
		}
		if len(args) != 2 {
			t.Errorf("Expecting: %v\nGot: %v", []int{1, 2}, args)
		}
	}()

	// where
	func() {
		q := New()
		q.WhereString("field1 = ? and field2 = ? and field3 in(?)", 1, 2, []int{11, 22, 33})
		sql, args, err := q.Do()

		if err != nil {
			panic(err)
		}
		expected := "where field1 = ? and field2 = ? and field3 in(?, ?, ?)"
		expectArgs := []interface{}{1, 2, 11, 22, 33}
		if sql != expected {
			t.Errorf("Expecting: %s\nGot: %s", expected, sql)
		}

		if !reflect.DeepEqual(expectArgs, args) {
			t.Errorf("Expecting: %v\nGot: %v", expectArgs, args)
		}

		q.Where(func(where *Condition) {
			where.And("field1 = ?", 1)
			where.And("field2 = ?", 2)
			where.And("field3 in(?)", []int{11, 22, 33})
		})

		sql, args, err = q.Do()
		if err != nil {
			panic(err)
		}

		if sql != expected {
			t.Errorf("Expecting: %s\nGot: %s", expected, sql)
		}

		if !reflect.DeepEqual(expectArgs, args) {
			t.Errorf("Expecting: %v\nGot: %v", expectArgs, args)
		}

		q.Where(func(where *Condition) {})

		sql, args, err = q.Do()
		if err != nil {
			panic(err)
		}
		expected = ""
		if sql != expected {
			t.Errorf("Expecting: %s\nGot: %s", expected, sql)
		}
		if len(args) != 0 {
			t.Errorf("Expecting: %v\nGot: %v", []int{}, args)
		}
	}()

	// group
	func() {
		q := New()
		q.GroupString("field1, field2, field3")
		sql, _, err := q.Do()

		if err != nil {
			panic(err)
		}
		expected := "group by field1, field2, field3"

		if sql != expected {
			t.Errorf("Expecting: %s\nGot: 1%s", expected, sql)
		}

		q.Group(func(group *Group) {
			group.Add("field1")
			group.Add("field2")
			group.Add("field3")
		})

		sql, _, err = q.Do()

		if err != nil {
			panic(err)
		}

		if sql != expected {
			t.Errorf("Expecting: %s\nGot: %s", expected, sql)
		}

		q.Group(func(group *Group) {})
		sql, _, err = q.Do()
		if sql != "" {
			t.Errorf("Expecting: %s\nGot: %s", "", sql)
		}
	}()

	// having
	func() {
		q := New()
		q.HavingString("field1 = ? and field2 = ? and field3 in(?)", 1, 2, []int{11, 22, 33})
		sql, args, err := q.Do()

		if err != nil {
			panic(err)
		}
		expected := "having field1 = ? and field2 = ? and field3 in(?, ?, ?)"
		expectArgs := []interface{}{1, 2, 11, 22, 33}
		if sql != expected {
			t.Errorf("Expecting: %s\nGot: %s", expected, sql)
		}

		if !reflect.DeepEqual(expectArgs, args) {
			t.Errorf("Expecting: %v\nGot: %v", expectArgs, args)
		}

		q.Having(func(having *Condition) {
			having.And("field1 = ?", 1)
			having.And("field2 = ?", 2)
			having.And("field3 in(?)", []int{11, 22, 33})
		})

		sql, args, err = q.Do()
		if err != nil {
			panic(err)
		}

		if sql != expected {
			t.Errorf("Expecting: %s\nGot: %s", expected, sql)
		}

		if !reflect.DeepEqual(expectArgs, args) {
			t.Errorf("Expecting: %v\nGot: %v", expectArgs, args)
		}

		q.Having(func(having *Condition) {})

		sql, args, err = q.Do()
		if err != nil {
			panic(err)
		}
		expected = ""
		if sql != expected {
			t.Errorf("Expecting: %s\nGot: %s", expected, sql)
		}
		if len(args) != 0 {
			t.Errorf("Expecting: %v\nGot: %v", []int{}, args)
		}
	}()

	// group
	func() {
		q := New()
		q.OrderString("field1 desc, field2 asc, field3 desc")
		sql, _, err := q.Do()

		if err != nil {
			panic(err)
		}
		expected := "order by field1 desc, field2 asc, field3 desc"

		if sql != expected {
			t.Errorf("Expecting: %s\nGot: %s", expected, sql)
		}

		q.Order(func(order *Order) {
			order.Add("field1", DESC)
			order.Add("field2", ASC)
			order.Add("field3", DESC)
		})

		sql, _, err = q.Do()

		if err != nil {
			panic(err)
		}

		if sql != expected {
			t.Errorf("Expecting: %s\nGot: %s", expected, sql)
		}

		q.Order(func(order *Order) {})
		sql, _, err = q.Do()
		if sql != "" {
			t.Errorf("Expecting: %s\nGot: %s", "", sql)
		}
	}()

	// DoCount
	func() {
		q := New()
		q.SelectString("field1, field2, filed3")
		q.From("table as a")
		q.WhereString("field1 = ? and field2 = ? and field3 in(?)", 1, 2, []int{11, 22, 33})
		q.GroupString("field1, field2")
		q.HavingString("field1 = ? and field2 = ? and field3 in(?)", 3, 4, []int{44, 55, 66})
		q.OrderString("field1 desc")
		q.Limit(0, 10)
		sql, args, err := q.Do()

		if err != nil {
			panic(err)
		}

		expected := "select field1, field2, filed3 from table as a where field1 = ? and field2 = ? and field3 in(?, ?, ?) group by field1, field2 having field1 = ? and field2 = ? and field3 in(?, ?, ?) order by field1 desc limit ?, ?"
		expectArgs := []interface{}{1, 2, 11, 22, 33, 3, 4, 44, 55, 66, uint(0), uint(10)}
		if sql != expected {
			t.Errorf("Expecting: %s\nGot: %s", expected, sql)
		}

		if !reflect.DeepEqual(expectArgs, args) {
			t.Errorf("Expecting: %v\nGot: %v", expectArgs, args)
		}

		countSql, countArgs, err := q.DoCount()
		if err != nil {
			panic(err)
		}
		expected = "select count(1) as count from table as a where field1 = ? and field2 = ? and field3 in(?, ?, ?) group by field1, field2 having field1 = ? and field2 = ? and field3 in(?, ?, ?)"
		expectArgs = []interface{}{1, 2, 11, 22, 33, 3, 4, 44, 55, 66}
		if countSql != expected {
			t.Errorf("Expecting: %s\nGot: %s", expected, countSql)
		}

		if !reflect.DeepEqual(expectArgs, countArgs) {
			t.Errorf("Expecting: %v\nGot: %v", expectArgs, countArgs)
		}

		countSql, countArgs, err = q.DoCount("count(id) count")
		if err != nil {
			panic(err)
		}
		expected = "select count(id) count from table as a where field1 = ? and field2 = ? and field3 in(?, ?, ?) group by field1, field2 having field1 = ? and field2 = ? and field3 in(?, ?, ?)"
		expectArgs = []interface{}{1, 2, 11, 22, 33, 3, 4, 44, 55, 66}

		if countSql != expected {
			t.Errorf("Expecting: %s\nGot: %s", expected, countSql)
		}

		if !reflect.DeepEqual(expectArgs, countArgs) {
			t.Errorf("Expecting: %v\nGot: %v", expectArgs, countArgs)
		}
	}()
}
