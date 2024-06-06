package mysql

import (
	"database/sql"
	"fmt"

	"Todo_Application/pkg/models"
)

type TodoModel struct { // Define a struct TodoModel which wraps a sql.DB connection pool.
	DB         *sql.DB
	InsertStmt *sql.Stmt
}

func (m *TodoModel) Insert(name, modified string) (int, error) { // This will insert new datas into the database.
	tx, err := m.DB.Begin() //calling begin to start db transaction
	if err != nil {
		return 0, err
	}

	stmt := `INSERT INTO Todo (name, created, modified) 
			VALUES(?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	InsertStmt, _ := tx.Prepare(stmt) //tx.Exec(stmt, name, modified) //Executing the insert into query from stmt and adding the name and modified values passed to fn
	result, err := InsertStmt.Exec(name, modified)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	ids, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	err = tx.Commit()
	defer InsertStmt.Close()
	return int(ids), err
}

func (m *TodoModel) Get(id int) (*models.Todo, error) { // This will return a specific name based on its id.
	stmt := `SELECT id, name, created, modified FROM Todo 
   			 WHERE modified > UTC_TIMESTAMP() AND id = ? `
	row := m.DB.QueryRow(stmt, id) //QueryRow take the query stmt and select a row of data
	s := &models.Todo{}
	err := row.Scan(&s.ID, &s.Name, &s.Created, &s.Modified) //row.Scan copy columns from matched row
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}

func (m *TodoModel) Delete(id int) (*models.Todo, error) { //This will delete the specified data
	stmt := `DELETE FROM Todo 
    		 WHERE id = ? `
	_, err := m.DB.Exec(stmt, id) //DB.Exec will perform deletion query
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (m *TodoModel) Update(name string, id int) (*models.Todo, error) { //Updates the specified name
	stmt := `UPDATE Todo
			SET name = ? 
			WHERE id =?`
	_, err := m.DB.Exec(stmt, name, id) //DB.Exec will perform updation query
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (m *TodoModel) GetAll() ([]*models.Todo, error) { // This will return all the created todo list
	stmt := `SELECT id, name, created, modified FROM Todo 
    		 WHERE modified > UTC_TIMESTAMP() `
	rows, err := m.DB.Query(stmt) //DB.Query selects the required fields.almost similar to DB.QueryRow
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	task_array := []*models.Todo{} //calling the todo model to task_array
	for rows.Next() {              //looping each rows
		s := &models.Todo{} //s to access the model elements
		err = rows.Scan(&s.ID, &s.Name, &s.Created, &s.Modified)
		if err != nil {
			return nil, err
		}
		task_array = append(task_array, s) // Populate slice of todos in []*models.Todo
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return task_array, nil // return the slice
}

func (m *TodoModel) ErrorManage(errors map[string]string) (*models.Todo, error) {
	//task_array := []*models.Todo{}
	s := &models.Todo{Errors: make(map[string]string)}

	if len(errors) > 0 {
		for x, _ := range errors {
			s.Errors[x] = errors[x]
			fmt.Println(s.Errors[x])
		}

	}
	// s.Errors = errors
	// task_array = append(task_array, s)
	// log.Println(task_array)
	return s, nil
}
