package formatters

import "automotive/models"

type RoleFormatter struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type RoleWithUsersFormatter struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	Description *string         `json:"description"`
	Users       []UserFormatter `json:"users"`
}

func FormatRole(role models.Role) RoleFormatter {
	roleFormatter := RoleFormatter{}

	roleFormatter.ID = role.ID
	roleFormatter.Name = role.Name
	roleFormatter.Description = role.Description

	return roleFormatter
}

func FormatRoles(roles []models.Role) []RoleFormatter {
	rolesFormatter := []RoleFormatter{}

	for _, role := range roles {
		roleFormatter := FormatRole(role)
		rolesFormatter = append(rolesFormatter, roleFormatter)
	}

	return rolesFormatter
}
