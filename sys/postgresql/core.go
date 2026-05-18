package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

/*
系统描述:postgresql 数据库驱动系统,基于 jackc/pgx/v5 的 pgxpool 连接池实现

使用方式:
 1. 全局单例:OnInit(config) 初始化默认实例后,直接使用包级函数(Exec/Query/...)
 2. 多实例:  NewSys(SetXxx...) 返回独立 ISys,适合连接多个 PG 实例的场景

参数占位符使用 PostgreSQL 原生风格 $1, $2, ...,而非 ? ;
所有方法均接收 context.Context,调用方需自行控制超时与取消
*/
type (
	// ISys postgresql 系统对外能力集合
	// 既覆盖 pgxpool 常用的 Exec/Query/Tx/Batch/Copy 等原生能力,
	// 也提供 Transaction 闭包等便捷封装
	ISys interface {
		// Pool 返回底层 pgxpool.Pool,用于需要直接操作 pgx 原生 API 的场景
		Pool() *pgxpool.Pool
		// Ping 健康检查:从池中借一条连接执行一次往返
		Ping(ctx context.Context) (err error)
		// Close 关闭连接池,等待所有借出连接归还
		Close()
		// Exec 执行不返回结果集的语句(INSERT/UPDATE/DELETE/DDL)
		Exec(ctx context.Context, sql string, args ...interface{}) (tag pgconn.CommandTag, err error)
		// Query 执行返回多行结果集的查询,rows 用完必须 Close
		Query(ctx context.Context, sql string, args ...interface{}) (rows pgx.Rows, err error)
		// QueryRow 执行单行查询,错误在 Scan 时返回
		QueryRow(ctx context.Context, sql string, args ...interface{}) (row pgx.Row)
		// Begin 以默认隔离级别开启事务,需手动 Commit/Rollback
		Begin(ctx context.Context) (tx pgx.Tx, err error)
		// BeginTx 以指定 TxOptions 开启事务(隔离级别、只读等)
		BeginTx(ctx context.Context, opts pgx.TxOptions) (tx pgx.Tx, err error)
		// SendBatch 单次往返发送一批 SQL,降低 RTT 开销
		SendBatch(ctx context.Context, b *pgx.Batch) (br pgx.BatchResults)
		// CopyFrom 使用 COPY 协议进行批量导入,性能远高于循环 INSERT
		CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (n int64, err error)
		// Transaction 闭包式事务:fn 返回 error 自动回滚,否则自动提交,内部 panic 也会回滚
		Transaction(ctx context.Context, fn func(tx pgx.Tx) error) (err error)
	}
)

var (
	// defsys 全局默认实例,由 OnInit 初始化,包级函数均基于此实例
	defsys ISys
	// ErrNoRows 单行查询无结果时的错误,等价于 pgx.ErrNoRows
	ErrNoRows = pgx.ErrNoRows
	// ErrTxClosed 事务已结束(已 Commit/Rollback)后继续操作时返回
	ErrTxClosed = pgx.ErrTxClosed
)

// OnInit 初始化全局默认实例,通常由框架启动阶段调用
// config 通过 mapstructure 解码到 Options,可再用 option 覆盖
func OnInit(config map[string]interface{}, option ...Option) (err error) {
	defsys, err = newSys(newOptions(config, option...))
	return
}

// NewSys 创建一个独立的 postgresql 实例,不影响全局 defsys
// 适用于需要同时连接多个 PG 实例或多个库的场景
func NewSys(option ...Option) (sys ISys, err error) {
	sys, err = newSys(newOptionsByOption(option...))
	return
}

// 以下为操作全局默认实例的包级快捷函数,使用前请确保已调用过 OnInit

func Pool() *pgxpool.Pool             { return defsys.Pool() }
func Ping(ctx context.Context) error  { return defsys.Ping(ctx) }
func Close()                          { defsys.Close() }
func Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return defsys.Exec(ctx, sql, args...)
}
func Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return defsys.Query(ctx, sql, args...)
}
func QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return defsys.QueryRow(ctx, sql, args...)
}
func Begin(ctx context.Context) (pgx.Tx, error) { return defsys.Begin(ctx) }
func BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
	return defsys.BeginTx(ctx, opts)
}
func SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return defsys.SendBatch(ctx, b)
}
func CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return defsys.CopyFrom(ctx, tableName, columnNames, rowSrc)
}
func Transaction(ctx context.Context, fn func(tx pgx.Tx) error) error {
	return defsys.Transaction(ctx, fn)
}
