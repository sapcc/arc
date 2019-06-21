package models

import (
	"database/sql"
	"encoding/json"
	"reflect"

	"github.com/sapcc/arc/api-server/auth"
)

type Db interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type JSONB map[string]interface{}

func (j JSONB) Value() (interface{}, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *JSONB) Scan(value interface{}) error {
	// nothing is set yet in the clumn
	if value == nil {
		return nil
	}

	// convert to json
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}

	return nil
}

func JSONBfromString(data string) (*JSONB, error) {
	jsonbObj := make(JSONB)
	err := json.Unmarshal([]byte(data), &jsonbObj)
	if err != nil {
		return nil, err
	}
	return &jsonbObj, nil
}

func JobUserToJSONB(user auth.User) (*JSONB, error) {
	userJson, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	userJSONB, err := JSONBfromString(string(userJson))
	if err != nil {
		return nil, err
	}
	return userJSONB, nil
}

func CompareUserWithJobUser(user auth.User, jobUser JSONB) (bool, error) {
	checkUser, err := JobUserToJSONB(user)
	if err != nil {
		return false, err
	}

	eq := reflect.DeepEqual(jobUser, *checkUser)
	if err != nil {
		return false, err
	}
	return eq, nil
}
