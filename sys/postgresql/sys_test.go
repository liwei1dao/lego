package postgresql_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/liwei1dao/lego/sys/postgresql"
)

const testDsn = "postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable"

// 测试 系统初始化 与 Ping
func Test_sys_Init(t *testing.T) {
	sys, err := postgresql.NewSys(
		postgresql.SetPostgresqlUrl(testDsn),
		postgresql.SetMaxConns(8),
		postgresql.SetMinConns(1),
	)
	if err != nil {
		t.Skipf("postgres not available, skip: %v", err)
		return
	}
	defer sys.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	if err = sys.Ping(ctx); err != nil {
		t.Fatalf("ping failed: %v", err)
	}
}

// 测试 增删改查 与 事务
func Test_sys_CRUD(t *testing.T) {
	sys, err := postgresql.NewSys(postgresql.SetPostgresqlUrl(testDsn))
	if err != nil {
		t.Skipf("postgres not available, skip: %v", err)
		return
	}
	defer sys.Close()

	ctx := context.Background()
	if _, err = sys.Exec(ctx, `CREATE TABLE IF NOT EXISTS sys_pg_test (id BIGSERIAL PRIMARY KEY, name TEXT NOT NULL, age INT)`); err != nil {
		t.Fatalf("create table: %v", err)
	}
	defer sys.Exec(ctx, `DROP TABLE IF EXISTS sys_pg_test`)

	if _, err = sys.Exec(ctx, `INSERT INTO sys_pg_test(name, age) VALUES ($1, $2)`, "alice", 18); err != nil {
		t.Fatalf("insert: %v", err)
	}

	var name string
	var age int
	if err = sys.QueryRow(ctx, `SELECT name, age FROM sys_pg_test WHERE name = $1`, "alice").Scan(&name, &age); err != nil {
		t.Fatalf("queryrow: %v", err)
	}
	fmt.Printf("queryrow => name=%s age=%d\n", name, age)

	err = sys.Transaction(ctx, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `INSERT INTO sys_pg_test(name, age) VALUES ($1, $2)`, "bob", 20); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `UPDATE sys_pg_test SET age = age + 1 WHERE name = $1`, "alice"); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatalf("transaction: %v", err)
	}

	rows, err := sys.Query(ctx, `SELECT name, age FROM sys_pg_test ORDER BY id`)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&name, &age); err != nil {
			t.Fatalf("scan: %v", err)
		}
		fmt.Printf("row => name=%s age=%d\n", name, age)
	}
}
