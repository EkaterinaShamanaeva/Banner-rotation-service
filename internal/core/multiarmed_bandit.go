package core

import (
	"context"
	"fmt"
	"github.com/EkaterinaShamanaeva/Banner-rotation-service/internal/storage/sqlstorage"
	"github.com/gofrs/uuid"
)

// GetBannerId реализует алгоритм "многорукий бандит".
func GetBannerId(ctx context.Context, s sqlstorage.Storage, slotId uuid.UUID, userGroupId uuid.UUID) (uuid.UUID, error) {
	// список баннеров в ротации указанного слота
	bannerIds, err := s.GetBannersFromSlot(ctx, slotId)
	if err != nil {
		return uuid.Nil, err
	}

	if len(bannerIds) == 0 {
		return uuid.Nil, fmt.Errorf("there are not banners in slot")
	}

	//for _, banner := range bannerIds {
	//
	//}

	return uuid.Nil, nil
}
