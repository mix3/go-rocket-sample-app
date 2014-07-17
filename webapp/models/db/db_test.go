package db

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mix3/go-rocket-sample-app/webapp/util"
	"github.com/stretchr/testify/assert"
)

func MockNow(str string) func() time.Time {
	return func() time.Time {
		ret, _ := time.Parse("2006-01-02 15:04:05 -0700", str)
		return ret
	}
}

func setup(t *testing.T) *DB {
	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("Environment `DATABASE_URL` undefined")
	}
	return GetDB()
}

func teardown() {
	//...
}

func Test(t *testing.T) {
	d := setup(t)
	defer func() {
		d.DropTables()
		d.DbMap.Db.Close()
	}()
	var hash, newHash string
	var err error
	var email Email
	util.FixedTime("2014-07-01 00:00:00 +0900", func() {
		// InterimRegisterEmail
		hash, err = d.InterimRegisterEmail("hoge@example.com")
		if err != nil {
			t.Fatal(err)
		}
		err = d.SelectOne(&email, "SELECT * FROM email WHERE hash = $1", hash)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "hoge@example.com", email.Email)
		assert.Equal(t, false, email.Status)
		assert.Equal(t, hash, email.Hash)
		assert.Equal(t, "2014-07-01 00:00:00 +0900 JST", fmt.Sprintf("%v", email.CreatedAt))
		assert.Equal(t, "2014-07-01 00:00:00 +0900 JST", fmt.Sprintf("%v", email.UpdatedAt))
	})

	util.FixedTime("2014-07-01 01:00:00 +0900", func() {
		// Duplicate InterimRegisterEmail
		newHash, err = d.InterimRegisterEmail("hoge@example.com")
		if err != nil {
			t.Fatal(err)
		}
		err = d.SelectOne(&email, "SELECT * FROM email WHERE hash = $1", newHash)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "hoge@example.com", email.Email)
		assert.Equal(t, false, email.Status)
		assert.Equal(t, newHash, email.Hash)
		assert.Equal(t, "2014-07-01 00:00:00 +0900 JST", fmt.Sprintf("%v", email.CreatedAt))
		assert.Equal(t, "2014-07-01 01:00:00 +0900 JST", fmt.Sprintf("%v", email.UpdatedAt))
	})

	util.FixedTime("2014-07-01 02:00:00 +0900", func() {
		// RegisterEmail
		err = d.RegisterEmail(newHash)
		if err != nil {
			t.Fatal(err)
		}
		err = d.SelectOne(&email, "SELECT * FROM email WHERE hash = $1", newHash)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "hoge@example.com", email.Email)
		assert.Equal(t, true, email.Status)
		assert.Equal(t, newHash, email.Hash)
		assert.Equal(t, "2014-07-01 00:00:00 +0900 JST", fmt.Sprintf("%v", email.CreatedAt))
		assert.Equal(t, "2014-07-01 02:00:00 +0900 JST", fmt.Sprintf("%v", email.UpdatedAt))
	})

	util.FixedTime("2014-07-01 02:00:00 +0900", func() {
		// RegisterRemind
		remindAt, err := util.GetTime("2014-08-01 00:00:00 +0900")
		if err != nil {
			t.Fatal(err)
		}
		err = d.RegisterRemind("fuga@example.com", "message", remindAt)
		if err == nil {
			t.Fatal(fmt.Errorf("fuga@example.com unregistered"))
		}
		err = d.RegisterRemind("hoge@example.com", "message", remindAt)
		if err != nil {
			t.Fatal(err)
		}
		var reminds []Remind
		_, err = d.Select(&reminds, "SELECT * FROM remind")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 1, len(reminds))
		assert.Equal(t, "hoge@example.com", reminds[0].To)
		assert.Equal(t, "message", reminds[0].Message)
		assert.Equal(t, "2014-07-01 02:00:00 +0900 JST", fmt.Sprintf("%v", reminds[0].CreatedAt))
		assert.Equal(t, "2014-08-01 00:00:00 +0900 JST", fmt.Sprintf("%v", reminds[0].RemindAt))
	})

	util.FixedTime("2014-07-01 03:00:00 +0900", func() {
		// RegisterRemind
		remindAt, err := util.GetTime("2014-08-01 01:00:00 +0900")
		if err != nil {
			t.Fatal(err)
		}
		err = d.RegisterRemind("hoge@example.com", "happy", remindAt)
		if err != nil {
			t.Fatal(err)
		}
		var reminds []Remind
		_, err = d.Select(&reminds, "SELECT * FROM remind")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 2, len(reminds))

		assert.Equal(t, "hoge@example.com", reminds[1].To)
		assert.Equal(t, "happy", reminds[1].Message)
		assert.Equal(t, "2014-07-01 03:00:00 +0900 JST", fmt.Sprintf("%v", reminds[1].CreatedAt))
		assert.Equal(t, "2014-08-01 01:00:00 +0900 JST", fmt.Sprintf("%v", reminds[1].RemindAt))

		assert.Equal(t, "hoge@example.com", reminds[0].To)
		assert.Equal(t, "message", reminds[0].Message)
		assert.Equal(t, "2014-07-01 02:00:00 +0900 JST", fmt.Sprintf("%v", reminds[0].CreatedAt))
		assert.Equal(t, "2014-08-01 00:00:00 +0900 JST", fmt.Sprintf("%v", reminds[0].RemindAt))
	})

	util.FixedTime("2014-08-01 00:00:00 +0900", func() {
		// RemindList
		list := d.RemindList()
		assert.Equal(t, 0, len(list))
	})
	util.FixedTime("2014-08-01 00:00:01 +0900", func() {
		// RemindList
		list := d.RemindList()
		assert.Equal(t, 1, len(list))
	})
	util.FixedTime("2014-08-01 01:00:00 +0900", func() {
		// RemindList
		list := d.RemindList()
		assert.Equal(t, 1, len(list))
	})
	util.FixedTime("2014-08-01 01:00:01 +0900", func() {
		// RemindList
		list := d.RemindList()
		assert.Equal(t, 2, len(list))
	})

	util.FixedTime("2014-08-01 01:00:01 +0900", func() {
		// DeleteRemind
		list := d.RemindList()
		assert.Equal(t, 2, len(list))
		for _, remind := range list {
			err := d.DeleteRemind(remind)
			if err != nil {
				t.Fatal(err)
			}
		}
		list = d.RemindList()
		assert.Equal(t, 0, len(list))
	})
}
