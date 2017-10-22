package tasks

// NOTE we do not log anything in this file! Leave it to the http handler to handle errors.

import (
	"database/sql"
	"time"
)

// TODO add done, done_by, c_time (maybe m_time?)
type Task struct {
	Key        int64     `json:"key"`
	ShortDesc  string    `json:"short_desc"`
	LongDesc   string    `json:"long_desc"`
	Done       bool      `json:"done"`
	DoneBy     string    `json:"done_by"`
	CreateTime time.Time `json:"c_time"`
	ChoreId    int64     `json:"chore_id"`
}

func setChoresNowStmt(db *sql.DB) (*sql.Stmt, error) {
	return db.Prepare("INSERT INTO tasks(chore_id, c_time) SELECT chore_id, NULL from chores where (morning = ? OR night = ?) AND ((dwm = 'd') OR (dwm = 'w' AND day = ?) OR (dwm = 'm' AND date = ?))")
}

func (h *handler) setChoresNow(t time.Time) error {
	var (
		isMorning int
		isNight   int
	)
	if t.Hour() >= 5 || t.Hour() < 17 {
		isMorning = 1
		isNight = -1
	} else {
		isNight = -1
		isMorning = 1
	}
	_, err := h.setChoresNowStmt.Exec(isMorning, isNight, t.Weekday(), t.Day())
	return err
}

func getAllTasksStmt(db *sql.DB) (*sql.Stmt, error) {
	return db.Prepare("SELECT task_id, chores.short_desc, chores.long_desc, done, done_by, tasks.c_time, tasks.chore_id from tasks left join chores on (tasks.chore_id = chores.chore_id)")
}

func (h *handler) getAllTasks() ([]*Task, error) {
	rows, err := h.getAllTasksStmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var (
		taskList   []*Task
		task_id    int64
		short_desc string
		long_desc  string
		done       bool
		done_by    string
		chore_id   int64
		c_time     time.Time
	)
	for rows.Next() {
		err := rows.Scan(&task_id, &short_desc, &long_desc, &done, &done_by, &c_time, &chore_id)
		if err != nil {
			return nil, err
		}
		t := &Task{task_id, short_desc, long_desc, done, done_by, c_time, chore_id}
		taskList = append(taskList, t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return taskList, nil
}

func modifyTaskStmt(db *sql.DB) (*sql.Stmt, error) {
	return db.Prepare("UPDATE tasks set done = ?, done_by = ? where task_id = ?")
}

func (h *handler) modifyTask(task *Task) error {
	_, err := h.modifyTaskStmt.Exec(
		task.Done,
		task.DoneBy,
	)
	return err
}