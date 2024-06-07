package mysql

import (
	"Todo_Application/pkg/models"
	"database/sql"
)

type SpecialTask struct {
	DB *sql.DB
}

func (s *SpecialTask) Insert(name, types, modified string) (int, error) { // This will insert new datas into the database.
	stmt := `INSERT INTO Special (name, created, modified, Type) 
			VALUES(?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY),?)`
	result, _ := s.DB.Exec(stmt, name, modified, types) //Executing the insert into query from stmt and adding the name and modified values passed to fn
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), err
}

func (s *SpecialTask) Delete(name string) (*models.Special, error) { //This will delete the specified data
	stmt2 := `SELECT name FROM Special` //selecting name from the special structure
	rows, err := s.DB.Query(stmt2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() { //looping through rows and compare the name for checking if present in database Special
		st := &models.Special{}
		if st.Name == name {
			stmt := `DELETE FROM Special
						WHERE name = ? `
			_, err := s.DB.Exec(stmt, name) //DB.Exec will perform deletion query
			if err != nil {
				return nil, err
			}
		}
	}
	return nil, nil
}

func (m *SpecialTask) GetSpecial() ([]*models.Special, error) { // This will return all the created todo list
	stmt := `SELECT id, name, created, modified,Type FROM Special 
    		 WHERE modified > UTC_TIMESTAMP() `
	rows, err := m.DB.Query(stmt) //DB.Query selects the required fields.almost similar to DB.QueryRow
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	task_special := []*models.Special{} //calling the special model to task_special array
	for rows.Next() {                   //looping each rows
		s := &models.Special{} //s to access the model elements
		err = rows.Scan(&s.ID, &s.Name, &s.Created, &s.Modified, &s.Type)
		if err != nil {
			return nil, err
		}
		task_special = append(task_special, s) // Populate slice of special tasks in []*models.special
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return task_special, nil // return the slice
}
