package postgres

import (
	"context"

	"github.com/TezzBhandari/frs"
	"github.com/jackc/pgx/v5"
)

var _ frs.FundRaiserService = (*FundRaiserService)(nil)

type FundRaiserService struct {
	db *DB
}

func NewFundRaiserService(db *DB) *FundRaiserService {
	return &FundRaiserService{db: db}
}

func (fr *FundRaiserService) CreateFundRaiser(ctx context.Context, fundRaiser *frs.FundRaiser) error {
	tx, err := fr.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	fundRaiser.ID = fr.db.snowflake.Generate().Int64()
	fundRaiser.CreatedAt = tx.Now
	fundRaiser.UpdatedAt = fundRaiser.CreatedAt

	err = createFundRaiser(ctx, tx, fundRaiser)
	if err != nil {
		return nil
	}
	return tx.Commit(ctx)
}

func (fr *FundRaiserService) FindFundRaiser(ctx context.Context, filterFundRaiser *frs.FilterFundRaiser) ([]*frs.FundRaiser, int, error) {
	return nil, 0, nil
}

func (fr *FundRaiserService) FindFundRaiserById(ctx context.Context, id int) (*frs.FundRaiser, error) {
	return nil, nil
}

func (fr *FundRaiserService) UpdateFundRaiser(ctx context.Context, id int, updFundRaiser *frs.UpdateFundRaiser) (*frs.FundRaiser, error) {
	return nil, nil
}

func (fr *FundRaiserService) DeleteFundRaiser(ctx context.Context, id int) error {
	return nil
}

func createFundRaiser(ctx context.Context, tx *Tx, fundRaiser *frs.FundRaiser) error {
	if err := fundRaiser.Validate(); err != nil {
		return err
	}

	insertFundRaiserQuery := `
		INSERT INTO fundraisers (id, title, story, cover_img, target_amount, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $6, $7);
	`

	_, err := tx.Exec(ctx, insertFundRaiserQuery, fundRaiser.ID, fundRaiser.Title, fundRaiser.Story, fundRaiser.CoverImg, fundRaiser.TargetAmount, fundRaiser.CreatedAt, fundRaiser.UpdatedAt)
	if err != nil {
		return err
	}

	return nil

}
