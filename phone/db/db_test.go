package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	. "gophercises/phone/normalizer"
	"regexp"
	"testing"
)

const TableName = "phone_numbers"

var PhoneNumbers = []PhoneNumber{
	{
		ID:         1,
		Number:     "(123) 456-7893",
		Owner:      "EE",
		NumberType: "Mobile",
	},
	{
		ID:         2,
		Number:     "1234567890",
		Owner:      "EE",
		NumberType: "Mobile",
	},
}

type DBTestSuite struct {
	suite.Suite
	unit   *PhoneDB
	gomega *GomegaWithT
}

func Test_DBTestSuite(t *testing.T) {
	var db *gorm.DB
	var err error
	if db, err = gorm.Open("postgres", "user=postgres dbname=postgres sslmode=disable"); err != nil {
		panic(fmt.Errorf("couldn't open postgres database: %s", err))
	}
	defer db.Close()
	db.DropTableIfExists(&PhoneNumber{})
	db = db.CreateTable(&PhoneNumber{})

	testSuite := DBTestSuite{unit: &PhoneDB{db: db}, gomega: NewGomegaWithT(t)}
	suite.Run(t, &testSuite)
}

func (s *DBTestSuite) Test_PhoneNumbersInsertedSuccessfully() {
	firstNumber := PhoneNumber{
		ID:         1,
		Number:     "1234567890",
		Owner:      "Vlad",
		NumberType: "Mobile",
	}
	s.unit.Insert(&firstNumber)
	secondNumber := PhoneNumber{
		ID:         2,
		Number:     "9876543210",
		Owner:      "Alan",
		NumberType: "Mobile",
	}
	s.unit.Insert(&secondNumber)
	actualNumber := s.unit.Query(1)
	s.gomega.Expect(actualNumber).Should(Equal(firstNumber))
}

func (s *DBTestSuite) Test_NumbersAreNormalized() {
	normalizationPattern := NormalizationPattern
	var matched bool
	var err error
	for _, number := range PhoneNumbers {
		number.Normalize()
		if matched, err = regexp.MatchString(normalizationPattern, number.Number); err != nil {
			panic(fmt.Errorf("error occurred while trying to match number %s with the pattern %s: %s\n",
				number.Number, normalizationPattern, err))
		}
		s.gomega.Expect(matched).Should(Equal(true))
	}
}
