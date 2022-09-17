package sqlstorage

import (
	"context"
	"fmt"
	"github.com/EkaterinaShamanaeva/Banner-rotation-service/internal/storage"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type Storage struct {
	Pool *pgxpool.Pool
}

func New() *Storage {
	return &Storage{Pool: nil}
}

// Connect - подключение к базе данных.
func Connect(ctx context.Context, dsn string) (dbPool *pgxpool.Pool, err error) {
	// проверка url к Postgresql
	conn, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		err = fmt.Errorf("failed to parse pg config: %w", err)
		return
	}
	// конфигурация
	conn.MaxConns = int32(5)
	conn.MinConns = int32(1)
	conn.HealthCheckPeriod = 1 * time.Minute
	conn.MaxConnLifetime = 24 * time.Hour
	conn.MaxConnIdleTime = 30 * time.Minute
	conn.ConnConfig.ConnectTimeout = 1 * time.Second
	// пул коннектов
	dbPool, err = pgxpool.ConnectConfig(ctx, conn) // pool
	if err != nil {
		err = fmt.Errorf("failed to connect config: %w", err)
		return
	}
	return
}

// Close - закрытие пула коннектов к базе данных.
func (s *Storage) Close(ctx context.Context) error {
	s.Pool.Close()
	return nil
}

// CreateBanner - создание нового баннера.
func (s *Storage) CreateBanner(ctx context.Context, description string) (uuid.UUID, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}
	query := `INSERT INTO banners(id, description)
				VALUES($1, $2) RETURNING id; `
	_, err = s.Pool.Exec(ctx, query, id, description)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

// CreateSlot - создание нового слота.
func (s *Storage) CreateSlot(ctx context.Context, description string) (uuid.UUID, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}
	query := `INSERT INTO slots(id, description)
				VALUES($1, $2) RETURNING id; `
	_, err = s.Pool.Exec(ctx, query, id, description)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

// CreateUserGroup - создание новой группы пользователей.
func (s *Storage) CreateUserGroup(ctx context.Context, description string) (uuid.UUID, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}
	query := `INSERT INTO user_groups(id, description)
				VALUES($1, $2) RETURNING id; `
	_, err = s.Pool.Exec(ctx, query, id, description)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

// AddBannerToSlot - добавить баннер в слот.
func (s *Storage) AddBannerToSlot(ctx context.Context, slotId uuid.UUID, bannerId uuid.UUID) error {
	query := `INSERT INTO banner_slot(slot_id, banner_id)
				VALUES($1, $2); `
	_, err := s.Pool.Exec(ctx, query, slotId, bannerId)
	if err != nil {
		return err
	}
	return nil
}

// DeleteBannerFromSlot - удалить баннер из слота.
func (s *Storage) DeleteBannerFromSlot(ctx context.Context, slotId uuid.UUID, bannerId uuid.UUID) error {
	query := `DELETE from banner_slot WHERE slot_id=$1 AND banner_id=$2); `
	_, err := s.Pool.Exec(ctx, query, slotId, bannerId)
	if err != nil {
		return err
	}
	return nil
}

// GetBannersFromSlot - получить все баннеры в ротации в слоте.
func (s *Storage) GetBannersFromSlot(ctx context.Context, slotId uuid.UUID) ([]uuid.UUID, error) {
	query := `SELECT banner_id FROM banner_slot WHERE slot_id=$1;`
	var bannersByQuery []*uuid.UUID
	err := pgxscan.Select(ctx, s.Pool, &bannersByQuery, query, slotId)
	if err != nil {
		return nil, err
	}
	bannersRes := make([]uuid.UUID, 0, len(bannersByQuery))
	for _, banner := range bannersByQuery {
		bannersRes = append(bannersRes, *banner)
	}
	return bannersRes, nil
}

func (s *Storage) ClickBanner(ctx context.Context, bannerId, slotId, userGroupId uuid.UUID) error {

	return nil
}

// GetBannerStatistic - получить статистику кликов и показов баннера в слоте.
// TODO изменить на массив или установить лимит на 1
func (s *Storage) GetBannerStatistic(ctx context.Context, bannerId uuid.UUID, userGroupId uuid.UUID) (storage.BannerStatistic, error) {
	emptyBannerStatistic := storage.BannerStatistic{
		BannerId:    bannerId,
		UserGroupId: userGroupId,
	}
	query := `SELECT * FROM banner_statistic WHERE banner_id=$1 AND user_group_id = $2;`
	var bannerByQuery []*storage.BannerStatistic
	err := pgxscan.Select(ctx, s.Pool, &bannerByQuery, query, bannerId, userGroupId)
	if err != nil {
		return emptyBannerStatistic, err
	}
	return *bannerByQuery[0], nil
}
