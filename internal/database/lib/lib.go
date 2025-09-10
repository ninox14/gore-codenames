package lib

import (
	"context"

	"github.com/google/uuid"
	"github.com/ninox14/gore-codenames/internal/database/sqlc"
)

func QuietFindAndDeleteUserById(ctx context.Context, queries *sqlc.Queries, userId uuid.UUID) {
	user, err := queries.GetUserByID(ctx, userId)

	if err != nil {
		// Do nothing
		return
	}

	if err := queries.DeleteUser(ctx, user.ID); err != nil {
		// Do nothing
		return
	}
}
