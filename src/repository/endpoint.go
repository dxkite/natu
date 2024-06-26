package repository

import (
	"context"

	"dxkite.cn/meownest/pkg/database"
	"dxkite.cn/meownest/src/entity"
	"gorm.io/gorm"
)

type Endpoint interface {
	Create(ctx context.Context, endpoint *entity.Endpoint) (*entity.Endpoint, error)
	Get(ctx context.Context, id uint64) (*entity.Endpoint, error)
	BatchGet(ctx context.Context, ids []uint64) ([]*entity.Endpoint, error)
	List(ctx context.Context, param *ListEndpointParam) (*ListEndpointResult, error)
	Update(ctx context.Context, id uint64, ent *entity.Endpoint) error
	Delete(ctx context.Context, id uint64) error
}

func NewEndpoint() Endpoint {
	return &endpoint{}
}

type endpoint struct {
}

func (r *endpoint) Create(ctx context.Context, endpoint *entity.Endpoint) (*entity.Endpoint, error) {
	if err := r.dataSource(ctx).Create(&endpoint).Error; err != nil {
		return nil, err
	}
	return endpoint, nil
}

func (r *endpoint) Get(ctx context.Context, id uint64) (*entity.Endpoint, error) {
	var cert entity.Endpoint
	if err := r.dataSource(ctx).Where("id = ?", id).First(&cert).Error; err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r *endpoint) BatchGet(ctx context.Context, ids []uint64) ([]*entity.Endpoint, error) {
	var items []*entity.Endpoint
	if err := r.dataSource(ctx).Where("id in ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

type ListEndpointParam struct {
	Name string
	// pagination
	Page         int
	PerPage      int
	IncludeTotal bool
}

type ListEndpointResult struct {
	Data  []*entity.Endpoint
	Total int64
}

func (r *endpoint) List(ctx context.Context, param *ListEndpointParam) (*ListEndpointResult, error) {
	var items []*entity.Endpoint
	db := r.dataSource(ctx)

	// condition
	condition := func(db *gorm.DB) *gorm.DB {
		if param.Name != "" {
			db = db.Where("name like ?", "%"+param.Name+"%")
		}
		return db
	}

	// pagination
	query := db.Scopes(condition)
	if param.Page > 0 && param.PerPage > 0 {
		query.Offset((param.Page - 1) * param.PerPage).Limit(param.PerPage)
	}

	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}

	rst := &ListEndpointResult{}
	rst.Data = items

	if param.IncludeTotal {
		if err := db.Model(entity.Endpoint{}).Scopes(condition).Count(&rst.Total).Error; err != nil {
			return nil, err
		}
	}

	return rst, nil
}

func (r *endpoint) Update(ctx context.Context, id uint64, ent *entity.Endpoint) error {
	if err := r.dataSource(ctx).Where("id = ?", id).Updates(&ent).Error; err != nil {
		return err
	}
	return nil
}

func (r *endpoint) Delete(ctx context.Context, id uint64) error {
	if err := r.dataSource(ctx).Where("id = ?", id).Delete(entity.Endpoint{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *endpoint) dataSource(ctx context.Context) *gorm.DB {
	return database.Get(ctx).Engine().(*gorm.DB)
}
