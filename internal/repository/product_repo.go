package repository

import (
	"context"
	"fmt"
	"math"
	"server/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	DB *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, p []models.Product) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Check if product already exists
	for _, prod := range p {
		query := `SELECT id FROM products WHERE icode = $1 AND item_name = $2 AND batch_no = $3 AND mrp = $4 AND barcode = $5 LIMIT 1`
		id := -1
		err := tx.QueryRow(ctx, query, prod.ICode, prod.ItemName, prod.BatchNo, prod.MRP, prod.Barcode).Scan(&id)

		if err == nil {
			return fmt.Errorf("product already exists (icode: %d, name: %s)", prod.ICode, prod.ItemName)
		} else if err.Error() != "no rows in result set" {
			return err
		}
	}

	// Create products
	for _, prod := range p {
		query := `INSERT INTO products (icode, item_name, batch_no, mrp, barcode) VALUES ($1,$2,$3,$4,$5)`
		_, err := tx.Exec(ctx, query, prod.ICode, prod.ItemName, prod.BatchNo, prod.MRP, prod.Barcode)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *ProductRepository) GetAllProducts(ctx context.Context, icode *int, page models.Pagination) (models.PaginatedProducts, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	where := []string{}
	args := []any{}
	argIdx := 1

	if icode != nil {
		where = append(where, fmt.Sprintf("icode = $%d", argIdx))
		args = append(args, *icode)
		argIdx++
	}

	whereClause := ""
	if len(where) > 0 {
		whereClause = " WHERE " + where[0]
	}

	countQuery := `SELECT COUNT(*) FROM products` + whereClause
	var totalCount int
	if err := r.DB.QueryRow(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		return models.PaginatedProducts{}, err
	}

	offset := (page.Page - 1) * page.PageSize
	dataArgs := append(args, page.PageSize, offset)
	dataQuery := fmt.Sprintf(
		`SELECT id, icode, item_name, batch_no, mrp, barcode FROM products%s ORDER BY id LIMIT $%d OFFSET $%d`,
		whereClause, argIdx, argIdx+1,
	)

	rows, err := r.DB.Query(ctx, dataQuery, dataArgs...)
	if err != nil {
		return models.PaginatedProducts{}, err
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.ICode, &p.ItemName, &p.BatchNo, &p.MRP, &p.Barcode); err != nil {
			return models.PaginatedProducts{}, err
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return models.PaginatedProducts{}, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(page.PageSize)))

	return models.PaginatedProducts{
		Data:       products,
		Page:       page.Page,
		PageSize:   page.PageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

func (r *ProductRepository) GetProductByBarcode(ctx context.Context, barcode string) ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `SELECT id, icode, item_name, batch_no, mrp, barcode FROM products WHERE barcode = $1`
	rows, err := r.DB.Query(ctx, query, barcode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.ICode, &p.ItemName, &p.BatchNo, &p.MRP, &p.Barcode); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, id int, data models.Product) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `UPDATE products SET icode=$1, item_name=$2, batch_no=$3, mrp=$4, barcode=$5 where id = $6`
	_, err := r.DB.Exec(ctx, query, data.ICode, data.ItemName, data.BatchNo, data.MRP, data.Barcode, id)
	return err
}

func (r *ProductRepository) DeleteProducts(ctx context.Context, ids []int) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, id := range ids {
		query := `DELETE FROM products WHERE id = $1`
		_, err := tx.Exec(ctx, query, id)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
