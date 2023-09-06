package base

import (
	"errors"
	"fmt"

	"gorm.io/hints"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// gorm 查询基类
// 最佳实践见：models/nativeapp/deviceinfo.go 或者 models/nativeapp/devicecheckin.go

const (
	// proxy 支持的hints
	HintsReadWrite = "/*#mode=READWRITE*/"
	HintsReadOnly  = "/*#mode=READONLY*/"
)

// 删除状态字段
const (
	DeletedNo  = iota //未删除
	DeletedYes        //已删除
)

type (
	// BaseModel db操作基础接口
	BaseModel interface {
		GetDB(ctx *gin.Context) *gorm.DB
		GetOne(ctx *gin.Context, dest interface{}, options ...OptionFunc) error
		GetByCond(ctx *gin.Context, dest interface{}, options ...OptionFunc) error
		GetById(ctx *gin.Context, id uint64) (dest interface{}, err error)
		Create(ctx *gin.Context, data interface{}, options ...OptionFunc) (rowAffects int64, err error)
		CreateInBatches(ctx *gin.Context, data interface{}, batchSize int, options ...OptionFunc) (rowAffects int64, err error)
		Update(ctx *gin.Context, update interface{}, options ...OptionFunc) (rowAffects int64, err error)
		Delete(ctx *gin.Context, data interface{}, options ...OptionFunc) (rowAffects int64, err error) // 谨慎操作，过于危险，不建议使用
		DeleteById(ctx *gin.Context, id uint64) (rowAffects int64, err error)
		Count(ctx *gin.Context, options ...OptionFunc) (int64, error)
		Clauses(ctx *gin.Context, cond ...clause.Expression) (tx *gorm.DB)
		Upsert(ctx *gin.Context, data interface{}, columns []clause.Column, doUpdate clause.Set, options ...OptionFunc) (rowAffects int64, err error)
	}

	defaultBaseModel struct {
		tableName string
		db        *gorm.DB
	}
)

func NewBaseModel(db *gorm.DB, tableName string) BaseModel {
	return &defaultBaseModel{
		tableName: tableName,
		db:        db,
	}
}

func (m *defaultBaseModel) GetDB(ctx *gin.Context) *gorm.DB {
	return m.db.WithContext(ctx).Table(m.tableName)
}

func (m *defaultBaseModel) buildOption(ctx *gin.Context, opts ...OptionFunc) *gorm.DB {
	db := m.GetDB(ctx)
	for _, op := range opts {
		db = op(db)
	}
	return db
}

func (m *defaultBaseModel) GetOne(ctx *gin.Context, dest interface{}, options ...OptionFunc) (err error) {
	return m.buildOption(ctx, options...).Take(dest).Error
}

func (m *defaultBaseModel) GetByCond(ctx *gin.Context, dest interface{}, options ...OptionFunc) (err error) {
	return m.buildOption(ctx, options...).Find(dest).Error
}

func (m *defaultBaseModel) GetById(ctx *gin.Context, id uint64) (dest interface{}, err error) {
	if err = m.buildOption(ctx, WithId(id)).Take(&dest).Error; err != nil {
		return nil, err
	}
	return
}

func (m *defaultBaseModel) Create(ctx *gin.Context, data interface{}, options ...OptionFunc) (rowAffects int64, err error) {
	db := m.buildOption(ctx, options...).Create(data)
	return db.RowsAffected, db.Error
}

func (m *defaultBaseModel) CreateInBatches(ctx *gin.Context, data interface{}, batchSize int, options ...OptionFunc) (rowAffects int64, err error) {
	db := m.buildOption(ctx, options...).CreateInBatches(data, batchSize)
	return db.RowsAffected, db.Error
}

func (m *defaultBaseModel) Update(ctx *gin.Context, data interface{}, options ...OptionFunc) (rowAffects int64, err error) {
	db := m.buildOption(ctx, options...).Updates(data)
	return db.RowsAffected, db.Error
}

func (m *defaultBaseModel) Delete(ctx *gin.Context, data interface{}, options ...OptionFunc) (rowAffects int64, err error) {
	db := m.buildOption(ctx, options...).Delete(data, options)
	return db.RowsAffected, db.Error
}

func (m *defaultBaseModel) DeleteById(ctx *gin.Context, id uint64) (rowAffects int64, err error) {
	db := m.buildOption(ctx, WithId(id)).Delete(nil)
	return db.RowsAffected, db.Error
}

func (m *defaultBaseModel) Count(ctx *gin.Context, options ...OptionFunc) (count int64, err error) {
	err = m.buildOption(ctx, options...).Count(&count).Error
	return
}

func (m *defaultBaseModel) Clauses(ctx *gin.Context, cond ...clause.Expression) (tx *gorm.DB) {
	return m.GetDB(ctx).Clauses(cond...)
}

func (m *defaultBaseModel) Upsert(ctx *gin.Context, data interface{}, columns []clause.Column, doUpdate clause.Set, options ...OptionFunc) (rowAffects int64, err error) {
	if len(columns) == 0 {
		return 0, errors.New("clause.Column must > 1")
	}
	if len(doUpdate) == 0 {
		return 0, errors.New("clause.Set must > 1")
	}

	db := m.buildOption(ctx, options...).Clauses(clause.OnConflict{
		Columns:   columns,
		DoUpdates: doUpdate,
	}).Create(data)

	return db.RowsAffected, db.Error
}

//---------------------------------------------Options---------------------------------------------

type OptionFunc func(*gorm.DB) *gorm.DB

func TableName(tableName string) OptionFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table(tableName)
	}
}

func Select(fields string, args ...interface{}) OptionFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(fields, args...)
	}
}

func SelectDistinct(fields string, args ...interface{}) OptionFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(fmt.Sprintf("distinct %s", fields), args...)
	}
}

func ReadMaster() OptionFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Clauses(hints.New(HintsReadWrite))
	}
}

func Where(query interface{}, args ...interface{}) OptionFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(query, args...)
	}
}

func Limit(limit int) OptionFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	}
}

func Offset(offset int) OptionFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset)
	}
}

func Preload(query string, args ...interface{}) OptionFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(query, args...)
	}
}

func Order(order string) OptionFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(order)
	}
}

func Group(group string) OptionFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Group(group)
	}
}

func WithId(id uint64) OptionFunc {
	return Where("id = ?", id)
}

func WithIds(ids []uint64) OptionFunc {
	return Where("id in (?)", ids)
}

func WithCuid(cuid string) OptionFunc {
	return Where("cuid = ?", cuid)
}

func WithDelete() OptionFunc {
	return Where("deleted = ?", DeletedYes)
}

func WithUnDelete() OptionFunc {
	return Where("deleted = ?", DeletedNo)
}
