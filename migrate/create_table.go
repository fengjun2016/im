package migrate

import (
	"imWebSocket/app"
	"imWebSocket/model"
)

func CreateTable() {
	app.DB.AutoMigrate(
		&model.User{},
		&model.Contact{},
		&model.ContactGroups{},
	)
}
