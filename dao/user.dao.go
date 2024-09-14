package dao

import (
	"gorm.io/gorm"
	"model"
)

type GroupDao struct {
	db *gorm.DB
}

func NewGroupDao(db *gorm.DB) *GroupDao {
	return &GroupDao{
		db: db,
	}
}

func (r *GroupDao) Create(m *model.Group) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create record
	if err := tx.Create(m).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// Delete record
func (r *GroupDao) Delete(m *model.Group) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Delete(m).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}

func (r *GroupDao) GetById(id *string) (*model.Group, error) {
	var m model.Group
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// Find records with pagination, search, and sorting
func (r *GroupDao) Find(offset, limit int, search, sort string) ([]model.Group, error) {
	var results []model.Group
	query := r.db.Model(&model.Group{})

	// Apply search condition
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	// Apply sorting
	if sort != "" {
		query = query.Order(sort)
	}

	// Apply pagination
	if err := query.Offset(offset).Limit(limit).Find(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

// Find Group with Cars
func (r *GroupDao) GetGroupWithCars(id *string) (*model.Group, error) {
	var m model.Group
	if err := r.db.Preload("Cars").First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// Find Group with Groups
func (r *GroupDao) GetGroupWithGroups(id *string) (*model.Group, error) {
	var m model.Group
	if err := r.db.Preload("Groups").First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// Find Group with Users
func (r *GroupDao) GetGroupWithUsers(id *string) (*model.Group, error) {
	var m model.Group
	if err := r.db.Preload("Users").First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}
