package services

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"testing"
	"time"
)

func BenchmarkCreateUserPgpoolNoLimits(b *testing.B) {
	userService := createUserService(b, true, 0, 0)
	benchmarkCreateUser(userService, b)
}

func BenchmarkCreateUserPgpoolConnectionLimits(b *testing.B) {
	userService := createUserService(b, true, 1, 1)
	benchmarkCreateUser(userService, b)
}

func BenchmarkCreateUserPostgresNoLimits(b *testing.B) {
	userService := createUserService(b, false, 0, 0)
	benchmarkCreateUser(userService, b)
}

func BenchmarkCreateUserPostgresIncreaseIdle(b *testing.B) {
	userService := createUserService(b, false, 0, 3)
	benchmarkCreateUser(userService, b)
}

func BenchmarkCreateUserPostgresIncreaseAll(b *testing.B) {
	userService := createUserService(b, false, 5, 5)
	benchmarkCreateUser(userService, b)
}

func benchmarkCreateUser(userService *PostgresUserService, b *testing.B) {
	user := User{
		ID:         "",
		Name:       "Suzy",
		Occupation: "Worker",
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = userService.Create(&user)
		}
	})
}

func createUserService(t testing.TB, usePgpool bool, maxConnLimits int, maxIdleLimits int) *PostgresUserService {
	var db *sql.DB
	var err error
	if usePgpool {
		db, err = sql.Open("postgres", "postgresql://customuser:custompassword@localhost:5432/postgres?sslmode=disable&fallback_application_name=pgpoolbenchmark")
	} else {
		db, err = sql.Open("postgres", "postgresql://customuser:custompassword@localhost:15432/postgres?sslmode=disable&fallback_application_name=postgresbenchmark")
	}
	if err != nil {
		t.Fatal(err)
	}

	if maxConnLimits > 0 {
		db.SetMaxOpenConns(maxConnLimits)
	}

	if maxIdleLimits > 0 {
		db.SetMaxIdleConns(maxIdleLimits)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	userService := PostgresUserService{db: db}

	return &userService
}
