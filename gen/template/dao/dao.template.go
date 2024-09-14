package dao

import (
    "gorm.io/gorm"
    {{- if .ProjectName}}
    "{{.ProjectName}}/{{.ModelPackage}}"
    {{- else}}
    "{{.ModelPackage}}"
    {{- end}}
)

type {{.ModelName}}Dao struct {
    db *gorm.DB
}

func New{{.ModelName}}Dao(db *gorm.DB) *{{.ModelName}}Dao {
    return &{{.ModelName}}Dao{
        db: db,
    }
}

func (r *{{.ModelName}}Dao) Create(m *{{.ModelPackage}}.{{.ModelName}}) error {
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

// Update record
func (r *{{.ModelName}}Dao) Update(m *{{.ModelPackage}}.{{.ModelName}}) error {
    tx := r.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    if err := tx.Model(m).Update(m).Error; err != nil {
        return err
    }

    return tx.Commit().Error
}

// Delete record
func (r *{{.ModelName}}Dao) Delete(m *{{.ModelPackage}}.{{.ModelName}}) error {
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

func (r *{{.ModelName}}Dao) GetById(id *{{.IdType}}) (*{{.ModelPackage}}.{{.ModelName}}, error) {
    var m {{.ModelPackage}}.{{.ModelName}}
    if err := r.db.First(&m, id).Error; err != nil {
        return nil, err
    }
    return &m, nil
}

// Find records with pagination, search, and sorting
func (r *{{.ModelName}}Dao) Find(offset, limit int, search, sort string) (int64, []{{.ModelPackage}}.{{.ModelName}}, error) {
    var total int64
    var results []{{.ModelPackage}}.{{.ModelName}}
    query := r.db.Model(&{{.ModelPackage}}.{{.ModelName}}{})

    // Apply search condition
    if search != "" {
        query = query.Where("name LIKE ?", "%"+search+"%")
    }

    // Get total count
    if err := query.Count(&total).Error; err != nil {
        return 0, nil, err
    }

    // Apply sorting
    if sort != "" {
        query = query.Order(sort)
    }

    // Apply pagination
    if err := query.Offset(offset).Limit(limit).Find(&results).Error; err != nil {
        return 0, nil, err
    }

    return total, results, nil
}

{{range .Relations}}
// Find {{$.ModelName}} with {{.Name}}
func (r *{{$.ModelName}}Dao) Get{{$.ModelName}}With{{.Name}}(id *{{$.IdType}}) (*{{$.ModelPackage}}.{{$.ModelName}}, error) {
    var m {{$.ModelPackage}}.{{$.ModelName}}
    if err := r.db.Preload("{{.Name}}").First(&m, id).Error; err != nil {
        return nil, err
    }
    return &m, nil
}

// Add {{.Name}} to {{$.ModelName}}
func (r *{{$.ModelName}}Dao) Add{{.Name}}To{{$.ModelName}}(m *{{$.ModelPackage}}.{{$.ModelName}}, {{.Name}} *{{$.ModelPackage}}.{{.ModelName}}) error {
    return r.db.Model(m).Association("{{.Name}}").Append({{.Name}})
}

// Remove {{.Name}} from {{$.ModelName}}
func (r *{{$.ModelName}}Dao) Remove{{.Name}}From{{$.ModelName}}(m *{{$.ModelPackage}}.{{$.ModelName}}, {{.Name}} *{{$.ModelPackage}}.{{.ModelName}}) error {
    return r.db.Model(m).Association("{{.Name}}").Delete({{.Name}})
}
{{end}}
