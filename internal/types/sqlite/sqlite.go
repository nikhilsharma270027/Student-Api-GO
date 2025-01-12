package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // if its getting used indirectly use "_" before
	"github.com/nikhilsharma270027/GOLang-student-api/internal/config"
	"github.com/nikhilsharma270027/GOLang-student-api/internal/types"
)

type Sqlite struct {
	Db *sql.DB
}

// use * taking its reference, copy everything here
func New(cfg *config.Config) (*Sqlite, error) {
	// we shd open database here
	// in open we need to specific sql driver, path og the storagr fleconst
	// the storage path is in config cfg.storagePath
	db, err := sql.Open("sqlite3", cfg.StoragePath) // it will return error and DB instance "db"
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil

}

// to make instance of sql
// go dont have constructors
// we can make a func

// install database driver from go sqlite driver

// implement storage in student struct
func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	// prepare insertion of data and ? means avoiding sql injection for security
	stmt, err := s.Db.Prepare("INSERT INTO students (name , email, age) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, nil
	}
	// we last entried of data by lastInsertedId
	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return lastId, nil
	// return 0, nil

}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	// stmt, err := s.Db.Prepare("SELECT * FROM students WHERE id = ? LIMIT 1") // search based on id
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students WHERE id = ? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}

	defer stmt.Close()

	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id: %s", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("query error: %d", err)
	}

	return student, nil
}

func (s *Sqlite) GetStudents() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var students []types.Student
	// Every call to [Rows.Scan], even the first one, must be preceded by a call to [Rows.Next].
	for rows.Next() {
		var student types.Student

		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	return students, err
}
func (s *Sqlite) UpdateStudentById(name string, email string, age int, id int64) error {
	// stmt, err := s.Db.Prepare("SELECT * FROM students WHERE id = ? LIMIT 1") // search based on id
	stmt, err := s.Db.Prepare("UPDATE students SET name = ?, email = ?, age = ? WHERE id = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	// Execute the prepared statement with the provided data
	result, err := stmt.Exec(name, email, age, id)
	if err != nil {
		return err // Return the error if execution fails
	}

	// Check the number of rows affected to ensure the update was successful
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err // Return the error if unable to fetch affected rows
	}

	// If no rows were affected, return an error indicating the student was not found
	if rowsAffected == 0 {
		return fmt.Errorf("no student found with the given ID: %d", id)
	}

	return nil // Return nil if the update was successful
}

func (s *Sqlite) DeleteStudentById(id int64) error {
	stmt, err := s.Db.Prepare("DELETE FROM students WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no student found with the given ID: %d", id)
	}

	return nil
}
