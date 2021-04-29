package driver

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlDao struct {
	db *gorm.DB
}

func (m *MysqlDao) Connect(dsn string) error {
	db, err := gorm.Open(mysql.New(
		mysql.Config{
			DSN:               dsn,
			DefaultStringSize: 256,
		}), &gorm.Config{})
	if err != nil {
		return err
	}
	m.db = db
	return nil
}

func (m *MysqlDao) Create(name string, model interface{}) error {
	return m.db.AutoMigrate(model)
}

func (m *MysqlDao) Select(id uint64, data interface{}) error {
	return nil
}
func (m *MysqlDao) Insert(id uint64, data interface{}) (uint64, error) {
	return 0, nil

}
func (m *MysqlDao) Update(id uint64, data interface{}) error {
	return nil

}
func (m *MysqlDao) Del(name string, id uint64) error {
	return nil
}

func init() {

}
