package model

import (
	"errors"
	"strings"
	"time"

	"github.com/rs/xid"

	"github.com/jinzhu/gorm"
)

type Model struct {
	ID        string     `gorm:"primary_key;type:varchar(32)" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}

func (m *Model) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("id", xid.New().String())
}

func buildWhere(rawQuery string, db *gorm.DB) (*gorm.DB, error) {
	if rawQuery != "" {
		querys := strings.Split(rawQuery, ",")
		for _, query := range querys {
			oneQuery := strings.Split(query, ":")
			if len(oneQuery) != 2 && len(oneQuery) != 3 {
				return db, errors.New("parseRawQuery error, rawQuery should like: 'title:=:golang,name:like:%jason%,id:<:100' , if the whereType is '=', you can omitted it: title:golang, notice: '%' after encode is %25")
			}
			if len(oneQuery) == 2 {
				field := oneQuery[0]
				whereType := "="
				value := oneQuery[1]
				db = db.Where("`"+field+"`"+" "+whereType+" "+"?", value)
			}
			if len(oneQuery) == 3 {
				field := oneQuery[0]
				whereType := oneQuery[1]
				value := oneQuery[2]
				db = db.Where(field+" "+whereType+" "+"?", value)
			}
		}
	}
	return db, nil
}

func buildOrder(rawOrder string, db *gorm.DB) (*gorm.DB, error) {

	if rawOrder != "" {
		orders := strings.Split(rawOrder, ",")
		for _, order := range orders {
			oneOrder := strings.Split(order, ":")
			if len(oneOrder) != 1 && len(oneOrder) != 2 {
				return db, errors.New("parse rawOrder error, rawOrder should like:'created:desc,id:asc,name', orderType default is asc")
			}

			if len(oneOrder) == 1 {
				field := oneOrder[0]
				db = db.Order(field)
			}
			if len(oneOrder) == 2 {
				field := oneOrder[0]
				orderType := oneOrder[1]
				db = db.Order(field + " " + orderType)
			}
		}
	}
	return db, nil
}
