package news

import (
	"context"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type Record struct {
	bun.BaseModel `bun:"table:news"`
	ID            uuid.UUID `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	Author        string    `bun:"author,nullzero,notnull"`
	Title         string    `bun:"title,nullzero,notnull"`
	Summary       string    `bun:"summary,nullzero,notnull"`
	Content       string    `bun:"content,nullzero,notnull"`
	Source        string    `bun:"source,nullzero,notnull"`
	Tags          []string  `bun:"tags,nullzero,notnull,array"`
	CreatedAt     time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	DeletedAt     time.Time `bun:"deleted_at,nullzero,soft_delete"`
}

func (s Store) Create(ctx context.Context, news Record) (createdNews Record, err error) {
	news.ID = uuid.New()
	err = s.db.NewInsert().Model(&news).Returning("*").Scan(ctx, &createdNews)
	if err != nil {
		return createdNews, err
	}
	return createdNews, nil
}

func (s Store) FindByID(ctx context.Context, id uuid.UUID) (news Record, err error) {
	err = s.db.NewSelect().Model(&news).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return news, err
	}
	return news, nil
}

func (s Store) FindAll(ctx context.Context) (news []Record, err error) {
	err = s.db.NewSelect().Model(&news).Scan(ctx, &news)
	if err != nil {
		return news, err
	}
	return news, nil
}

func (s Store) DeleteByID(ctx context.Context, id uuid.UUID) (err error) {
	_, err = s.db.NewDelete().Model(&Record{}).Where("id = ?", id).Returning("NULL").Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s Store) UpdateByID(ctx context.Context, id uuid.UUID, news Record) (err error) {
	_, err = s.db.NewUpdate().Model(&news).Where("id = ?", id).Returning("NULL").Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
