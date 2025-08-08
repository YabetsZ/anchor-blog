package usersvc

import (
	AppError "anchor-blog/internal/errors"
	"context"
	"log"
)

func (us *UserServices) PromoteUserToAdmin(ctx context.Context, promoterID, targetUserID string) error {
	// Fetch the user to be promoted
	targetUser, err := us.userRepo.GetUserByID(ctx, targetUserID)
	if err != nil {
		log.Println("user promotion failed because target user isn't found: ", err)
		return err
	}

	// Check if user is already an admin
	if targetUser.Role == "admin" {
		return AppError.ErrUserAlreadyAdmin
	}

	return us.userRepo.UpdateUserRole(ctx, promoterID, targetUserID, "admin")
}

func (us *UserServices) DemoteAdminToUser(ctx context.Context, demoterID, targetAdminID string) error {
	// Safety check: an admin cannot demote themselves
	if demoterID == targetAdminID {
		return AppError.ErrCannotDemoteThemselves
	}

	// Fetch the user to be demoted
	targetUser, err := us.userRepo.GetUserByID(ctx, targetAdminID)
	if err != nil {
		return err
	}

	// Check if user is already a regular user
	if targetUser.Role == "user" {
		return AppError.ErrUserNotAdmin
	}

	return us.userRepo.UpdateUserRole(ctx, demoterID, targetAdminID, "user")
}
