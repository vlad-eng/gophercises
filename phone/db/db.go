package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	. "gophercises/phone/normalizer"
)

type PhoneDB struct {
	db *gorm.DB
}

const (
	user    = "postgres"
	dbname  = "postgres"
	sslmode = "disable"
)

func (p *PhoneDB) Create(values interface{}) {
	pSqlInfo := fmt.Sprintf("user=%s dbname=%s sslmode=%s", user, dbname, sslmode)
	var err error

	p.db, err = gorm.Open("postgres", pSqlInfo)
	must(err)
	p.ResetTable(values)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func (p *PhoneDB) ResetTable(values interface{}) {
	p.db.DropTable(values)
	p.db.CreateTable(values)
}

func (p *PhoneDB) Query(id int) PhoneNumber {
	number := PhoneNumber{}
	p.db.First(&number, id)
	return number
}

func (p *PhoneDB) QueryAll() []PhoneNumber {
	entries := make([]PhoneNumber, 0)
	p.db.Model(&PhoneNumber{}).Scan(&entries)
	return entries
}

func (p *PhoneDB) Insert(number *PhoneNumber) {
	p.db = p.db.Create(number)
}

func (p *PhoneDB) Update(number *PhoneNumber) {
	p.db.Model(&PhoneNumber{}).Where("id = ?", number.ID).Update("number", number.Number)
}

func (p *PhoneDB) Delete(number *PhoneNumber) {
	p.db.Model(&PhoneNumber{}).Delete(number)
}

func (p *PhoneDB) Drop(values ...interface{}) {
	p.db.DropTableIfExists(values)
}

func (p *PhoneDB) Close() {
	p.db.Close()
}
