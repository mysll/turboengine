package driver

import (
	"testing"
	"time"
	"turboengine/common/datatype"
)

type TestObject struct {
	ID        uint64
	Name      string `gorm:"size:64"`
	Mp        int32
	Hp        float32
	MaxHp     float64
	Sex       int           `gorm:"size:1"`
	Pos       datatype.Vec3 `gorm:"type:varbinary(64)"`
	Pos2      datatype.Vec2 `gorm:"type:varbinary(32)"`
	Data      []byte        `gorm:"type:varbinary(128)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t TestObject) DBId() uint64 {
	return t.ID
}

func TestMysqlDao_Connect(t *testing.T) {
	db := new(MysqlDao)
	err := db.Connect("root:123456@tcp(127.0.0.1:3306)/turbo?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		t.Fatal(err)
	}
	db.Create("TestObject", &TestObject{})
	obj := &TestObject{
		Name:  "sll",
		Mp:    1,
		MaxHp: 3,
		Hp:    2,
		Sex:   1,
		Pos:   datatype.Vec3{1, 2, 3},
		Pos2:  datatype.Vec2{4, 5},
		Data:  []byte("test"),
	}

	t.Log(db.Save(obj))
	var obj2 TestObject
	db.Find(obj.ID, &obj2)
	t.Log(obj2)
	obj2.Name = "sll2"
	db.Update(obj2)
	var obj3 TestObject
	db.FindBy(&obj3, "name=?", "hello")
	t.Log(obj3)

	var objs []TestObject
	db.FindAll(&objs, "")
	t.Log(objs)

	db.DelBy(&TestObject{}, "name=?", "sll2")
}
