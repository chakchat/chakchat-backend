package storage

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/shared/go/postgres"
	"github.com/chakchat/chakchat-backend/user-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type FieldRestrictions struct {
	Field          string
	OpenTo         models.Restriction
	SpecifiedUsers []uuid.UUID
}

type RestrictionStorage struct {
	db postgres.SQLer
}

func NewRestrictionStorage(db postgres.SQLer) *RestrictionStorage {
	return &RestrictionStorage{
		db: db,
	}
}

func (s *RestrictionStorage) GetAllowedUserIDs(ctx context.Context, id uuid.UUID, field string) ([]uuid.UUID, error) {
	var fieldRestriction FieldRestrictions
	q := `SELECT field_name, permitted_user_id
	 FROM users.field_restrictions 
	 WHERE owner_user_id = $1 
		AND field_name = $2`

	rows, err := s.db.Query(ctx, q, id, field)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var specifiedUsers []uuid.UUID
	for rows.Next() {
		var userID uuid.UUID
		if err := rows.Scan(&fieldRestriction.Field, &userID); err != nil {
			return nil, err
		}
		specifiedUsers = append(specifiedUsers, userID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return specifiedUsers, nil
}

func (s *RestrictionStorage) UpdateRestrictions(ctx context.Context, id uuid.UUID, restrictions FieldRestrictions) (*FieldRestrictions, error) {

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var updateQuery string
	if restrictions.Field == "phone" {
		updateQuery = `UPDATE users.user SET phone_visibility = $1 WHERE id = $2`
	} else {
		updateQuery = `UPDATE users.user SET date_of_birth_visibility = $1 WHERE id = $2`
	}

	_, err = tx.Exec(ctx, updateQuery, restrictions.OpenTo, id)
	if err != nil {
		return nil, err
	}

	if restrictions.OpenTo == models.RestrictionSpecified {
		var currentSpecifiedUsers []uuid.UUID
		q := `SELECT permitted_user_id FROM users.field_restrictions WHERE owner_user_id = $1 AND field_name = $2::users.user_field`
		rows, err := tx.Query(ctx, q, id, restrictions.Field)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var userID uuid.UUID
			if err := rows.Scan(&userID); err != nil {
				return nil, err
			}
			currentSpecifiedUsers = append(currentSpecifiedUsers, userID)
		}

		if rows.Err() != nil {
			return nil, errors.New(rows.Err().Error())
		}

		add := recordMisses(currentSpecifiedUsers, restrictions.SpecifiedUsers)
		del := recordMisses(restrictions.SpecifiedUsers, currentSpecifiedUsers)

		if len(del) > 0 {
			q := `DELETE FROM users.field_restrictions WHERE owner_user_id = $1 AND field_name = $2::users.user_field AND permitted_user_id = ANY($3::uuid[])`
			_, err = tx.Exec(ctx, q, id, restrictions.Field, del)
			if err != nil {
				return nil, err
			}
		}

		if len(add) > 0 {
			q := `INSERT INTO users.field_restrictions (owner_user_id, field_name, permitted_user_id) VALUES ($1, $2::users.user_field, $3)`

			for _, userID := range add {
				_, err = tx.Exec(ctx, q, id, restrictions.Field, userID)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &restrictions, nil
}

func recordMisses(orig, comp []uuid.UUID) []uuid.UUID {
	compMap := make(map[uuid.UUID]bool, len(comp))
	for _, t := range comp {
		compMap[t] = true
	}

	var misses []uuid.UUID

	for _, t := range orig {
		if !compMap[t] {
			misses = append(misses, t)
		}
	}

	return misses
}
