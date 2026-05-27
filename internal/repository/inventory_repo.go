package repository

import (
	"context"
	"fmt"
	"math"
	"server/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryRepository struct {
	DB *pgxpool.Pool
}

func NewInventoryRepository(db *pgxpool.Pool) *InventoryRepository {
	return &InventoryRepository{DB: db}
}

func (r *InventoryRepository) CreateLog(ctx context.Context, log models.InventoryLog) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `INSERT INTO inventory_logs (product_id, quantity) VALUES ($1, $2)`
	_, err := r.DB.Exec(ctx, query, log.ProductID, log.Quantity)
	return err
}

func (r *InventoryRepository) GetAllLogs(ctx context.Context, filter models.LogFilter, page models.Pagination) (models.PaginatedLogs, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	where := []string{}
	args := []any{}
	argIdx := 1

	if filter.Updated != nil {
		where = append(where, fmt.Sprintf("updated = $%d", argIdx))
		args = append(args, *filter.Updated)
		argIdx++
	}
	if filter.ProductID != nil {
		where = append(where, fmt.Sprintf("product_id = $%d", argIdx))
		args = append(args, *filter.ProductID)
		argIdx++
	}
	if filter.DateFrom != nil {
		where = append(where, fmt.Sprintf("created_at >= $%d", argIdx))
		args = append(args, *filter.DateFrom)
		argIdx++
	}
	if filter.DateTo != nil {
		where = append(where, fmt.Sprintf("created_at <= $%d", argIdx))
		args = append(args, *filter.DateTo)
		argIdx++
	}

	whereClause := ""
	if len(where) > 0 {
		whereClause = " WHERE "
		for i, clause := range where {
			if i > 0 {
				whereClause += " AND "
			}
			whereClause += clause
		}
	}

	countQuery := `SELECT COUNT(*) FROM inventory_logs` + whereClause
	var totalCount int
	if err := r.DB.QueryRow(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		return models.PaginatedLogs{}, err
	}

	offset := (page.Page - 1) * page.PageSize
	dataArgs := append(args, page.PageSize, offset)
	dataQuery := fmt.Sprintf(
		`SELECT id, product_id, quantity, updated, created_at FROM inventory_logs%s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
		whereClause, argIdx, argIdx+1,
	)

	rows, err := r.DB.Query(ctx, dataQuery, dataArgs...)
	if err != nil {
		return models.PaginatedLogs{}, err
	}
	defer rows.Close()

	logs := []models.InventoryLog{}
	for rows.Next() {
		var log models.InventoryLog
		if err := rows.Scan(&log.ID, &log.ProductID, &log.Quantity, &log.Updated, &log.CreatedAt); err != nil {
			return models.PaginatedLogs{}, err
		}
		logs = append(logs, log)
	}
	if err := rows.Err(); err != nil {
		return models.PaginatedLogs{}, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(page.PageSize)))

	return models.PaginatedLogs{
		Data:       logs,
		Page:       page.Page,
		PageSize:   page.PageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}


func (r *InventoryRepository) MarkLogsUpdated(ctx context.Context, ids []int) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var query string
	var args []any

	if len(ids) == 0 {
		query = `UPDATE inventory_logs SET updated = true WHERE updated = false`
	} else {
		query = `UPDATE inventory_logs SET updated = true WHERE id = ANY($1)`
		args = append(args, ids)
	}

	_, err := r.DB.Exec(ctx, query, args...)
	return err
}