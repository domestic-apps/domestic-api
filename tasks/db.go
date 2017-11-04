package tasks

// NOTE we do not log anything in this file! Leave it to the http handler to handle errors.

import (
	"database/sql"
	"time"
)

// Task holds all information related to a task.
// TODO add done, done_by, c_time (maybe m_time?)
type Task struct {
	Key        int64     `json:"key"`
	ShortDesc  string    `json:"short_desc"`
	Done       bool      `json:"done"`
	DoneBy     string    `json:"done_by"`
	CreateTime time.Time `json:"c_time"`
	ChoreID    int64     `json:"chore_id"`
}

func setChoresNowStmt(db *sql.DB) (*sql.Stmt, error) {
	return db.Prepare("INSERT INTO tasks(chore_id, c_time) SELECT chore_id, NULL from chores where (morning = ? OR night = ?) AND ((dwm = 'd') OR (dwm = 'w' AND day = ?) OR (dwm = 'm' AND date = ?))")
}

func (h *Handler) setChoresNow(t time.Time) error {
	var (
		isMorning int
		isNight   int
	)
	if t.Hour() >= 5 && t.Hour() < 17 {
		isMorning = 1
		isNight = -1
	} else {
		isMorning = -1
		isNight = 1
	}
	_, err := h.setChoresNowStmt.Exec(isMorning, isNight, t.Weekday(), t.Day())
	return err
}

func getAllTasksStmt(db *sql.DB) (*sql.Stmt, error) {
	return db.Prepare("SELECT task_id, chores.short_desc, done, done_by, tasks.c_time, tasks.chore_id from tasks left join chores on (tasks.chore_id = chores.chore_id)")
}

func (h *Handler) getAllTasks() ([]*Task, error) {
	rows, err := h.getAllTasksStmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var (
		taskID    int64
		shortDesc string
		done      bool
		doneByRaw sql.NullString
		doneBy    string
		choreID   int64
		cTime     time.Time
	)
	taskList := make([]*Task, 0)
	for rows.Next() {
		err := rows.Scan(&taskID, &shortDesc, &done, &doneByRaw, &cTime, &choreID)
		if err != nil {
			return nil, err
		}

		if doneByRaw.Valid {
			doneBy = doneByRaw.String
		} else {
			doneBy = ""
		}

		t := &Task{taskID, shortDesc, done, doneBy, cTime, choreID}
		taskList = append(taskList, t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return taskList, nil
}

// TODO prevent cheekiness: could have done_by contention
func modifyTaskStmt(db *sql.DB) (*sql.Stmt, error) {
	return db.Prepare("UPDATE tasks set done = ?, done_by = ? where task_id = ?")
}

func (h *Handler) modifyTask(task *Task) error {
	_, err := h.modifyTaskStmt.Exec(
		task.Done,
		task.DoneBy,
		task.Key,
	)
	return err
}
