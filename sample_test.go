package main

import (
	"context"
	"database/sql"
	"github.com/ericlagergren/decimal"
	"github.com/jwilner/sqlboiler-bug/models"
	_ "github.com/lib/pq" // pull in driver
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/types"
	"os"
	"testing"
)

type other struct {
	models.Example `boil:",bind"`
	SomeVal        null.Int64 `boil:"val"`
}

func TestModel(t *testing.T) {
	r := require.New(t)
	ctx := context.Background()

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_DSN"))
	r.NoError(err)

	x1 := &models.Example{
		FieldA: types.NewNullDecimal(decimal.New(4096, 0)),
		FieldB: types.NewNullDecimal(decimal.New(4096, 0)),
		FieldC: types.NewNullDecimal(decimal.New(4096, 0)),
	}

	x2 := &models.Example{
		FieldA: types.NewNullDecimal(decimal.New(256, 0)),
		FieldB: types.NewNullDecimal(decimal.New(256, 0)),
		FieldC: types.NewNullDecimal(decimal.New(256, 0)),
	}

	r.NoError(x1.Insert(ctx, db, boil.Infer()))
	r.NoError(x2.Insert(ctx, db, boil.Infer()))

	var exes []*other
	r.NoError(queries.Raw("SELECT * FROM public.example").Bind(ctx, db, &exes))
	r.Len(exes, 2)

	if exes[0].ID != x1.ID {
		exes[0], exes[1] = exes[1], exes[0]
	}

	r.Equal(int64(4096), forceInt64(exes[0].FieldA))
	r.Equal(int64(4096), forceInt64(exes[0].FieldB))
	r.Equal(int64(4096), forceInt64(exes[0].FieldC))
	r.Equal(int64(256), forceInt64(exes[1].FieldA))
	r.Equal(int64(256), forceInt64(exes[1].FieldB))
	r.Equal(int64(256), forceInt64(exes[1].FieldC))
}

func forceInt64(x types.NullDecimal) int64 {
	i, _ := x.Int64()
	return i
}
