package usersvc

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	AppError "anchor-blog/internal/errors"
	"anchor-blog/pkg/hashutil"
)

func (us *UserServices) Register(ctx context.Context, userDto *UserDTO) (string, error) {

	user := DTOToEntity(*userDto)

	passwordHash, err := hashutil.HashPassword(userDto.Password)
	if err != nil {
		return "", err
	}
	user.PasswordHash = passwordHash

	if strings.Trim(user.Username, " ") == "" {
		return "", errors.New("username cannot be empty")
	}
	if strings.Trim(user.Email, " ") == "" {
		return "", errors.New("email cannot be empty")
	}

	username, email := strings.Trim(user.Username, " "), strings.Trim(user.Email, " ")
	var validUsername = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*$`)
	if (!validUsername.MatchString(username)) && len(username) < 3 {
		return "", AppError.ErrInvalidUsername
	}

	exists, err := us.userRepo.CheckUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if exists {
		return "", AppError.ErrUsernameTaken
	}
	exists, err = us.userRepo.CheckEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if exists {
		return "", AppError.ErrEmailAlreadyExists
	}

	firstName, lastName := strings.Trim(user.FirstName, " "), strings.Trim(user.LastName, " ")
	if firstName == "" || lastName == "" || len(firstName) < 3 || len(lastName) < 3 {
		return "", AppError.ErrNameCannotEmpty
	}

	user.FirstName, user.LastName = firstName, lastName

	user.Username = username
	user.Email = email

	user.Activated = false
	user.Role = "unverified"

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return us.userRepo.CreateUser(ctx, &user)
}
