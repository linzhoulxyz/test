package test

import (
	"fmt"
	"log"
	"os"

	"github.com/DATA-DOG/go-txdb"
	"github.com/jinzhu/gorm"
	"github.com/khaiql/dbcleaner"
	"github.com/khaiql/dbcleaner/engine"
	"github.com/romanyx/polluter"
	"github.com/stretchr/testify/suite"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Suite struct {
	suite.Suite

	Dns          string
	Db           *gorm.DB
	Models       []interface{} // 初始化GORM模型 *指针*
	Cleaner      dbcleaner.DbCleaner
	SeedFileName string
	Polluter     *polluter.Polluter
	Verbose      bool // 是否输出SQL语句

	modelTableNames []string // 模型表名
	txDbDriverName  string
}

// init suite info
func (s *Suite) Init(dns string, models []interface{}, seedFileName string, verbose bool) {
	if dns == "" {
		log.Println("without mysql dns, skip init")
		return
	}

	var err error

	// init db
	s.Dns = dns
	s.Db, err = gorm.Open("mysql", s.Dns)
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}
	if s.Verbose {
		s.Db.LogMode(true)
	}

	// init txdb
	s.txDbDriverName = "mysqltx"
	txdb.Register(s.txDbDriverName, "mysql", s.Dns)

	// init dbcleaner
	s.Cleaner = dbcleaner.New(dbcleaner.SetLockFileDir(os.TempDir()))
	cleanDb := engine.NewMySQLEngine(s.Dns)
	s.Cleaner.SetEngine(cleanDb)

	// init models
	s.Models = models
	s.modelTableNames = getModelsTablesName(s.Db, s.Models)
	log.Println("tables:", s.modelTableNames)

	// init seed
	s.SeedFileName = seedFileName
	s.Polluter = polluter.New(polluter.MySQLEngine(s.Db.DB()), polluter.JSONParser)
}

func (s *Suite) SetupSuite() {
	log.Println("SetupSuite")

	// init models
	if len(s.Models) > 0 {
		s.Db.AutoMigrate(s.Models...)
	}
}

func (s *Suite) TearDownSuite() {
	log.Println("TearDownSuite")
}

func (s *Suite) SetupTest() {
	log.Println("==========SetupTest==========")

	// lock tables and import seed
	if len(s.Models) > 0 {
		s.Cleaner.Acquire(s.modelTableNames...)
	}

	// import seed data
	if s.SeedFileName != "" {
		if f, err := os.Open(s.SeedFileName); err != nil {
			s.T().Fatalf("failed to open seed file: %s, %s", s.SeedFileName, err)
		} else if err := s.Polluter.Pollute(f); err != nil {
			s.T().Fatalf("failed to pollute: %s", err)
		}
	}
}

func (s *Suite) TearDownTest() {
	log.Println("=========TearDownTest=========")

	// clean tables data and release table
	if len(s.Models) > 0 {
		s.Cleaner.Clean(s.modelTableNames...)
	}
}
