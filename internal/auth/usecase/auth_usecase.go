package usecase

import (
	"context"
	"cool-games/internal/domain"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
    userRepo       domain.UserRepository
    customerRepo   domain.CustomerRepository
    jwtSecret      string
    contextTimeout time.Duration
}

func NewAuthUsecase(
    repo domain.UserRepository, 
    cRepo domain.CustomerRepository,
    secret string, 
    timeout time.Duration,
) domain.AuthUsecase {
    return &authUsecase{
        userRepo:      repo,
        customerRepo:  cRepo,
        jwtSecret:     secret,
        contextTimeout: timeout,
    }
}

func (u *authUsecase) Register(ctx context.Context, user *domain.User) (domain.AuthResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return domain.AuthResponse{}, err
    }
    user.HashedPassword = string(hashedPassword)

    if err := u.userRepo.Create(ctx, user); err != nil {
        return domain.AuthResponse{}, err
    }

    switch user.Role {
case "customer":
        customer := &domain.Customer{
            UserID: user.ID,
            CustomerName: user.Email,
        }
        _ = u.customerRepo.Create(ctx, customer)
    case "publisher":
        err := u.userRepo.CreatePublisher(ctx, user.ID, user.Email)
        if err != nil {
            return domain.AuthResponse{}, err
        }
    }

    token, _ := u.generateJWT(*user)
    return domain.AuthResponse{Token: token, User: *user}, nil
}

func (u *authUsecase) Login(ctx context.Context, req domain.LoginRequest) (domain.AuthResponse, error) {
	c, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	user, err := u.userRepo.GetByEmail(c, req.Email)
	if err != nil {
		return domain.AuthResponse{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password)); err != nil {
		return domain.AuthResponse{}, err
	}

	token, _ := u.generateJWT(user)
	return domain.AuthResponse{Token: token, User: user}, nil
}

func (u *authUsecase) generateJWT(user domain.User) (string, error) {
    claims := jwt.MapClaims{
        "user_id": user.ID,
        "role":    user.Role,
        "exp":     time.Now().Add(time.Hour * 72).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(u.jwtSecret))
}