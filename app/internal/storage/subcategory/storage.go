package subcategory

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"log"
)

var table = "subcategory"

type SubcategoryStorage struct {
	c *pgxpool.Pool
}

func NewSubcategoryStorage(c *pgxpool.Pool) *SubcategoryStorage {
	return &SubcategoryStorage{c: c}
}

func (c *SubcategoryStorage) GetAll(ctx context.Context, CatID int) ([]map[string]interface{}, error) {
	var subcats []Subcategory
	query, args, err := squirrel.Select("id", "title", "subitem_id").Where(squirrel.Eq{"category_id": CatID}).PlaceholderFormat(squirrel.Dollar).From(table).ToSql()
	if err != nil {
		log.Printf("error make query: %v", err)
		return nil, err
	}

	rows, err := c.c.Query(ctx, query, args...)
	if err != nil {
		log.Printf("error query: %v", err)
		return nil, err
	}
	var arrayMap []map[string]interface{}
	for rows.Next() {
		var cat Subcategory

		err = rows.Scan(&cat.ID, &cat.Title, &cat.SubItemID)
		if err != nil {
			log.Printf("error scan: %v", err)
			return nil, err
		}
		cat.CategoryID = CatID
		subcats = append(subcats, cat)
	}
	subcategories := Sort(subcats)

	for _, subcategory := range subcategories {
		arrayMap = append(arrayMap, subcategory.ToMap())
	}

	return arrayMap, nil
}

func Sort(subcats []Subcategory) []Subcategory {
	for i := 0; i < len(subcats)-1; i++ {
		if subcats[i].Title > subcats[i+1].Title {
			subcats[i+1].Title = subcats[i].Title
			subcats[i].Title = subcats[i+1].Title
		}
	}
	return subcats
}

func (c *SubcategoryStorage) GetCount(ctx context.Context, CatID int) (int, error) {
	var count int
	query := fmt.Sprintf("SELECT count(id) FROM %s WHERE category_id=$1", table)

	row := c.c.QueryRow(ctx, query, CatID)

	err := row.Scan(&count)
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
	query, args, err := squirrel.Update(table).Where(squirrel.Eq{"id": SubCatID}).
		PlaceholderFormat(squirrel.Dollar).Suffix("RETURNING id").SetMap(m).ToSql()
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

	query, args, err := squirrel.Select("title", "category_id", "subitem_id").Where(squirrel.Eq{"id": SubCatID}).From(table).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	row := c.c.QueryRow(ctx, query, args...)
	err = row.Scan(&cat.Title, &cat.CategoryID, &cat.SubItemID)
	if err != nil {
		return nil, err
	}

	cat.ID = SubCatID
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
		logrus.Debugf("Subcat error: %v", err)
		return 0, err
	}
	return id, nil
}
