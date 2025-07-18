package seed

import (
	"catalog-gin/config"
	"catalog-gin/model"
)

func SeedRolesAndPermissions() {
	// Content Types
	ct := model.ContentType{Name: "user"}
	config.DB.FirstOrCreate(&ct, model.ContentType{Name: "user"})

	// Permissions
	perm := model.Permission{Name: "view_admin", ContentTypeID: ct.ID}
	config.DB.FirstOrCreate(&perm, perm)

	// Permissions for user
	userPerm := model.Permission{Name: "view_user", ContentTypeID: ct.ID}
	config.DB.FirstOrCreate(&userPerm, userPerm)

	// Role
	admin := model.Role{Name: "admin"}
	config.DB.FirstOrCreate(&admin, admin)
	config.DB.Model(&admin).Association("Permissions").Append(&perm)

	// Role
	user := model.Role{Name: "user"}
	config.DB.FirstOrCreate(&user, user)
	config.DB.Model(&user).Association("Permissions").Append(&userPerm)
}
