package core

import (
	"context"
	"fmt"
	"github.com/EkaterinaShamanaeva/Banner-rotation-service/internal/storage/sqlstorage"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBandit(t *testing.T) {
	s := sqlstorage.New()
	dsn := "postgres://postgres:password@localhost:5432/banner_rotator?sslmode=%s"
	ctx := context.Background()
	dbPool, _ := sqlstorage.Connect(ctx, dsn)
	s.Pool = dbPool
	defer s.Close(ctx)

	bannerId, _ := s.CreateBanner(ctx, "shop banner")
	bannerIdSecond, _ := s.CreateBanner(ctx, "cafe banner")

	userGroupId, _ := s.CreateUserGroup(ctx, "men 30+")

	slotId, _ := s.CreateSlot(ctx, "slot A")

	// slotIdSec, _ := s.CreateSlot(ctx, "slot B")

	_ = s.AddBannerToSlot(ctx, slotId, bannerId)
	_ = s.AddBannerToSlot(ctx, slotId, bannerIdSecond)

	//_ = s.AddBannerToSlot(ctx, slotIdSec, bannerIdSecond)

	_ = s.NewBannerStatistic(ctx, bannerId, slotId, userGroupId)
	_ = s.NewBannerStatistic(ctx, bannerIdSecond, slotId, userGroupId)

	_ = s.ShowBanner(ctx, bannerId, slotId, userGroupId)
	_ = s.ShowBanner(ctx, bannerIdSecond, slotId, userGroupId)

	_ = s.ClickBanner(ctx, bannerIdSecond, slotId, userGroupId)

	id, err := GetBannerId(ctx, *s, slotId, userGroupId)
	fmt.Println(id)
	require.NoError(t, err)

	_ = s.ClickBanner(ctx, id, slotId, userGroupId)
	_ = s.ShowBanner(ctx, bannerId, slotId, userGroupId)
	_ = s.ShowBanner(ctx, bannerIdSecond, slotId, userGroupId)

	id, err = GetBannerId(ctx, *s, slotId, userGroupId)
	fmt.Println(id)
	require.NoError(t, err)
}
