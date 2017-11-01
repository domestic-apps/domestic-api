package chores

// NOTE we do not log anything in this file! Leave it to the http handler to handle errors.

import (
	"database/sql"
)

// Chore holds information about a recurring household task
type Chore struct {
	Key       int64  `json:"key"`
	ShortDesc string `json:"short_desc"`
	LongDesc  string `json:"long_desc"`
	Morning   bool   `json:"morning"`
	Night     bool   `json:"night"`
	Period    string `json:"dwm"`
	Day       int    `json:"day"`
	Date      int    `json:"date"`
}

func addChoreStmt(db *sql.DB) (*sql.Stmt, error) {
	return db.Prepare("INSERT INTO chores(short_desc, long_desc, morning, night, dwm, day, date, c_time) VALUES(?,?,?,?,?,?,?,?)")
}

func (h *handler) addChore(chore *Chore) error {
	res, err := h.addChoreStmt.Exec(
		chore.ShortDesc,
		chore.LongDesc,
		chore.Morning,
		chore.Night,
		chore.Period,
		chore.Day,
		chore.Date,
		nil,
	)
	if err != nil {
		return err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	chore.Key = lastID
	return nil
}

func getAllChoresStmt(db *sql.DB) (*sql.Stmt, error) {
	return db.Prepare("SELECT chore_id, short_desc, long_desc, morning, night, dwm, day, date from chores")
}

func (h *handler) getAllChores() ([]*Chore, error) {
	rows, err := h.getAllChoresStmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var (
		choreID   int64
		shortDesc string
		longDesc  string
		morning   bool
		night     bool
		dwm       string
		day       int
		date      int
	)
	choreList := make([]*Chore, 0)
	for rows.Next() {
		err := rows.Scan(&choreID, &shortDesc, &longDesc, &morning, &night, &dwm, &day, &date)
		if err != nil {
			return nil, err
		}
		c := &Chore{choreID, shortDesc, longDesc, morning, night, dwm, day, date}
		choreList = append(choreList, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return choreList, nil
}

func getRecentChoresStmt(db *sql.DB) (*sql.Stmt, error) {
	// TODO Make it actually correct tho!
	return db.Prepare("SELECT * from chores")
}

/*
func getRecentChores(db *sql.DB, since *time.Time, excludedIds ...int) []Chore, error {

}
*/
func modifyChoreStmt(db *sql.DB) (*sql.Stmt, error) {
	return db.Prepare("UPDATE chores set short_desc = ?, long_desc = ?, morning = ?, night = ?, dwm = ?, day = ?, date = ? where chore_id = ?")
}

func (h *handler) modifyChore(chore *Chore) error {
	_, err := h.modifyChoreStmt.Exec(
		chore.ShortDesc,
		chore.LongDesc,
		chore.Morning,
		chore.Night,
		chore.Period,
		chore.Day,
		chore.Date,
		chore.Key,
	)
	return err
}

func deleteChoreStmt(db *sql.DB) (*sql.Stmt, error) {
	// TODO Also delete linked chores in tasks table (in a separate statement tho)
	return db.Prepare("DELETE FROM chores where chore_id = ?")
}

func (h *handler) deleteChore(choreID int64) error {
	_, err := h.deleteChoreStmt.Exec(choreID)
	return err
}
