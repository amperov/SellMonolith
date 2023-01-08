package subcategory

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
)

var table = "subcategories"

type SubcategoryStorage struct {
	c *pgxpool.Pool
}

func NewSubcategoryStorage(c *pgxpool.Pool) *SubcategoryStorage {
	return &SubcategoryStorage{c: c}
}

func (c *SubcategoryStorage) GetAll(ctx context.Context, CatID int) ([]map[string]interface{}, error) {
	var cats []Subcategory

	query, args, err := squirrel.Select("title", "category_id").Where(squirrel.Eq{"user_id": CatID}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := c.c.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	var arrayMap []map[string]interface{}
	for rows.Next() {
		var cat Subcategory
		var um map[string]interface{}

		err = rows.Scan(&cat.Title, &cat.CategoryID)
		if err != nil {
			return nil, err
		}
		um[strconv.Itoa(cat.ID)] = cat.ToMap()
		arrayMap = append(arrayMap, um)
		cats = append(cats, cat)
	}

	return arrayMap, nil
}

func (c *SubcategoryStorage) GetCount(ctx context.Context, CatID int) (int, error) {
	var count int
	query, args, err := squirrel.Select("id").Prefix("SELECT EXISTS(").Suffix(")").
		From(table).Where(squirrel.Eq{"category_id": CatID}).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return 0, err
	}

	row := c.c.QueryRow(ctx, query, args...)

	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c *SubcategoryStorage) Create(ctx context.Context, m map[string]interface{}) (int, error) {
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

func (c *SubcategoryStorage) Update(ctx context.Context, m map[string]interface{}, SubCatID int) (int, error) {
	var id int
	query, args, err := squirrel.Update(table).Where(squirrel.Eq{"id": SubCatID}).PlaceholderFormat(squirrel.Dollar).Suffix("RETURNING id").SetMap(m).ToSql()
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

func (c *SubcategoryStorage) Delete(ctx context.Context, SubCatID int) error {
	query, args, err := squirrel.Delete(table).Where(squirrel.Eq{"id": SubCatID}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}
	_, err = c.c.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (c *SubcategoryStorage) Get(ctx context.Context, SubCatID int) (map[string]interface{}, error) {
	var cat Subcategory

	query, args, err := squirrel.Select("title", "title_eng", "description", "user_id").Where(squirrel.Eq{"id": SubCatID}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	row := c.c.QueryRow(ctx, query, args...)
	err = row.Scan(&cat.Title, &cat.CategoryID)
	if err != nil {
		return nil, err
	}

	return cat.ToMap(), nil
}

func (c *SubcategoryStorage) GetID(ctx context.Context, SubcategoryName string, CategoryID int) (int, error) {
	var id int
	query, args, err := squirrel.Select("id").From(table).
		Where(squirrel.Eq{"title": SubcategoryName, "category_id": CategoryID}).
		PlaceholderFormat(squirrel.Dollar).ToSql()
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
