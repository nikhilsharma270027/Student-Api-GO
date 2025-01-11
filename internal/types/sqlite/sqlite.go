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
