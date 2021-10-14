package builders

import "profira-backend/usecase/viewmodel"

type IUserTokenJwtVmBuilder interface {
	SetToken(token string) IUserTokenJwtVmBuilder

	Token() string

	SetExpiredToken(expiredToken string) IUserTokenJwtVmBuilder

	ExpiredToken() string

	SetRefreshToken(refreshToken string) IUserTokenJwtVmBuilder

	RefreshToken() string

	SetExpiredRefreshToken(expiredRefreshToken string) IUserTokenJwtVmBuilder

	ExpiredRefreshToken() string

	SetIsActive(isActive bool) IUserTokenJwtVmBuilder

	IsActive() bool

	GetUserTokenJwtVm() viewmodel.UserJwtTokenVm
}

type UserTokenJwtVmBuilder struct {
	token               string
	expiredToken        string
	refreshToken        string
	expiredRefreshToken string
	isActive            bool
}

func NewUserTokenJwtVmBuilder() *UserTokenJwtVmBuilder {
	return &UserTokenJwtVmBuilder{}
}

func (b *UserTokenJwtVmBuilder) SetToken(token string) IUserTokenJwtVmBuilder {
	b.token = token

	return b
}

func (b *UserTokenJwtVmBuilder) Token() string {
	return b.token
}

func (b *UserTokenJwtVmBuilder) SetExpiredToken(expiredToken string) IUserTokenJwtVmBuilder {
	b.expiredToken = expiredToken

	return b
}

func (b *UserTokenJwtVmBuilder) ExpiredToken() string {
	return b.expiredToken
}

func (b *UserTokenJwtVmBuilder) SetRefreshToken(refreshToken string) IUserTokenJwtVmBuilder {
	b.refreshToken = refreshToken

	return b
}

func (b *UserTokenJwtVmBuilder) RefreshToken() string {
	return b.refreshToken
}

func (b *UserTokenJwtVmBuilder) SetExpiredRefreshToken(expiredRefreshToken string) IUserTokenJwtVmBuilder {
	b.expiredRefreshToken = expiredRefreshToken

	return b
}

func (b *UserTokenJwtVmBuilder) ExpiredRefreshToken() string {
	return b.expiredRefreshToken
}

func (b *UserTokenJwtVmBuilder) SetIsActive(isActive bool) IUserTokenJwtVmBuilder {
	b.isActive = isActive

	return b
}

func (b *UserTokenJwtVmBuilder) IsActive() bool {
	return b.isActive
}

func (b *UserTokenJwtVmBuilder) GetUserTokenJwtVm() viewmodel.UserJwtTokenVm {
	return viewmodel.UserJwtTokenVm{
		Token:           b.token,
		ExpTime:         b.expiredToken,
		RefreshToken:    b.refreshToken,
		ExpRefreshToken: b.expiredRefreshToken,
		IsActive:        b.isActive,
	}
}
