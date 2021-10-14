package actions

import (
	"database/sql"
	"gorm.io/gorm"
	"profira-backend/db/models"
	"profira-backend/db/repositories/contracts"
)

type CategoryRepository struct {
	DB     *sql.DB
	GormDB *gorm.DB
}

func NewCategoryRepository(DB *sql.DB) contracts.ICategoryRepository {
	return &CategoryRepository{DB: DB}
}

const (
	categorySelectStatement = `select c.id,c.slug,c.name,c.parent_id,c.is_active,c.file_icon_id,fi.path,c.file_background_id,fb.path,c.created_at,c.updated_at`
	categoryJoinStatement   = `left join files fi on fi.id = c.file_icon_id and fi.deleted_at is null
                             left join files fb on fb.id = c.file_background_id and fb.deleted_at is null`
)

func (repository CategoryRepository) scanRows(rows *sql.Rows) (res models.Category, err error) {
	err = rows.Scan(&res.ID, &res.Slug, &res.Name, &res.ParentID, &res.IsActive, &res.FileIconID, &res.FileIconPath, &res.FileBackgroundID, &res.FileBackgroundPath, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (repository CategoryRepository) scanRow(row *sql.Row) (res models.Category, err error) {
	err = row.Scan(&res.ID, &res.Slug, &res.Name, &res.ParentID, &res.IsActive, &res.FileIconID, &res.FileIconPath, &res.FileBackgroundID, &res.FileBackgroundPath, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (repository CategoryRepository) Browse(search, order, sort string, limit, offset int) (data []models.Category, count int, err error) {
	repository.GormDB.Find(&data)

	return data, count, nil
}

func (repository CategoryRepository) BrowseTree(search string) (data []models.Category, err error) {
	statement := `WITH RECURSIVE treatment_categories AS (
	SELECT categories.id,categories.slug,categories.name,categories.parent_id,categories.is_active,categories.file_icon_id,fi.path,
	categories.file_background_id,fb.path,categories.created_at,categories.updated_at from categories
	LEFT JOIN files fi on fi.id=categories.file_icon_id and fi.deleted_at is null
	left join files fb on fb.id = categories.file_background_id and fb.deleted_at is null 
	WHERE categories.parent_id IS NULL and categories.deleted_at is null
	UNION
	SELECT child.id,child.slug,child.name,child.parent_id,child.is_active,child.file_icon_id,fi.path,
	child.file_background_id,fb.path,child.created_at,child.updated_at from categories child
	INNER JOIN treatment_categories tc on tc.id = child.parent_id
	LEFT JOIN files fi on fi.id=child.file_icon_id and fi.deleted_at is null
	left join files fb on fb.id = child.file_background_id and fb.deleted_at is null 
    where child.deleted_at is null
) SELECT treatment_categories.id,treatment_categories.slug,treatment_categories.name,
treatment_categories.parent_id,treatment_categories.is_active,treatment_categories.file_icon_id,fi.path,
	treatment_categories.file_background_id,fb.path,treatment_categories.created_at,
	treatment_categories.updated_at,array_to_string(array_agg(mt.id ||':'|| mt.name),',') FROM treatment_categories 
	LEFT JOIN files fi on fi.id=treatment_categories.file_icon_id and fi.deleted_at is null
	left join files fb on fb.id = treatment_categories.file_background_id and fb.deleted_at is null
	left JOIN master_treatments mt on mt.category_id = treatment_categories.id
	GROUP BY treatment_categories.id,treatment_categories.slug,treatment_categories.name,fi.id,
	treatment_categories.parent_id,treatment_categories.is_active,treatment_categories.file_icon_id,fb.id,
	treatment_categories.file_background_id,treatment_categories.created_at,
	treatment_categories.updated_at`

	rows, err := repository.DB.Query(statement)
	if err != nil {
		return data, err
	}

	for rows.Next() {
		var temp models.Category
		err = rows.Scan(&temp.ID, &temp.Slug, &temp.Name, &temp.ParentID, &temp.IsActive, &temp.FileIconID, &temp.FileIconPath, &temp.FileBackgroundID,
			&temp.FileBackgroundPath, &temp.CreatedAt, &temp.UpdatedAt, &temp.Treatments)
		if err != nil {
			return data, err
		}
		data = append(data, temp)
	}

	return data, nil
}

func (repository CategoryRepository) BrowseByParentID(parentID string) (data []models.Category, err error) {
	statement := categorySelectStatement + ` from categories c ` + categoryJoinStatement + ` where c.parent_id =$1 and c.deleted_at is null order by c.name asc`
	rows, err := repository.DB.Query(statement, parentID)
	if err != nil {
		return data, err
	}

	for rows.Next() {
		temp, err := repository.scanRows(rows)
		if err != nil {
			return data, err
		}
		data = append(data, temp)
	}

	return data, nil
}

func (repository CategoryRepository) ReadBy(column, value, operator string) (data models.Category, err error) {
	statement := categorySelectStatement + ` from categories c ` + categoryJoinStatement + ` where ` + column + `` + operator + `$1 and c.deleted_at is null`
	row := repository.DB.QueryRow(statement, value)
	data, err = repository.scanRow(row)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (repository CategoryRepository) Edit(model models.Category) (res string, err error) {
	statement := `update categories set name=$1, slug=$2, parent_id=$3, is_active=$4, file_icon_id=$5, file_background_id=$6, updated_at=$7 where id=$8 returning id`
	err = repository.DB.QueryRow(statement, model.Name, model.Slug, model.ParentID.String, model.IsActive, model.FileIconID.String, model.FileBackgroundID.String,
		model.UpdatedAt, model.ID).Scan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (repository CategoryRepository) Add(model models.Category) (res string, err error) {
	statement := `insert into categories (name,slug,parent_id,is_active,file_icon_id,file_background_id,created_at,updated_at) values ($1,$2,$3,$4,$5,$6,$7,$8) returning id`
	err = repository.DB.QueryRow(statement, model.Name, model.Slug, model.ParentID.String, model.IsActive, model.FileIconID.String, model.FileBackgroundID.String,
		model.CreatedAt, model.UpdatedAt).Scan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (repository CategoryRepository) Delete(model models.Category) (res string, err error) {
	statement := `update categories set updated_at=$1, deleted_at=$2 where id=$3 returning id`
	err = repository.DB.QueryRow(statement, model.UpdatedAt, model.DeletedAt.Time, model.ID).Scan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}
