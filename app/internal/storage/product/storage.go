package product

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductStorage struct {
	c *pgxpool.Pool
}

func NewProductStorage(c *pgxpool.Pool) *ProductStorage {
	return &ProductStorage{c: c}
}

var (
	prodTable        = "products"
	transactionTable = "transactions"
)

func (p *ProductStorage) SearchByUniqueCode(ctx context.Context, UniqueCode string) (map[string]interface{}, bool) {
	var input ProdForClient
	var m map[string]interface{}
	var arrayMap []map[string]interface{}
	query, args, err := squirrel.Select("id", "content", "category", "subcategory", "date_check").From(transactionTable).Where(squirrel.Eq{"unique_code": UniqueCode}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, false
	}
	rows, err := p.c.Query(ctx, query, args...)
	if err != nil {
		return nil, false
	}
	for rows.Next() {
		err = rows.Scan(&input.ID, &input.Content, &input.Category, &input.Subcategory, &input.DateCheck)
		if err != nil {
			continue
		}
		arrayMap = append(arrayMap, input.ToMap())
	}
	m["products"] = arrayMap
	return m, true
}

func (p *ProductStorage) Create(ctx context.Context, m map[string]interface{}) (int, error) {
	var id int
	query, args, err := squirrel.Insert(prodTable).PlaceholderFormat(squirrel.Dollar).SetMap(m).Suffix("RETURNING id").ToSql()
	if err != nil {
		return 0, err
	}
	row := p.c.QueryRow(ctx, query, args)

	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (p *ProductStorage) Update(ctx context.Context, m map[string]interface{}, ProdID int) (int, error) {
	var id int
	query, args, err := squirrel.Update(prodTable).PlaceholderFormat(squirrel.Dollar).SetMap(m).Suffix("RETURNING id").ToSql()
	if err != nil {
		return 0, err
	}
	row := p.c.QueryRow(ctx, query, args)

	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (p *ProductStorage) Delete(ctx context.Context, SubCatID, ProdID int) error {
	query, args, err := squirrel.Delete(prodTable).PlaceholderFormat(squirrel.Dollar).Where(squirrel.Eq{"subcategory_id": SubCatID, "id": ProdID}).ToSql()
	if err != nil {
		return err
	}
	_, err = p.c.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductStorage) GetAll(ctx context.Context, SubCatID int) ([]map[string]interface{}, error) {
	var m []map[string]interface{}
	var i ProdForClient
	query, args, err := squirrel.Select("id", "content_key").From(prodTable).Where(squirrel.Eq{"subcategory_id": SubCatID}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := p.c.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err := rows.Scan(&i.ID, &i.Content)
		if err != nil {
			return nil, err
		}
		m = append(m, i.ToMapForSeller())
	}
	return m, nil
}

func (p *ProductStorage) GetCount(ctx context.Context, SubCatID int) (int, error) {
	var count int
	query, args, err := squirrel.Select("id").From(prodTable).PlaceholderFormat(squirrel.Dollar).Prefix("SELECT count(").Suffix(")").Where(squirrel.Eq{"subcategory_id": SubCatID}).ToSql()
	if err != nil {
		return 0, err
	}

	row := p.c.QueryRow(ctx, query, args...)
	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (p *ProductStorage) GetSomeProducts(ctx context.Context, SubcatID int, Count int) ([]map[string]interface{}, error) {
	var m []map[string]interface{}
	var i ProdForClient
	query, args, err := squirrel.Select("id", "content_key").PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"subcategory_id": SubcatID}).Suffix(fmt.Sprintf("LIMIT %d", Count)).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := p.c.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(&i.ID, &i.Content)
		if err != nil {
			return nil, err
		}
		i.SubcategoryID = SubcatID

		m = append(m, i.ToMapForSeller())
	}

	return m, nil
}

func (p *ProductStorage) DeleteOne(ctx context.Context, ProdID int) error {
	query, args, err := squirrel.Delete(prodTable).Where(squirrel.Eq{"id": ProdID}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}
	_, err = p.c.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}