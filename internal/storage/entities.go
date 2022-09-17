package storage

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type Storage interface {
	Close(context.Context) error
	Connect(context.Context, string) (*pgxpool.Pool, error)
	CreateSlot(context.Context, string) (uuid.UUID, error)
	CreateUserGroup(context.Context, string) (uuid.UUID, error)
	AddBannerToSlot(ctx context.Context, slotId uuid.UUID, bannerId uuid.UUID) error
	DeleteBannerFromSlot(ctx context.Context, slotId uuid.UUID, bannerId uuid.UUID) error
	GetBannersFromSlot(ctx context.Context, slotId uuid.UUID) ([]uuid.UUID, error)
	ClickBanner(ctx context.Context, bannerId, slotId, userGroupId uuid.UUID) error
}

// Banner - рекламный/информационный элемент, который показывается в слоте.
type Banner struct {
	Id          uuid.UUID `json:"id"`
	Description string    `json:"description"`
}

// Slot - место на сайте, на котором мы показываем баннер.
type Slot struct {
	Id          uuid.UUID `json:"id"`
	Description string    `json:"description"`
}

// UserGroup - это группа пользователей сайта со схожими интересами, например "девушки 20-25" или "дедушки 80+".
type UserGroup struct {
	Id          uuid.UUID `json:"id"`
	Description string    `json:"description"`
}

// Rotation - баннер в ротации в данном слоте.
type Rotation struct {
	SlotId   uuid.UUID `json:"slot_id"`
	BannerId uuid.UUID `json:"banner_id"`
}

// Event - событие по показу баннера или переходу.
type Event struct {
	BannerId    uuid.UUID `json:"banner_id"`
	UserGroupId uuid.UUID `json:"user_group_id"`
	Date        time.Time `json:"date"`
	Action      Action    `json:"action"`
}

// Action - переход или показ баннера
type Action string

const (
	click Action = "Click"
	show  Action = "Show"
)

// BannerStatistic - статистика по переходу и показу баннера.
type BannerStatistic struct {
	BannerId    uuid.UUID `json:"banner_id"`
	SlotId      uuid.UUID `json:"slot_id"`
	UserGroupId uuid.UUID `json:"user_group_id"`
	ShowCount   int       `json:"show_count"`
	ClickCount  int       `json:"click_count"`
}
