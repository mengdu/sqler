package sqler

import (
	"testing"
)

func TestBlock(t *testing.T) {
	b := Block{}
	str, args := b.Join("")

	if str != "" || len(args) != 0 {
		t.Errorf("Expected: void(\"\")")
	}

	b1 := Block{}
	b1.Add("field1 = ?", 1)
	b1.Add("field2 = ?", 2)
	b1.Add("field3 = ?")
	str1, args1 := b1.Join(", ")

	if str1 != "field1 = ?, field2 = ?, field3 = ?" {
		t.Errorf("Expected: %s", "field1 = ?, field2 = ?, field3 = ?")
	}

	if len(args1) != 2 {
		t.Errorf("Expected: args len = %d", 2)
	}
}

func TestCondition(t *testing.T) {
	c := NewCondition("where")
	str, args := c.Do()

	if str != "" || len(args) != 0 {
		t.Errorf("Expected: void(\"\"), but: %s", str)
	}

	c1 := NewCondition("where")
	c1.And("field1 = ?", 1)
	c1.And("field2 = ?", 2)
	c1.And("field3 = ? and field4 = ?", 3, 4)
	c1.Or(func(or *Block) {
		or.Add("field5 = ?", 5)
		or.Add("field6 = ?", 6)
	})

	str1, args1 := c1.Do()
	expected := "where field1 = ? and field2 = ? and field3 = ? and field4 = ? and (field5 = ? or field6 = ?)"
	if str1 != str1 {
		t.Errorf("Expected: %s", expected)
	}

	if len(args1) != 6 {
		t.Errorf("Expected: args len = %d", 6)
	}
}

func TestOrder(t *testing.T) {
	o := &Order{}
	o.Add("id", DESC)
	o.Add("field1", ASC)
	o.Add("field2", DESC)

	expected := "order by id desc, field1 asc, field2 desc"
	if o.String() != expected {
		t.Errorf("Expected: %s", expected)
	}

	o2 := &Order{}
	if o2.String() != "" {
		t.Errorf("Expected: void(\"\")")
	}
}

func TestSqler(t *testing.T) {
	// todo
}
