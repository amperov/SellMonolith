package product

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
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

func (p *ProductStorage) SearchByUniqueCode(ctx context.Context, UniqueCode string) ([]map[string]interface{}, bool, error) {
	var input ProdForClient
	log.Println("Searching in ", transactionTable, "for unique code: ", UniqueCode)
	var arrayMap []map[string]interface{}

	query, args, err := squirrel.Select("id", "content_key", "category_name", "subcategory_name", "date_check").
		From(transactionTable).
		Where(squirrel.Eq{"unique_code": UniqueCode}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		log.Println(err)
		return nil, false, err
	}

	rows, err := p.c.Query(ctx, query, args...)
	if err != nil {
		log.Println(err)
		return nil, false, err
	}

	for rows.Next() {
		err = rows.Scan(&input.ID, &input.Content, &input.Category, &input.Subcategory, &input.DateCheck)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("TransactID: %v", input.ID)
		arrayMap = append(arrayMap, input.ToMap())
	}

	//////
	var exists bool
	query, args, err = squirrel.Select("id", "content_key", "category_name", "subcategory_name", "date_check").
		From(transactionTable).
		Prefix("SELECT EXISTS(").Suffix(")").
		Where(squirrel.Eq{"unique_code": UniqueCode}).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		log.Println(err)
		return nil, false, err
	}
	row := p.c.QueryRow(ctx, query, args...)
	err = row.Scan(&exists)
	if err != nil {
		return nil, exists, err
	}

	return arrayMap, exists, err
}

func (p *ProductStorage) Create(ctx context.Context, m map[string]interface{}) (int, error) {
	var id int
	query, args, err := squirrel.Insert(prodTable).PlaceholderFormat(squirrel.Dollar).SetMap(m).Suffix("RETURNING id").ToSql()
	if err != nil {
		return 0, err
	}
	row := p.c.QueryRow(ctx, query, args...)

	err = row.Scan(&id)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return id, nil
}

func (p *ProductStorage) Update(ctx context.Context, m map[string]interface{}, ProdID int) (int, error) {
	var id int
	query, args, err := squirrel.Update(prodTable).PlaceholderFormat(squirrel.Dollar).SetMap(m).Suffix("RETURNING id").Where(squirrel.Eq{"id": ProdID}).ToSql()
	if err != nil {
		return 0, err
	}
	row := p.c.QueryRow(ctx, query, args...)

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
		log.Println(err)
		return nil, err
	}
	rows, err := p.c.Query(ctx, query, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for rows.Next() {
		i.SubcategoryID = SubCatID
		err := rows.Scan(&i.ID, &i.Content)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		m = append(m, i.ToMapForSeller())
	}
	return m, nil
}

func (p *ProductStorage) GetCount(ctx context.Context, SubCatID int) (int, error) {
	var count int
	query := fmt.Sprintf("SELECT count(id) FROM %s WHERE subcategory_id=$1", prodTable)

	row := p.c.QueryRow(ctx, query, SubCatID)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (p *ProductStorage) GetSomeProducts(ctx context.Context, SubcatID int, Count int) ([]map[string]interface{}, error) {
	var m []map[string]interface{}
	var i ProdForClient

	query, args, err := squirrel.Select("content_key", "id").PlaceholderFormat(squirrel.Dollar).From(prodTable).
		Where(squirrel.Eq{"subcategory_id": SubcatID}).Suffix(fmt.Sprintf("LIMIT %d", Count)).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := p.c.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(&i.Content, &i.ID)
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
func (p *ProductStorage) Check(ctx context.Context, ItemID int) (bool, error) {
	var exists bool

	query, args, err := squirrel.Select("id").Prefix("SELECT EXISTS(").Suffix(")").From(prodTable).Where(squirrel.Eq{"subcategory_id": id}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		log.Println(err)
		return false, err
	}

	row := p.c.QueryRow(ctx, query, args...)
	err = row.Scan(&exists)
	if err != nil {
		log.Printf("scan: %v", err)
		return false, err
	}
	return exists, nil
}
