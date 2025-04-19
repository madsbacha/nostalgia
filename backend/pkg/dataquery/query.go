package dataquery

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Scanner[T any] interface {
	GetSelectProperties() []string
	Scan(s func(dest ...interface{}) error) (T, error)
}

type QueryBuilder[T any] struct {
	ctx            context.Context
	scanner        Scanner[T]
	db             *sql.DB
	tableName      string
	WhereString    string
	WhereParams    []interface{}
	HasWhere       bool
	hasOrdering    bool
	orderBy        string
	orderDirection string
	hasLimit       bool
	limit          int
	hasOffset      bool
	offset         int
	postfixParams  []interface{}
}

type QueryResult[T any] struct {
	items []T
}

func (qr *QueryResult[T]) First() T {
	return qr.items[0]
}

func (qb *QueryBuilder[T]) ScanRow(row *sql.Row) (QueryResult[T], error) {
	item, err := qb.scanner.Scan(row.Scan)
	if err != nil {
		return QueryResult[T]{}, err
	}

	qr := QueryResult[T]{
		items: []T{item},
	}

	return qr, nil
}

func (qb *QueryBuilder[T]) ScanRows(row *sql.Rows) (QueryResult[T], error) {
	itemList := make([]T, 0)

	for row.Next() {
		item, err := qb.scanner.Scan(row.Scan)
		if err != nil {
			return QueryResult[T]{}, err
		}
		itemList = append(itemList, item)
	}

	qr := QueryResult[T]{
		items: itemList,
	}

	return qr, nil
}

func QueryContext[T any](ctx context.Context, db *sql.DB, tableName string, scanner Scanner[T]) *QueryBuilder[T] {
	return &QueryBuilder[T]{
		ctx:            ctx,
		scanner:        scanner,
		db:             db,
		tableName:      tableName,
		WhereString:    "",
		WhereParams:    []interface{}{},
		HasWhere:       false,
		hasOrdering:    false,
		orderBy:        "",
		orderDirection: "ASC",
		hasLimit:       false,
		limit:          0,
		hasOffset:      false,
		offset:         0,
		postfixParams:  []interface{}{},
	}
}

func Query[T any](db *sql.DB, tableName string, scanner Scanner[T]) *QueryBuilder[T] {
	return QueryContext(context.Background(), db, tableName, scanner)
}

func (qb *QueryBuilder[T]) GenerateSelectQuery() string {
	el := strings.Join(qb.scanner.GetSelectProperties(), ",")

	query := "SELECT " + el + " FROM " + qb.tableName

	if qb.HasWhere {
		query += " WHERE " + qb.WhereString
	}

	return query
}

func (qb *QueryBuilder[T]) GeneratePostfixQuery() string {
	query := ""
	var params []interface{}

	if qb.hasOrdering {
		query += " ORDER BY " + qb.orderBy + " " + qb.orderDirection
	}

	if qb.hasLimit {
		query += " LIMIT ?"
		params = append(params, qb.limit)
	}

	if qb.hasOffset {
		query += " OFFSET ?"
		params = append(params, qb.offset)
	}

	qb.postfixParams = params

	return query
}

func (qb *QueryBuilder[T]) Where(where string, params ...interface{}) *QueryBuilder[T] {
	qb.WhereString = where
	qb.WhereParams = params
	qb.HasWhere = true

	return qb
}

func (qb *QueryBuilder[T]) OrderByAsc(col string) *QueryBuilder[T] {
	qb.hasOrdering = true
	qb.orderBy = col
	qb.orderDirection = "ASC"
	return qb
}

func (qb *QueryBuilder[T]) OrderByDesc(col string) *QueryBuilder[T] {
	qb.hasOrdering = true
	qb.orderBy = col
	qb.orderDirection = "DESC"
	return qb
}

func (qb *QueryBuilder[T]) Limit(limit int) *QueryBuilder[T] {
	qb.hasLimit = true
	qb.limit = limit
	return qb
}

func (qb *QueryBuilder[T]) Offset(offset int) *QueryBuilder[T] {
	qb.hasOffset = true
	qb.offset = offset
	return qb
}

