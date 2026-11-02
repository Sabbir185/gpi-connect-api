package designation

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/Sabbir185/gpi/entity"
)

type DesignationRepository interface {
	Search(ctx context.Context, query string, limit int) ([]entity.Designation, error)
	FindByID(ctx context.Context, id int64) (*entity.Designation, error)
}

type designationRepository struct {
	db *sql.DB
}

func NewDesignationRepository(db *sql.DB) DesignationRepository {
	return &designationRepository{db: db}
}

var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

// Search ranks titles that start with the query first, then any other match.
func (r *designationRepository) Search(ctx context.Context, query string, limit int) ([]entity.Designation, error) {
	pattern := likeEscaper.Replace(strings.ToLower(query))
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, hint, created_at, updated_at
		FROM designations
		WHERE LOWER(title) LIKE '%' || $1 || '%'
		ORDER BY (LOWER(title) LIKE $1 || '%') DESC, title
		LIMIT $2;
	`, pattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var designations []entity.Designation
	for rows.Next() {
		var d entity.Designation
		if err := rows.Scan(&d.ID, &d.Title, &d.Hint, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		designations = append(designations, d)
	}
	return designations, rows.Err()
}

func (r *designationRepository) FindByID(ctx context.Context, id int64) (*entity.Designation, error) {
	var d entity.Designation
	err := r.db.QueryRowContext(ctx, `
		SELECT id, title, hint, created_at, updated_at FROM designations WHERE id = $1;
	`, id).Scan(&d.ID, &d.Title, &d.Hint, &d.CreatedAt, &d.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}
