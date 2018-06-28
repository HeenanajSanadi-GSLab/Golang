package daos

import (
	"database/sql"

	"../models"
	"log"
	
)

// UserImplMysql implementation of user from Mysql
type UserImplMysql struct {
}

// Create a user in mysql
func (dao UserImplMysql) Create(u *models.User) error {
	query1 := "INSERT INTO department (department_id, department_name) VALUES(?, ?)"
	db := get()
	//defer db.Close()
	stmt, err := db.Prepare(query1)
	result, _ := stmt.Exec(u.Department.Department_id, u.Department.Department_name)

	query2 := "INSERT INTO allusers (name, gender, age, department_id) VALUES(?, ?, ?, ?)"
	db = get()
	defer db.Close()
	stmt, err = db.Prepare(query2)

	if err != nil {
		return err
	}

	//defer stmt.Close()
	result, err2 := stmt.Exec(u.Name, u.Gender, u.Age, u.Department.Department_id)
	if err2 != nil {
		return err
	}

	id, err3 := result.LastInsertId()
	if err3 != nil {
		return err
	}

	u.Id = int(id)
	return nil
}

// GetAll users in mysql
func (dao UserImplMysql) GetAll() ([]models.User, error) {
	query := "SELECT allusers.id, allusers.name, allusers.gender, allusers.age, department.department_id, department.department_name FROM allusers, department where allusers.department_id = department.department_id"
	users := make([]models.User, 0)
	db := get()
	defer db.Close()

	stmt, err := db.Prepare(query)
	if err != nil {
		return users, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return users, err
	}
	for rows.Next() {
		var row models.User
		var rowdept models.Department
		err := rows.Scan(&row.Id, &row.Name, &row.Gender, &row.Age, &rowdept.Department_id, &rowdept.Department_name)
		if err != nil {
			return nil, err
		}
		row.Department = &rowdept
		users = append(users, row)
	}
	return users, nil
}

// Delete user in mysql
func (dao UserImplMysql) Delete(id int) error {
	query := "DELETE FROM allusers WHERE id = ?"
	db := get()
	defer db.Close()
	stmt, err := db.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()
	_ , err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

// Get user in mysql
func (dao UserImplMysql) Get(id int) (models.User, error) {

	query := "SELECT allusers.id, allusers.name, allusers.gender, allusers.age FROM allusers WHERE allusers.id = ?"
	//query := "SELECT allusers.id, allusers.name, allusers.gender, allusers.age, department.department_id from allusers JOIN department ON department.department_id = allusers.department_id WHERE allusers.id=?"
	db := get()
	defer db.Close()
	stmt, err := db.Prepare(query)
	if err != nil {
		return models.User{}, err
	}
	//defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		return models.User{}, err
	}
	var row models.User
	for rows.Next() {
		err := rows.Scan(&row.Id, &row.Name, &row.Gender, &row.Age)
		if err != nil {
			return models.User{}, err
		}
		}
		
	query2 := "SELECT department.department_id, department.department_name from department JOIN allusers ON allusers.department_id = department.department_id WHERE allusers.id=?"
	db = get()
	//defer db.Close()
	stmt1, _ := db.Prepare(query2)
	
	defer stmt1.Close()
	var row1 models.Department
	rows1, err := stmt1.Query(id)
	for rows1.Next() {
		err := rows1.Scan(&row1.Department_id, &row1.Department_name)
		if err != nil {
			return models.User{}, err
		}
	}
	row.Department = &row1

	return row, err
}

// Update user in mysql
func (dao UserImplMysql) Update(u *models.User) error {
	//query := "UPDATE allusers SET name = ?, gender = ?, age = ? WHERE id=?"
	query := "UPDATE allusers SET name = ?, gender = ?, age = ?, department_id = ? WHERE id=?"
	db := get()
	defer db.Close()
	stmt, err := db.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Query(NewNullString(u.Name), NewNullString(u.Gender), u.Age, u.Department.Department_id, u.Id)
	if err != nil {
		return err
	}
	return nil
}

// NewNullString check for null
func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
/*func (dao UserImplMysql) DeleteDepartmentByUserId(id int) error { 
	//get_dept_query := "SELECT department_id from allusers where id=?"
	q := fmt.Sprintf("SELECT * from allusers where id=%d ", id)
	db := get()
	dept_row, err1 := db.Query(q)
  if err1 != nil {
    log.Fatal(err1)
  }
  defer dept_row.Close()
  var d models.User
  for dept_row.Next() {
    dept_row.Scan(&d.Department.Department_id)
}
	log.Print(d.Department.Department_id)
	deptId:= d.Department.Department_id
	query := "select allusers.id from allusers, department where department.department_id = ? and allusers.department_id = department.department_id"
	db = get()
	defer db.Close()
	stmt, _ := db.Prepare(query)
	result, _ := stmt.Exec(deptId)
	//defer stmt.Close()
	//count, _ = stmt.Exec(id)
	log.Print("*************",result)
	if result == nil {
		log.Print(deptId)
		query := "DELETE FROM department WHERE id = ?"
		db := get()
		stmt, _ := db.Prepare(query)
		_ , e := stmt.Exec(deptId)
		if e != nil {
		return e
	}
	}
	return nil
	}*/


func (dao UserImplMysql) CheckDuplicateUser(name string) bool{
		
		query := "SELECT * from allusers where name=?"
		db := get()
		stmt, _ := db.Prepare(query)
		rows, _ := stmt.Query(name)
			
		log.Printf("inside CheckDuplicateUser() method >>>>")
		if rows.Next() {
			return true;	
		}
		return false

	}


