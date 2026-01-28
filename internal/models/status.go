package models

import (
	"encoding/json"
	"fmt"
)

type Status string

const (
	StatusPending Status = "status_pending"
	StatusSuccess Status = "status_success"
	StatusFailed  Status = "status_failed"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusPending, StatusFailed, StatusSuccess:
		return true
	default:
		return false
	}
}

func (s *Status) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data,&str); err != nil{
		return  err
	}
	check := Status(str)
	if !check.IsValid() {
		return fmt.Errorf("Выберите одно из:Pending,Success,Failed")
	}
	 *s = check
	 return nil  
}
