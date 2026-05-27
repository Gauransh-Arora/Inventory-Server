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

func (r *ProductRepository) CreateProduct(ctx context.Context, p models.Product) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `INSERT INTO products (icode, item_name, batch_no, mrp, barcode) VALUES ($1,$2,$3,$4,$5)`
	_, err := r.DB.Exec(ctx, query, p.ICode, p.ItemName, p.BatchNo, p.MRP, p.Barcode)
	return err
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
		`SELECT id, icode, item_name, batch_no, mrp, barcode FROM products%s ORDER BY id DESC LIMIT $%d OFFSET $%d`,
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