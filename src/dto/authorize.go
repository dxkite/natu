package dto

import (
	"time"

	"dxkite.cn/meownest/pkg/identity"
	"dxkite.cn/meownest/src/constant"
	"dxkite.cn/meownest/src/entity"
	"dxkite.cn/meownest/src/value"
)

// Authorize
type Authorize struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 鉴权备注
	Name string `json:"name"`
	// 鉴权类型
	Type string `json:"type"`
	// 描述
	Description string `json:"description"`
	// 鉴权属性
	Attribute *value.AuthorizeAttribute `json:"attribute"`
}

func NewAuthorize(ent *entity.Authorize) *Authorize {
	obj := new(Authorize)
	obj.Id = identity.Format(constant.AuthorizePrefix, ent.Id)
	obj.CreatedAt = ent.CreatedAt
	obj.UpdatedAt = ent.UpdatedAt
	obj.Name = ent.Name
	obj.Description = ent.Description
	obj.Type = ent.Type
	obj.Attribute = ent.Attribute
	return obj
}
