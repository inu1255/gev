package gev

import (
	"log"

	"github.com/go-xorm/xorm"
)

type IAddressModel interface {
	ISearchModel
}

// Address Entity
type AddressModel struct {
	SearchModel `xorm:"extends"`
	Center      string `json:"center,omitempty" xorm:"not null" gev:"中心经纬度"`
	Citycode    string `json:"citycode,omitempty" xorm:"not null" gev:"城市区号"`
	Level       string `json:"level,omitempty" xorm:"not null" gev:"级别"`
	Name        string `json:"name,omitempty" xorm:"not null" gev:"城市名"`
	ParentId    int    `json:"parent_id,omitempty" xorm:"not null" gev:"父地址"`
	Value       string `json:"value,omitempty" xorm:"not null" gev:"三级地址名|分隔"`
}

func (a *AddressModel) TableName() string {
	return "address"
}

func (a *AddressModel) GetCondition() ISearch {
	return &SearchAddress{}
}

func (a *AddressModel) SearchSession(session *xorm.Session, condition ISearch) error {
	search := condition.(*SearchAddress)
	if search.Keyword != "" {
		session.Where("value like ?", search.Keyword+"%")
	}
	if search.ParentId != 0 {
		session.Where("parent_id=?", search.ParentId)
	}
	return nil
}

func (m *AddressModel) Bind(g ISwagRouter, self IModel) {
	if self == nil {
		self = m
	}
	m.SearchModel.Bind(g, self)
	// 导入地址数据
	if n, err := Db.Limit(10).Count(self); err == nil && n < 1 {
		res, err := Db.ImportFile(PkgPath + "/address.sql")
		log.Println(res, err)
	}
}
