package utils

import (
	"fmt"

	"github.com/kaleemubarok/recam/backend/pkg/repository"
)

// GetCredentialsByRole func for getting credentials from a role name.
func GetCredentialsByRole(role string) ([]string, error) {
	// Define credentials variable.
	var credentials []string

	// Switch given role.
	switch role {
	case repository.AdminRoleName:
		// Admin credentials (all access).
		credentials = []string{
			repository.RouteCreateCredential,
			repository.RouteUpdateCredential,
			repository.RouteDeleteCredential,
		}
	case repository.ModeratorRoleName:
		// Moderator credentials (only route creation and update).
		credentials = []string{
			repository.RouteCreateCredential,
			repository.RouteUpdateCredential,
		}
	case repository.UserRoleName:
		// Simple user credentials (only route creation).
		credentials = []string{
			repository.RouteCreateCredential,
		}
	default:
		// Return error message.
		return nil, fmt.Errorf("role '%v' does not exist", role)
	}

	return credentials, nil
}
