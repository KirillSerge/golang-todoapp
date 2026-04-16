package users_service

import (
	"context"
	"fmt"

	"github.com/KirillSerge/golang-todoapp/internal/core/domain"
)

func (s *UsersService) GetUser(ctx context.Context, id int) (domain.User, error) {
	user, err := s.userRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("gwt user from repository: %w", err)
	}

	return user, nil
}
