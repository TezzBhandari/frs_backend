package frs

import (
	"context"
	"time"
)

type FundRaiser struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Story        string    `json:"story"`
	CoverImg     string    `json:"cover_img"`
	TargetAmount float64   `json:"target_amount"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type FilterFundRaiser struct {
	ID        *int64     `json:"id"`
	Title     *string    `json:"title"`
	CreatedAt *time.Time `json:"created_at"`

	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type UpdateFundRaiser struct {
	Title        *string  `json:"title"`
	Story        *string  `json:"story"`
	CoverImg     *string  `json:"cover_img"`
	TargetAmount *float64 `json:"target_amount"`
}

type FundRaiserService interface {
	CreateFundRaiser(ctx context.Context, FundRaiser *FundRaiser) error
	FindFundRaiser(ctx context.Context, filterFundRaiser *FilterFundRaiser) ([]*FundRaiser, int, error)
	FindFundRaiserById(ctx context.Context, id int) (*FundRaiser, error)
	UpdateFundRaiser(ctx context.Context, id int, updFundRaiser *UpdateFundRaiser) (*FundRaiser, error)
	DeleteFundRaiser(ctx context.Context, id int) error
}

func (fr *FundRaiser) Validate() error {
	if fr.Title == "" {
		return Errorf(EBADREQUEST, "fund raiser title is required")
	}

	if fr.Story == "" {
		return Errorf(EBADREQUEST, "fund raiser story is required")
	}

	if fr.CoverImg == "" {
		return Errorf(EBADREQUEST, "fund raiser cover image is required")
	}

	if fr.TargetAmount == 0.0 {
		return Errorf(EBADREQUEST, "fund raiser target amount is required")
	}

	return nil

}
