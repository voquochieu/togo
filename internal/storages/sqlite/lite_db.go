package sqllite

import (
	"context"
	"database/sql"

	"github.com/manabie-com/togo/internal/storages"
)

type (
	// LiteDB for working with sqllite
	LiteDB struct {
		DB *sql.DB
	}

	UserLimitAndTasks struct {
		MaxTodo   int
		TaskCount int
	}
)

// RetrieveTasks returns tasks if match userID AND createDate.
func (l *LiteDB) RetrieveTasks(ctx context.Context, userID, createdDate string) ([]*storages.Task, error) {
	stmt := `SELECT id, content, user_id, created_date FROM tasks WHERE user_id = ? AND created_date = ?`
	rows, err := l.DB.QueryContext(ctx, stmt, l.convertToSQLNullString(userID), l.convertToSQLNullString(createdDate))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*storages.Task
	for rows.Next() {
		t := &storages.Task{}
		err := rows.Scan(&t.ID, &t.Content, &t.UserID, &t.CreatedDate)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// AddTask adds a new task to DB
func (l *LiteDB) AddTask(ctx context.Context, t *storages.Task) error {
	stmt := `INSERT INTO tasks (id, content, user_id, created_date) VALUES (?, ?, ?, ?)`
	_, err := l.DB.ExecContext(ctx, stmt, &t.ID, &t.Content, &t.UserID, &t.CreatedDate)
	if err != nil {
		return err
	}

	return nil
}

// ValidateUser returns tasks if match userID AND password
func (l *LiteDB) ValidateUser(ctx context.Context, userID, password string) bool {
	stmt := `SELECT id FROM users WHERE id = ? AND password = ?`
	row := l.DB.QueryRowContext(ctx, stmt, l.convertToSQLNullString(userID), l.convertToSQLNullString(password))
	u := &storages.User{}
	err := row.Scan(&u.ID)
	if err != nil {
		return false
	}

	return true
}

func (l *LiteDB) RetrieveUserMaxTodoAndTaskCount(ctx context.Context, userID, created_date string) (int, int, error) {
	stmt := `SELECT u.max_todo, COUNT(t.id) AS task_count 
			FROM users AS u 
			LEFT JOIN tasks AS t ON t.user_id = u.id AND t.created_date = ? 
			WHERE u.id = ?`
	row := l.DB.QueryRowContext(ctx, stmt, l.convertToSQLNullString(created_date), l.convertToSQLNullString(userID))
	info := &UserLimitAndTasks{}
	err := row.Scan(&info.MaxTodo, &info.TaskCount)
	if err != nil {
		return 0, 0, err
	}

	return info.MaxTodo, info.TaskCount, nil
}

func (l *LiteDB) convertToSQLNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  len(s) > 0,
	}
}
