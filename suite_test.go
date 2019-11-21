package test

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// gorm model
type Order struct {
	ID        int `gorm:"primary_key"`
	CreatedAt time.Time
	OrderNo   string `gorm:"type:varchar(32)"`
}

type demoTest struct {
	Suite
}

func (d *demoTest) TestA() {
	log.Println("Testing A")
	// test nothing
	acter := &Order{}
	d.Db.Where("order_no = ?", "not exists").First(acter)
	d.Assert().Equal(0, acter.ID)
}

func (d *demoTest) TestB() {
	log.Println("Testing B")
	// test exists
	acter := &Order{}
	d.Db.Where("order_no = ?", "abc").First(acter)
	d.Assert().Equal(1, acter.ID)
}

func TestSuite(t *testing.T) {
	demoSuite := new(demoTest)
	dns := "root:123456@tcp(127.0.0.1:33306)/test?charset=utf8&parseTime=true"
	models := []interface{}{&Order{}}
	seedFile := "seed.json"
	demoSuite.Init(dns, models, seedFile, false)
	suite.Run(t, demoSuite)
}
