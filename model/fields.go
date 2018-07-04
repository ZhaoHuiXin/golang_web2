package model

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/go-sql-driver/mysql"
)

type NullString struct {
	sql.NullString
}

func ToNullString(s string) NullString {
	return NullString{
		sql.NullString{
			String: s,
			Valid:  true,
		},
	}
}

func (v *NullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	} else {
		return json.Marshal(nil)
	}
}

func (v NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		v.Valid = true
		v.String = *s
	} else {
		v.Valid = false
	}
	return nil
}

type NullInt64 struct {
	sql.NullInt64
}

func (v *NullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	} else {
		return json.Marshal(nil)
	}
}

func (v NullInt64) UnmarshalJSON(data []byte) error {
	var s *int64
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		v.Valid = true
		v.Int64 = *s
	} else {
		v.Valid = false
	}
	return nil
}

type NullFloat64 struct {
	sql.NullFloat64
}

func (v *NullFloat64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Float64)
	} else {
		return json.Marshal(nil)
	}
}

func (v NullFloat64) UnmarshalJSON(data []byte) error {
	var s *float64
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		v.Valid = true
		v.Float64 = *s
	} else {
		v.Valid = false
	}
	return nil
}

type NullBool struct {
	sql.NullBool
}

func (v *NullBool) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Bool)
	} else {
		return json.Marshal(nil)
	}
}

func (v NullBool) UnmarshalJSON(data []byte) error {
	var s *bool
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		v.Valid = true
		v.Bool = *s
	} else {
		v.Valid = false
	}
	return nil
}

type NullTime struct {
	mysql.NullTime
}

func (v *NullTime) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Time)
	} else {
		return json.Marshal(nil)
	}
}

func (v NullTime) UnmarshalJSON(data []byte) error {
	var s *time.Time
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		v.Valid = true
		v.Time = *s
	} else {
		v.Valid = false
	}
	return nil
}