func (qb *QueryBuilder[T]) First() (T, error) {
	query := qb.GenerateSelectQuery()
	query += qb.GeneratePostfixQuery()
	var params []interface{}
	params = append(params, qb.WhereParams...)
	params = append(params, qb.postfixParams...)
	row := qb.db.QueryRow(query, params...)

	qr, err := qb.ScanRow(row)
	if err != nil {
		return *new(T), err
	}
	return qr.First(), nil
}

func (qb *QueryBuilder[T]) Exists() (bool, error) {
	if !qb.HasWhere {
		return false, fmt.Errorf("a where statement is required, when using exists")
	}

	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s)", qb.tableName, qb.WhereString)

	var exists bool
	err := qb.db.QueryRow(query, qb.WhereParams...).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (qb *QueryBuilder[T]) GetAll() ([]T, error) {
	query := qb.GenerateSelectQuery()
	query += qb.GeneratePostfixQuery()
	var params []interface{}
	params = append(params, qb.WhereParams...)
	params = append(params, qb.postfixParams...)

	rows, err := qb.db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	qr, err := qb.ScanRows(rows)
	if err != nil {
		return nil, err
	}
	return qr.items, nil
}

func (qb *QueryBuilder[T]) Insert(props map[string]interface{}) (sql.Result, error) {
	query := "INSERT INTO " + qb.tableName
	keys := make([]string, 0)
	questionMark := make([]string, 0)
	values := make([]interface{}, 0)
	for key, val := range props {
		keys = append(keys, key)
		questionMark = append(questionMark, "?")
		values = append(values, val)
	}
	query += "(" + strings.Join(keys, ",") + ") values (" + strings.Join(questionMark, ",") + ")"

	stmt, err := qb.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	return stmt.Exec(values...)
}

func (qb *QueryBuilder[T]) Update(props map[string]interface{}) (sql.Result, error) {
	query := "UPDATE " + qb.tableName + " SET "

	keys := make([]string, 0)
	values := make([]interface{}, 0)
	for key, val := range props {
		keys = append(keys, key+" = ?")
		values = append(values, val)
	}
	query += strings.Join(keys, ", ")

	if qb.HasWhere {
		query += " WHERE " + qb.WhereString
		values = append(values, qb.WhereParams...)
	}

	stmt, err := qb.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	return stmt.Exec(values...)
}

func (qb *QueryBuilder[T]) AddToColumn(column string, val int) (sql.Result, error) {
	query := "UPDATE " + qb.tableName + " SET " + column + " = " + column + " + ?"

	values := make([]interface{}, 1)
	values[0] = val

	if qb.HasWhere {
		query += " WHERE " + qb.WhereString
		values = append(values, qb.WhereParams...)
	}

	stmt, err := qb.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	return stmt.Exec(values...)
}

func (qb *QueryBuilder[T]) Delete() (sql.Result, error) {
	if !qb.HasWhere {
		return nil, fmt.Errorf("A where clause is required when deleting")
	}
	query := "DELETE FROM " + qb.tableName + " WHERE " + qb.WhereString

	stmt, err := qb.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	return stmt.Exec(qb.WhereParams...)
}

type ColumnQuery[T any] struct {
	QueryBuilder *QueryBuilder[T]
	ColumnName   string
}

func (qb *QueryBuilder[T]) Column(name string) *ColumnQuery[T] {
	return &ColumnQuery[T]{
		QueryBuilder: qb,
		ColumnName:   name,
	}
}

func (cq *ColumnQuery[T]) GetAllStrings() ([]string, error) {
	query := "SELECT " + cq.ColumnName + " FROM " + cq.QueryBuilder.tableName

	if cq.QueryBuilder.HasWhere {
		query += " WHERE " + cq.QueryBuilder.WhereString
	}

	rows, err := cq.QueryBuilder.db.Query(query, cq.QueryBuilder.WhereParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resultList := make([]string, 0)

	for rows.Next() {
		var rowString string
		err := rows.Scan(&rowString)
		if err != nil {
			return nil, err
		}
		resultList = append(resultList, rowString)
	}

	return resultList, nil
}
