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

    // Role
    admin := model.Role{Name: "admin"}
    config.DB.FirstOrCreate(&admin, admin)
    config.DB.Model(&admin).Association("Permissions").Append(&perm)
}
