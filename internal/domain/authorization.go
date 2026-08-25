package domain

import (
	"fmt"
	"strings"
)

type Role string

const (
	RoleOperator Role = "operator"
	RoleReviewer Role = "reviewer"
	RoleManager  Role = "manager"
)

func ParseRole(value string) (Role, error) {
	role := Role(strings.ToLower(strings.TrimSpace(value)))
	switch role {
	case RoleOperator, RoleReviewer, RoleManager:
		return role, nil
	default:
		return "", fmt.Errorf("unknown role %q", value)
	}
}

func CanCreate(role Role) bool { return role == RoleOperator || role == RoleManager }

func CanReview(role Role) bool { return role == RoleReviewer || role == RoleManager }

func CanArchiveRole(role Role) bool { return role == RoleManager }

func ActionAllowed(role Role, action string) bool {
	switch strings.ToLower(strings.TrimSpace(action)) {
	case "create", "update":
		return CanCreate(role)
	case "review", "confirm":
		return CanReview(role)
	case "archive":
		return CanArchiveRole(role)
	case "publish":
		return CanReview(role)
	default:
		return false
	}
}

func PermissionMessage(role Role, action string) string {
	if ActionAllowed(role, action) {
		return "allowed"
	}
	return fmt.Sprintf("role %s cannot %s", role, action)
}
