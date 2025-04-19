package datab

import (
	"database/sql"
	"fmt"
	"strings"
)

type Repository[T any] interface {
	GetTableName() string
	GetSelectProperties() []string
	Db() *sql.DB
	Scan(s func(dest ...interface{}) error) (T, error)
}

type QueryBuilder[T any] struct {
	Repo           Repository[T]
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
	item, err := qb.Repo.Scan(row.Scan)
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
		item, err := qb.Repo.Scan(row.Scan)
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

func Query[T any](repo Repository[T]) *QueryBuilder[T] {
	return &QueryBuilder[T]{
		Repo:           repo,
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

func (qb *QueryBuilder[T]) GenerateSelectQuery() string {
	el := strings.Join(qb.Repo.GetSelectProperties(), ",")

	query := "SELECT " + el + " FROM " + qb.Repo.GetTableName()

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

func (qb *QueryBuilder[T]) GetFirst() (T, error) {
	query := qb.GenerateSelectQuery()
	query += qb.GeneratePostfixQuery()
	var params []interface{}
	params = append(params, qb.WhereParams...)
	params = append(params, qb.postfixParams...)
	row := qb.Repo.Db().QueryRow(query, params...)

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

	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s)", qb.Repo.GetTableName(), qb.WhereString)

	var exists bool
	err := qb.Repo.Db().QueryRow(query, qb.WhereParams...).Scan(&exists)
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

	rows, err := qb.Repo.Db().Query(query, params...)
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
	query := "INSERT INTO " + qb.Repo.GetTableName()
	keys := make([]string, 0)
	questionMark := make([]string, 0)
	values := make([]interface{}, 0)
	for key, val := range props {
		keys = append(keys, key)
		questionMark = append(questionMark, "?")
		values = append(values, val)
	}
	query += "(" + strings.Join(keys, ",") + ") values (" + strings.Join(questionMark, ",") + ")"

	stmt, err := qb.Repo.Db().Prepare(query)
	if err != nil {
		return nil, err
	}

	return stmt.Exec(values...)
}

func (qb *QueryBuilder[T]) Update(props map[string]interface{}) (sql.Result, error) {
	query := "UPDATE " + qb.Repo.GetTableName() + " SET "

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

	stmt, err := qb.Repo.Db().Prepare(query)
	if err != nil {
		return nil, err
	}

	return stmt.Exec(values...)
}

func (qb *QueryBuilder[T]) AddToColumn(column string, val int) (sql.Result, error) {
	query := "UPDATE " + qb.Repo.GetTableName() + " SET " + column + " = " + column + " + ?"

	values := make([]interface{}, 1)
	values[0] = val

	if qb.HasWhere {
		query += " WHERE " + qb.WhereString
		values = append(values, qb.WhereParams...)
	}

	stmt, err := qb.Repo.Db().Prepare(query)
	if err != nil {
		return nil, err
	}

	return stmt.Exec(values...)
}

func (qb *QueryBuilder[T]) Delete() (sql.Result, error) {
	if !qb.HasWhere {
		return nil, fmt.Errorf("A where clause is required when deleting")
	}
	query := "DELETE FROM " + qb.Repo.GetTableName() + " WHERE " + qb.WhereString

	stmt, err := qb.Repo.Db().Prepare(query)
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
	query := "SELECT " + cq.ColumnName + " FROM " + cq.QueryBuilder.Repo.GetTableName()

	if cq.QueryBuilder.HasWhere {
		query += " WHERE " + cq.QueryBuilder.WhereString
	}

	rows, err := cq.QueryBuilder.Repo.Db().Query(query, cq.QueryBuilder.WhereParams...)
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
