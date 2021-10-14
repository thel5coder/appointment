package actions

import (
	"database/sql"
	"profira-backend/db/models"
	"profira-backend/db/repositories/contracts"
)

type FileRepository struct {
	DB *sql.DB
}


func NewFileRepository(DB *sql.DB) contracts.IFileRepository {
	return &FileRepository{DB: DB}
}

const (
	fileSelectStatement = `select id,path,created_at,updated_at`
)

func (repository FileRepository) scanRow(row *sql.Row) (res models.File, err error) {
	err = row.Scan(&res.ID, &res.Path, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (repository FileRepository) ReadBy(column, value, operator string) (data models.File, err error) {
	statement := fileSelectStatement + ` from "files" where ` + column + `` + operator + `$1`

	row := repository.DB.QueryRow(statement, value)
	data, err = repository.scanRow(row)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (repository FileRepository) Add(model models.File) (res string, err error) {
	statement := `insert into "files" (path,created_at,updated_at) values($1,$2,$3) returning id`

	err = repository.DB.QueryRow(statement, model.Path, model.CreatedAt, model.UpdatedAt).Scan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (repository FileRepository) Delete(model models.File) (res string, err error) {
	statement := `update "files" set updated_at=$1, deleted_at=$2 where id=$1 returning id`
	err = repository.DB.QueryRow(statement, model.UpdatedAt, model.DeletedAt, model.ID).Scan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (FileRepository) CountBy(column, value string) (res int, err error) {
	panic("implement me")
}
