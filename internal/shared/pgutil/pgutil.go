package pgutil

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// --------------------------------------------------------
// ARRAYS  - TEXT[], INT[], UUID[], etc,
// --------------------------------------------------------
func ToPostgresArray[T any](elements []T) pgtype.Array[T] {
	return pgtype.Array[T]{
		Elements: elements,
		Valid:    len(elements) > 0,
	}
}

func FromPostgresArray[T any](arr pgtype.Array[T]) []T {
	if !arr.Valid || len(arr.Elements) == 0 {
		return []T{}
	}
	return arr.Elements
}

//  ------------------------------------------------
//  JSONB - any struct <-> jsonb column
//  ------------------------------------------------

func ToJSONB(v any) ([]byte, error) {
	if v == nil {
		return nil, nil
	}
	return json.Marshal(v)
}

func FromJSONB[T any](data []byte) (*T, error) {
	if len(data) == 0 {
		return nil, nil
	}
	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

//  ------------------------------------------------
//  NULLABLE TEXT - *string <-> pgtype.Text
//  ------------------------------------------------

func ToNullableText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func FromNullableText(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

//  -----------------------------------------------
//  NULLABLE TIMESTAMPTZ - *time.Time <-> pgtype.Timestamptz
//  -----------------------------------------------

func ToNullableTimestampt(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

func FromNullableTimestampt(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

func ToTimestampt(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}
func FromTimestampt(t pgtype.Timestamptz) time.Time {
	return t.Time
}

//  ---------------------------------------------------
//  NULLABLE UUID - *uuid.UUID <-> pgtype.UUID
//  ---------------------------------------------------

func ToNullableUUID(u *uuid.UUID) pgtype.UUID {
	if u == nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Valid: true, Bytes: *u}
}

func FromNullableUUID(u pgtype.UUID) *uuid.UUID {
	if !u.Valid {
		return nil
	}
	id := uuid.UUID(u.Bytes)
	return &id
}

//  -------------------------------------------------
//  UUID - uuid.UUID <-> pgtype.UUID
//  -------------------------------------------------

func ToUUID(u uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: u, Valid: true}
}

func FromUUID(u pgtype.UUID) uuid.UUID {
	return u.Bytes
}

//  -----------------------------------------------
//  NULLABLE INT - *int64 <-> pgtype.Int8
//  -----------------------------------------------

func ToNullableInt(i *int64) pgtype.Int8 {
	if i == nil {
		return pgtype.Int8{Valid: false}
	}
	return pgtype.Int8{Int64: *i, Valid: true}
}

func FromNullableInt(i pgtype.Int8) *int64 {
	if !i.Valid {
		return nil
	}
	return &i.Int64
}

//  ------------------------------------------------
//  NULLABLE BOOL - *bool <-> pgtype.Bool
//  ------------------------------------------------

func ToNullableBool(b *bool) pgtype.Bool {
	if b == nil {
		return pgtype.Bool{Valid: false}
	}
	return pgtype.Bool{Bool: *b, Valid: true}
}

func FromNullableBool(b pgtype.Bool) *bool {
	if !b.Valid {
		return nil
	}
	return &b.Bool
}

// ---------------------------------------------------------------
// NULLABLE FLOAT  — *float64 ↔ pgtype.Float8
// ---------------------------------------------------------------

func ToNullableFloat(f *float64) pgtype.Float8 {
	if f == nil {
		return pgtype.Float8{Valid: false}
	}
	return pgtype.Float8{Float64: *f, Valid: true}
}

func FromNullableFloat(f pgtype.Float8) *float64 {
	if !f.Valid {
		return nil
	}
	return &f.Float64
}

// ---------------------------------------------------------------
// NULLABLE NUMERIC  — *float64 ↔ pgtype.Numeric (DECIMAL columns)
// ---------------------------------------------------------------

func ToNullableNumeric(f *float64) pgtype.Numeric {
	if f == nil {
		return pgtype.Numeric{Valid: false}
	}
	n := pgtype.Numeric{}
	_ = n.Scan(*f)
	return n
}
