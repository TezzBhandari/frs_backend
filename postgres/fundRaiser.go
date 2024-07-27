package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/TezzBhandari/frs"
	"github.com/TezzBhandari/frs/utils"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
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

	log.Debug().Msg("hi there fund raising works")

	err = createFundRaiser(ctx, tx, fundRaiser)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (fr *FundRaiserService) FindFundRaiser(ctx context.Context, filterFundRaiser *frs.FilterFundRaiser) ([]*frs.FundRaiser, int, error) {
	tx, err := fr.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, 0, err
	}

	defer tx.Rollback(ctx)
	fundRaisers, n, err := findFundRaiser(ctx, tx, filterFundRaiser)
	if err != nil {
		return nil, 0, err
	}
	return fundRaisers, n, nil
}

func (fr *FundRaiserService) FindFundRaiserById(ctx context.Context, id int64) (*frs.FundRaiser, error) {
	tx, err := fr.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)
	log.Debug().Bool("hi ", true).Msg("")
	fundRaiser, err := findFundRaiserById(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	return fundRaiser, nil
}

func (fr *FundRaiserService) UpdateFundRaiser(ctx context.Context, id int64, updFundRaiser *frs.UpdateFundRaiser) (*frs.FundRaiser, error) {
	return nil, nil
}

func (fr *FundRaiserService) DeleteFundRaiser(ctx context.Context, id int64) error {
	return nil
}

func createFundRaiser(ctx context.Context, tx *Tx, fundRaiser *frs.FundRaiser) error {
	if err := fundRaiser.Validate(); err != nil {
		fmt.Println(err)
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

func findFundRaiser(ctx context.Context, tx *Tx, filterFundRaiser *frs.FilterFundRaiser) ([]*frs.FundRaiser, int, error) {
	where := []string{"1 = 1"}
	args := []any{}
	i := 1
	if filterFundRaiser.ID != nil {
		where = append(where, fmt.Sprintf("id = $%d", i))
		args = append(args, *filterFundRaiser.ID)
		i++
	}

	if filterFundRaiser.Title != nil {
		where = append(where, fmt.Sprintf("title = $%d", i))
		args = append(args, *filterFundRaiser.Title)
		i++
	}

	if filterFundRaiser.CreatedAt != nil {
		where = append(where, fmt.Sprintf("created_at = $%d", i))
		args = append(args, *filterFundRaiser.CreatedAt)
		i++
	}

	whereClause := strings.Join(where, " AND ")

	findFundRaiserQuery := `
		SELECT id, title, story, target_amount, cover_img, created_at, updated_at FROM fundraisers WHERE
	` + whereClause + `
		ORDER BY created_at DESC
	` + formatLimitAndOffset(filterFundRaiser.Limit, filterFundRaiser.Offset)

	log.Debug().Str("sql query", findFundRaiserQuery).Msg("")

	rows, err := tx.Query(ctx, findFundRaiserQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	fundRaisers := make([]*frs.FundRaiser, 0)

	for rows.Next() {
		fundRaiser := frs.FundRaiser{}
		if err := rows.Scan(&fundRaiser.ID, &fundRaiser.Title, &fundRaiser.Story, &fundRaiser.TargetAmount, &fundRaiser.CoverImg, &fundRaiser.CreatedAt, &fundRaiser.UpdatedAt); err != nil {
			return nil, 0, err
		}
		fundRaisers = append(fundRaisers, &fundRaiser)
		if err := rows.Err(); err != nil {
			return nil, 0, err
		}
	}

	return fundRaisers, len(fundRaisers), nil
}

func findFundRaiserById(ctx context.Context, tx *Tx, id int64) (*frs.FundRaiser, error) {
	log.Debug().Str("hey", "hey hey").Msg("")

	fundRaiser, n, err := findFundRaiser(ctx, tx, &frs.FilterFundRaiser{ID: &id})
	log.Debug().Bool("no err", true)
	if err != nil {
		log.Error().Err(err).Msg("erro while executing query")
		return nil, err
	}

	if n == 0 {
		return nil, frs.Errorf(frs.ENOTFOUND, utils.DoesNotExistMsg("fundraiser"))
	}

	return fundRaiser[0], nil
}
