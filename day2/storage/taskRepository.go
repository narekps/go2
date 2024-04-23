package storage

import (
	"fmt"
	"github.com/narekps/go2/day2/internal/app/models"
	"strings"
)

type TaskRepository struct {
	storage *Storage
}

func (tr *TaskRepository) SelectAllTasks() ([]*models.Task, error) {
	stmt, err := tr.storage.db.Prepare(fmt.Sprintf("SELECT id, task, tags, due FROM task"))
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tasks := make([]*models.Task, 0)
	for rows.Next() {
		task := &models.Task{}
		var tags string
		err := rows.Scan(&task.ID, &task.Text, &tags, &task.Due)
		if err != nil {
			return nil, err
		}
		task.Tags = strings.Split(tags, ",")
		tasks = append(tasks, task)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (tr *TaskRepository) FindTasksByTag(tag string) ([]*models.Task, error) {
	tasks, err := tr.SelectAllTasks()
	if err != nil {
		return nil, err
	}

	filteredTasks := make([]*models.Task, 0, len(tasks))
	for _, task := range tasks {
		for _, t := range task.Tags {
			if t == tag {
				filteredTasks = append(filteredTasks, task)
				break
			}
		}
	}

	return filteredTasks, nil
}

func (tr *TaskRepository) SelectTasksByDate(date string) ([]*models.Task, error) {
	stmt, err := tr.storage.db.Prepare(fmt.Sprintf("SELECT id, task, tags, due FROM task WHERE due BETWEEN ? AND ?"))
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(date, date+"T23:59:59")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tasks := make([]*models.Task, 0)
	for rows.Next() {
		task := &models.Task{}
		var tags string
		err := rows.Scan(&task.ID, &task.Text, &tags, &task.Due)
		if err != nil {
			return nil, err
		}
		task.Tags = strings.Split(tags, ",")
		tasks = append(tasks, task)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (tr *TaskRepository) Create(t *models.Task) (*models.Task, error) {
	tx, err := tr.storage.db.Begin()
	if err != nil {
		return nil, err
	}

	// Prepare prepared statement that can be reused.
	stmt, err := tx.Prepare("INSERT INTO task(task, tags, due) VALUES(?, ?, ?)")
	if err != nil {
		return nil, err
	}

	// close statement before exiting program.
	defer stmt.Close()

	// Execute statements for each tasks
	res, err := stmt.Exec(t.Text, strings.Join(t.Tags, ","), t.Due)
	if err != nil {
		return nil, err
	}

	// Commit the transaction, so that inserts are permanent.
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	t.ID = id

	return t, nil
}

func (tr *TaskRepository) FindById(id int) (*models.Task, bool, error) {
	stmt, err := tr.storage.db.Prepare(fmt.Sprintf("SELECT id, task, tags, due FROM task WHERE id = ?"))
	if err != nil {
		return nil, false, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	tasks := make([]*models.Task, 0)

	for rows.Next() {
		task := &models.Task{}
		var tags string
		err := rows.Scan(&task.ID, &task.Text, &tags, &task.Due)
		if err != nil {
			return nil, false, err
		}
		task.Tags = strings.Split(tags, ",")
		tasks = append(tasks, task)
	}

	err = rows.Err()
	if err != nil {
		return nil, false, err
	}

	if len(tasks) == 0 {
		return nil, false, nil
	}

	return tasks[0], true, nil
}

func (tr *TaskRepository) UpdateTask(t *models.Task) (*models.Task, error) {
	tx, err := tr.storage.db.Begin()
	if err != nil {
		return nil, err
	}

	// Prepare prepared statement that can be reused.
	stmt, err := tx.Prepare("UPDATE task SET task = ?, tags = ?, due = ? WHERE id = ?")
	if err != nil {
		return nil, err
	}

	// close statement before exiting program.
	defer stmt.Close()

	// Execute statements for each tasks
	_, err = stmt.Exec(t.Text, strings.Join(t.Tags, ","), t.Due, t.ID)
	if err != nil {
		return nil, err
	}

	// Commit the transaction, so that inserts are permanent.
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return t, nil
}

func (tr *TaskRepository) DeleteTask(id int) (*models.Task, error) {
	task, found, err := tr.FindById(id)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, nil
	}

	tx, err := tr.storage.db.Begin()
	if err != nil {
		return nil, err
	}

	// Prepare prepared statement that can be reused.
	stmt, err := tx.Prepare("DELETE FROM task WHERE id = ?")
	if err != nil {
		return nil, err
	}

	// close statement before exiting program.
	defer stmt.Close()

	// Execute statements for each tasks
	_, err = stmt.Exec(task.ID)
	if err != nil {
		return nil, err
	}

	// Commit the transaction, so that inserts are permanent.
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return task, nil
}

func (tr *TaskRepository) DeleteAllTasks() error {
	tx, err := tr.storage.db.Begin()
	if err != nil {
		return err
	}

	// Prepare prepared statement that can be reused.
	stmt, err := tx.Prepare("DELETE FROM task")
	if err != nil {
		return err
	}

	// close statement before exiting program.
	defer stmt.Close()

	// Execute statements for each tasks
	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	// Commit the transaction, so that inserts are permanent.
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
