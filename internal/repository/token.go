package repository

import (
	"context"

	"github.com/hunderaweke/metsasft/internal/domain"
	"github.com/sv-tools/mongoifc"
)

type tokenRepository struct {
	collec mongoifc.Collection
	ctx    context.Context
}

func NewTokenRepository(db mongoifc.Database, ctx context.Context) domain.TokenRepository {
	return &tokenRepository{collec: db.Collection("tokens"), ctx: ctx}
}

func (r *tokenRepository) CreateToken(token domain.Token) error {
	t, err := r.GetTokenByEmail(token.Email)
	if t != (domain.Token{}) {
		r.DeleteToken(token.Email)
	}
	_, err = r.collec.InsertOne(r.ctx, token, nil)
	if err != nil {
		return err
	}
	return nil
}

func (r *tokenRepository) GetTokenByEmail(email string) (domain.Token, error) {
	var token domain.Token
	res := r.collec.FindOne(r.ctx, map[string]interface{}{"email": email}, nil)
	if res.Err() != nil {
		return domain.Token{}, res.Err()
	}
	err := res.Decode(&token)
	if err != nil {
		return domain.Token{}, err
	}
	return token, nil
}

func (r *tokenRepository) DeleteToken(email string) error {
	_, err := r.collec.DeleteOne(r.ctx, map[string]interface{}{"email": email}, nil)
	if err != nil {
		return err
	}
	return nil
}
