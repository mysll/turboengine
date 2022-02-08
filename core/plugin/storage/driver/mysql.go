package driver

import (
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"turboengine/gameplay/dao"
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

func (m *MysqlDao) Create(name string, model dao.Persistent) error {
	return m.db.AutoMigrate(model)
}

func (m *MysqlDao) Find(id uint64, data dao.Persistent) error {
	if err := m.db.Find(data, id).Error; err != nil {
		return err
	}
	return nil
}

func (m *MysqlDao) FindBy(data dao.Persistent, where string, args ...any) error {
	if err := m.db.Where(where, args...).Find(data).Error; err != nil {
		return err
	}
	return nil
}

func (m *MysqlDao) FindAll(data any, where string, args ...any) error {
	if err := m.db.Where(where, args...).Find(data).Error; err != nil {
		return err
	}
	return nil
}

func (m *MysqlDao) Save(data dao.Persistent) (uint64, error) {
	if err := m.db.Create(data).Error; err != nil {
		return 0, err
	}
	return data.DBId(), nil
}

func (m *MysqlDao) Update(data dao.Persistent) error {
	if data.DBId() == 0 {
		return errors.New("data dbid is zero")
	}
	if err := m.db.Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func (m *MysqlDao) Del(data dao.Persistent) error {
	if err := m.db.Delete(data).Error; err != nil {
		return err
	}
	return nil
}

func (m *MysqlDao) DelBy(data dao.Persistent, where string, args ...any) error {
	if err := m.db.Where(where, args...).Delete(data).Error; err != nil {
		return err
	}
	return nil
}

func init() {

}
