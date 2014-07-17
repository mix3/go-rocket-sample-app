package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"github.com/mix3/go-rocket-sample-app/webapp/util"
)

type Email struct {
	Email     string    `db:"email"`
	Status    bool      `db:"status"`
	Hash      string    `db:"hash"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Remind struct {
	Id        int64     `db:"id"`
	To        string    `db:"to"`
	Message   string    `db:"message"`
	RemindAt  time.Time `db:"remind_at"`
	CreatedAt time.Time `db:"created_at"`
}

func (r *Remind) PreInsert(s gorp.SqlExecutor) error {
	r.CreatedAt = now()
	return nil
}

type DB struct {
	*gorp.DbMap
}

var db *DB

func init() {
	dbConn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	dbmap := &gorp.DbMap{Db: dbConn, Dialect: gorp.PostgresDialect{}}
	dbmap.TraceOn("[gorp]", log.New(os.Stdout, "myapp:", log.Lmicroseconds))
	dbmap.AddTableWithName(Email{}, "email").SetKeys(false, "email")
	dbmap.AddTableWithName(Remind{}, "remind").SetKeys(true, "id")
	dbmap.CreateTablesIfNotExists()
	db = &DB{dbmap}
}

func GetDB() *DB {
	return db
}

func now() time.Time {
	return util.Now()
}

func nowFormat() string {
	return now().Format("2006-01-02 15:04:05")
}

func (d *DB) InterimRegisterEmail(email string) (string, error) {
	//"UPDATE email SET hash = $1, updated_at = $2 WHERE email = $3", hash, now, email
	//"INSERT INTO email (email, hash, created_at, updated_at)
	// SELECT $1, $2, $3, $4
	// WHERE NOT EXISTS (SELECT * FROM email WHERE email = $5)", email, hash, now, now, email
	hash := util.GenUUID()
	var err error
	_, err = d.Exec("UPDATE email SET hash = $1, updated_at = $2 WHERE email = $3", hash, nowFormat(), email)
	if err != nil {
		return "", err
	}
	_, err = d.Exec(`
INSERT INTO email (email, status, hash, created_at, updated_at)
SELECT $1, FALSE, $2, $3, $4
WHERE NOT EXISTS (SELECT * FROM email WHERE email = $5)`, email, hash, nowFormat(), nowFormat(), email)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func (d *DB) RegisterEmail(hash string) error {
	//"UPDATE email SET status = TRUE WHERE email = $2 AND hash = $3 AND status = FALSE", email, hash
	_, err := d.Exec("UPDATE email SET status = TRUE, updated_at = $1 WHERE hash = $2 AND status = FALSE", nowFormat(), hash)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) RegisterRemind(to, message string, remindAt time.Time) error {
	var email Email
	err := d.SelectOne(&email, "SELECT * FROM email WHERE email = $1 AND status = true", to)
	if err != nil {
		return err
	}
	remind := Remind{
		To:       to,
		Message:  message,
		RemindAt: remindAt,
	}
	return d.Insert(&remind)
}

func (d *DB) RemindList() []Remind {
	// SELECT * FROM remind WHERE AND remind_at < NOW()
	var reminds []Remind
	d.Select(&reminds, "SELECT * FROM remind WHERE remind_at < $1 ORDER BY id DESC", nowFormat())
	return reminds
}

func (d *DB) DeleteRemind(remind Remind) error {
	// DELETE FROM remind WHERE Id = $1
	_, err := d.Delete(&remind)
	if err != nil {
		return err
	}
	return nil
}
