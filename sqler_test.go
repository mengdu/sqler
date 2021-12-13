package sqler

import (
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
	o.Add("id", DESC)
	o.Add("field1", ASC)
	o.Add("field2", DESC)

	expected := "order by id desc, field1 asc, field2 desc"
	if o.String() != expected {
		t.Errorf("Expecting: %s\nGot: %s", expected, o.String())
	}

	o2 := &Order{}
	if o2.String() != "" {
		t.Errorf("Expecting: \"\"\nGot: %s", o2.String())
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

	// todo
}
