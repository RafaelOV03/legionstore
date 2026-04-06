package repository

import (
	"database/sql"
	"fmt"
	"strings"
)

// QueryBuilder facilita la construcción dinámica de queries SQL
type QueryBuilder struct {
	query      string
	args       []interface{}
	conditions []string
}

// NewQueryBuilder crea un nuevo constructor de queries
func NewQueryBuilder(baseQuery string) *QueryBuilder {
	return &QueryBuilder{
		query: baseQuery,
		args:  []interface{}{},
	}
}

// Where agrega una condición WHERE
func (qb *QueryBuilder) Where(condition string, args ...interface{}) *QueryBuilder {
	qb.conditions = append(qb.conditions, condition)
	qb.args = append(qb.args, args...)
	return qb
}

// Build construye la query final completa
func (qb *QueryBuilder) Build() (string, []interface{}) {
	finalQuery := qb.query
	if len(qb.conditions) > 0 {
		finalQuery += " WHERE " + strings.Join(qb.conditions, " AND ")
	}
	return finalQuery, qb.args
}

// ScanError es un wrapper para errores de scanning
func ScanError(err error, fieldName string) error {
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error scanning %s: %w", fieldName, err)
	}
	return err
}

// DBError es un wrapper para errores de base de datos
func DBError(err error, operation string) error {
	if err != nil {
		return fmt.Errorf("%s failed: %w", operation, err)
	}
	return nil
}

// CheckRowsAffected valida que se modificaron filas
func CheckRowsAffected(result sql.Result, operation string) error {
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("%s failed: no rows affected", operation)
	}
	return nil
}
