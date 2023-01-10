package category

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

var table = "category"

type CategoryStorage struct {
	c *pgxpool.Pool
}

func NewCategoryStorage(c *pgxpool.Pool) *CategoryStorage {
	return &CategoryStorage{c: c}
}

func (c *CategoryStorage) GetID(ctx context.Context, CategoryName string) (int, error) {
	var id int
	query, args, err := squirrel.Select("id").From(table).Where(squirrel.Eq{"title": CategoryName}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return 0, err
	}

	row := c.c.QueryRow(ctx, query, args...)

	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (c *CategoryStorage) Create(ctx context.Context, m map[string]interface{}) (int, error) {
	var id int
	query, args, err := squirrel.Insert(table).SetMap(m).PlaceholderFormat(squirrel.Dollar).Suffix("RETURNING id").ToSql()
	if err != nil {
		return 0, err
	}

	row := c.c.QueryRow(ctx, query, args...)
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (c *CategoryStorage) Update(ctx context.Context, m map[string]interface{}, CatID int) (int, error) {
	var id int
	query, args, err := squirrel.Update(table).Where(squirrel.Eq{"id": CatID}).PlaceholderFormat(squirrel.Dollar).Suffix("RETURNING id").SetMap(m).ToSql()
	if err != nil {
		return 0, err
	}

	row := c.c.QueryRow(ctx, query, args...)
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (c *CategoryStorage) Delete(ctx context.Context, CatID int) error {
	query, args, err := squirrel.Delete(table).Where(squirrel.Eq{"id": CatID}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}
	_, err = c.c.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (c *CategoryStorage) GetOne(ctx context.Context, CatID int) (map[string]interface{}, error) {
	var cat Category

	query, args, err := squirrel.Select("title_ru", "title_eng", "description", "user_id").
		Where(squirrel.Eq{"id": CatID}).From(table).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	row := c.c.QueryRow(ctx, query, args...)
	err = row.Scan(&cat.TitleRu, &cat.TitleEng, &cat.Description, &cat.UserID)
	if err != nil {
		return nil, err
	}

	cat.ID = CatID
	return cat.ToMap(), nil
}

func (c *CategoryStorage) GetAll(ctx context.Context, UserID int) ([]map[string]interface{}, error) {

	query, args, err := squirrel.Select("id", "title_ru", "title_eng", "description", "user_id").
		Where(squirrel.Eq{"user_id": UserID}).
		PlaceholderFormat(squirrel.Dollar).From(table).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := c.c.Query(ctx, query, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var arrayMap []map[string]interface{}

	for rows.Next() {
		var cat Category

		err = rows.Scan(&cat.ID, &cat.TitleRu, &cat.TitleEng, &cat.Description, &cat.UserID)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		log.Printf("CatID: %v", cat.ID)
		arrayMap = append(arrayMap, cat.ToMap())

	}

	return arrayMap, nil
}
