package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"

	"github.com/jinzhu/gorm"

	mocket "github.com/selvatico/go-mocket"
	v2 "gopkg.in/validator.v2"
)

func TestUserDtoValidation(t *testing.T) {
	upr := UsernamePassRequest{
		Username: "strahinjamalabosna",
		Password: "strale84",
	}

	t.Run("All is well", func(t *testing.T) {
		err := v2.Validate(upr)
		if err != nil {
			t.Errorf("Unexpected error: %s", err.Error())
		}
	})

	t.Run("Username and password empty", func(t *testing.T) {
		err := v2.Validate(UsernamePassRequest{})
		if err == nil {
			t.Error("Should have failed")
		} else if !strings.Contains(err.Error(), "Username: zero value, less than min") ||
			!strings.Contains(err.Error(), "Password: zero value, less than min") {
			t.Errorf("Unexpected error: %s", err.Error())
		}
	})

	t.Run("Only small letters are allowed for username", func(t *testing.T) {
		err := v2.Validate(UsernamePassRequest{Username: "strale84", Password: "strale84"})
		if err == nil {
			t.Error("Should have failed")
		} else if !strings.Contains(err.Error(), "Username: regular expression mismatch") {
			t.Errorf("Unexpected error: %s", err.Error())
		}
	})
}

func SetupTests() *gorm.DB {
	mocket.Catcher.Register()
	mocket.Catcher.Logging = true
	gdb, _ := gorm.Open(mocket.DriverName, "somestring")
	//db.SetDBForTesting(gdb)
	return gdb
}

func TestSignup(t *testing.T) {
	t.Run("Successful signup", func(t *testing.T) {
		SetupTests()
		searchMock := mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "users"  WHERE (username = strahinjamalabosna)`).WithError(gorm.ErrRecordNotFound)
		insertMock := mocket.Catcher.NewMock().WithQuery(`INSERT INTO "users"`).WithID(1).WithRowsNum(1)
		body := `{"username":"strahinjamalabosna","password":"strale84"}`
		request := httptest.NewRequest("POST", "http://example.org/user/signup", strings.NewReader(body))
		response := httptest.NewRecorder()
		SignupUser(response, request, httprouter.Params{})
		if !searchMock.Triggered {
			t.Error("Select statement should have ran!")
		} else if searchMock.Error != gorm.ErrRecordNotFound {
			t.Error("No record should be found!")
		}
		if !insertMock.Triggered {
			t.Error("Insert statement should have been executed!")
		}
		if response.Code != http.StatusNoContent {
			t.Error("Wrong status code returned!")
		}
	})
}

func TestCsvReading(t *testing.T) {
	f, err := os.Open("/home/s.dobrijevic/Downloads/cities.txt")
	defer f.Close()
	r := csv.NewReader(f)
	fields, err := r.Read()
	for err == nil && fields != nil {
		fmt.Println(fields, len(fields))
		fields, err = r.Read()
	}
}
