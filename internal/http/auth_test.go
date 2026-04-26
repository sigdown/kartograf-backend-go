package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	apiauth "github.com/sigdown/kartograf-backend-go/internal/auth"
	"github.com/sigdown/kartograf-backend-go/internal/domain"
	"github.com/sigdown/kartograf-backend-go/internal/usecase"
)

func TestGetCurrentUserReturnsAuthorizedUser(t *testing.T) {
	ginTokenManager := apiauth.NewTokenManager("secret", time.Hour)
	token, err := ginTokenManager.Generate(domain.User{
		ID:       7,
		Username: "alice",
		Role:     domain.RoleUser,
	})
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	authService := usecase.NewAuthService(&usecaseTestUserRepo{
		getByIDFn: func(userID int64) domain.User {
			return domain.User{
				ID:          userID,
				Username:    "alice",
				DisplayName: "Alice",
				Email:       "alice@example.com",
				Role:        domain.RoleUser,
			}
		},
	}, ginTokenManager)

	router := NewRouter(Services{
		Auth:   authService,
		Tokens: ginTokenManager,
	})

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}

	var user domain.User
	if err := json.Unmarshal(resp.Body.Bytes(), &user); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if user.ID != 7 {
		t.Fatalf("unexpected user id: %d", user.ID)
	}
	if user.Email != "alice@example.com" {
		t.Fatalf("unexpected email: %s", user.Email)
	}
}

type usecaseTestUserRepo struct {
	getByIDFn func(userID int64) domain.User
}

func (r *usecaseTestUserRepo) Create(_ context.Context, user domain.User) (domain.User, error) {
	return user, nil
}

func (r *usecaseTestUserRepo) GetByID(_ context.Context, userID int64) (domain.User, error) {
	return r.getByIDFn(userID), nil
}

func (r *usecaseTestUserRepo) FindByLogin(_ context.Context, _ string) (domain.User, error) {
	return domain.User{}, nil
}

func (r *usecaseTestUserRepo) Update(_ context.Context, _ int64, _ usecase.UpdateAccountInput) (domain.User, error) {
	return domain.User{}, nil
}

func (r *usecaseTestUserRepo) Delete(_ context.Context, _ int64) error {
	return nil
}
