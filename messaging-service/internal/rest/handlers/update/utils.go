package update

import (
	"context"

	"github.com/chakchat/chakchat-backend/shared/go/auth"
	"github.com/google/uuid"
)

func getUserID(ctx context.Context) uuid.UUID {
	return uuid.MustParse(auth.GetClaims(ctx)[auth.ClaimId].(string))
}
