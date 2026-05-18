package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// newSys 根据传入的 Options 构造并初始化一个 Postgresql 实例
// 返回的 *Postgresql 已完成连接池建立与首次 Ping 检测,可直接使用
func newSys(options Options) (sys *Postgresql, err error) {
	sys = &Postgresql{options: options}
	err = sys.init()
	return
}

// Postgresql 基于 pgxpool 的 PostgreSQL 驱动封装
//   - options: 选项配置(连接串、连接池参数、日志等)
//   - pool:    pgx 提供的协程安全的连接池,所有 SQL 操作均经由该池获取连接
type Postgresql struct {
	options Options
	pool    *pgxpool.Pool
}

// init 完成连接池的构建:
//  1. 通过 pgxpool.ParseConfig 解析 DSN(支持 URL 和 key=value 两种格式)
//  2. 使用 Options 中显式配置的连接池参数覆盖默认值(0 值视为不覆盖)
//  3. 在 ConnectTimeout 上下文中真正建立连接池并进行 Ping,失败则返回错误
func (this *Postgresql) init() (err error) {
	var cfg *pgxpool.Config
	if cfg, err = pgxpool.ParseConfig(this.options.PostgresqlUrl); err != nil {
		this.options.Log.Errorln(err)
		return
	}
	// 仅当用户显式配置(>0)时才覆盖,避免把默认值清零
	if this.options.MaxConns > 0 {
		cfg.MaxConns = this.options.MaxConns
	}
	if this.options.MinConns > 0 {
		cfg.MinConns = this.options.MinConns
	}
	if this.options.MaxConnLifetime > 0 {
		cfg.MaxConnLifetime = this.options.MaxConnLifetime
	}
	if this.options.MaxConnIdleTime > 0 {
		cfg.MaxConnIdleTime = this.options.MaxConnIdleTime
	}
	if this.options.HealthCheckPeriod > 0 {
		cfg.HealthCheckPeriod = this.options.HealthCheckPeriod
	}
	if this.options.ConnectTimeout > 0 {
		cfg.ConnConfig.ConnectTimeout = this.options.ConnectTimeout
	}

	// 用 ConnectTimeout 限制建池+首次 Ping 总耗时,防止启动阶段长时间阻塞
	ctx, cancel := context.WithTimeout(context.Background(), this.options.ConnectTimeout)
	defer cancel()
	if this.pool, err = pgxpool.NewWithConfig(ctx, cfg); err != nil {
		this.options.Log.Errorln(err)
		return
	}
	if err = this.pool.Ping(ctx); err != nil {
		this.options.Log.Errorf("postgresql ping failed: %v", err)
		return
	}
	return
}

// Pool 暴露底层 *pgxpool.Pool,供需要直接使用 pgx 原生能力的调用方使用
// (例如自定义类型注册、ScanRow、AcquireFunc 等高级用法)
func (this *Postgresql) Pool() *pgxpool.Pool {
	return this.pool
}

// Ping 通过连接池获取一条连接并执行一次往返,用于健康检查
func (this *Postgresql) Ping(ctx context.Context) (err error) {
	return this.pool.Ping(ctx)
}

// Close 关闭连接池,等待所有已借出的连接归还并释放底层 TCP 资源
// 进程退出时应当调用,重复调用是安全的(pool 为 nil 时直接返回)
func (this *Postgresql) Close() {
	if this.pool != nil {
		this.pool.Close()
	}
}

// Exec 执行不返回结果集的 SQL(INSERT / UPDATE / DELETE / DDL 等)
// 返回的 CommandTag 可用于读取受影响的行数
// 参数占位符使用 PostgreSQL 风格:$1, $2, ...
func (this *Postgresql) Exec(ctx context.Context, sql string, args ...interface{}) (tag pgconn.CommandTag, err error) {
	tag, err = this.pool.Exec(ctx, sql, args...)
	return
}

// Query 执行返回多行结果集的查询
// 调用方必须在使用完毕后调用 rows.Close() 归还连接,否则会导致连接泄漏
func (this *Postgresql) Query(ctx context.Context, sql string, args ...interface{}) (rows pgx.Rows, err error) {
	rows, err = this.pool.Query(ctx, sql, args...)
	return
}

// QueryRow 执行只期望返回单行的查询
// 错误(包括 ErrNoRows)会延迟到 Scan 时才返回,调用方无需手动关闭
func (this *Postgresql) QueryRow(ctx context.Context, sql string, args ...interface{}) (row pgx.Row) {
	row = this.pool.QueryRow(ctx, sql, args...)
	return
}

// Begin 以默认隔离级别开启一个事务,需调用方手动 Commit/Rollback
// 若只是简单事务,优先使用 Transaction 闭包以避免遗漏回滚
func (this *Postgresql) Begin(ctx context.Context) (tx pgx.Tx, err error) {
	tx, err = this.pool.Begin(ctx)
	return
}

// BeginTx 以指定的隔离级别/访问模式开启一个事务
// 适用于需要 ReadOnly、Serializable 等特殊语义的场景
func (this *Postgresql) BeginTx(ctx context.Context, opts pgx.TxOptions) (tx pgx.Tx, err error) {
	tx, err = this.pool.BeginTx(ctx, opts)
	return
}

// SendBatch 通过单次网络往返发送一批 SQL 语句,显著降低高并发小语句的 RTT 开销
// 调用方需对返回的 BatchResults 依次 Exec/Query/QueryRow,并在最后 Close
func (this *Postgresql) SendBatch(ctx context.Context, b *pgx.Batch) (br pgx.BatchResults) {
	br = this.pool.SendBatch(ctx, b)
	return
}

// CopyFrom 使用 PostgreSQL COPY 协议进行批量数据导入,性能远高于循环 INSERT
//   - tableName:   目标表名(可带 schema,例如 pgx.Identifier{"public", "users"})
//   - columnNames: 写入的列名顺序需与 rowSrc.Values() 一致
//   - rowSrc:      数据源,通常使用 pgx.CopyFromRows / pgx.CopyFromSlice 构造
func (this *Postgresql) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (n int64, err error) {
	n, err = this.pool.CopyFrom(ctx, tableName, columnNames, rowSrc)
	return
}

// Transaction 以闭包形式执行一个事务,fn 返回 error 时自动回滚,否则自动提交
// 流程:
//  1. Begin 失败直接返回
//  2. fn 内 panic 时 recover 并回滚,err 被替换为对应的 panic 信息
//  3. fn 返回 error 时回滚,回滚自身的错误仅记录日志不覆盖业务错误
//  4. fn 正常返回则 Commit,Commit 错误作为最终 err 返回
func (this *Postgresql) Transaction(ctx context.Context, fn func(tx pgx.Tx) error) (err error) {
	var tx pgx.Tx
	if tx, err = this.pool.Begin(ctx); err != nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(ctx)
			err = fmt.Errorf("postgresql transaction panic: %v", r)
		}
	}()
	if err = fn(tx); err != nil {
		// ErrTxClosed 说明事务已经结束(例如 fn 内部已显式提交/回滚),忽略即可
		if rbErr := tx.Rollback(ctx); rbErr != nil && rbErr != pgx.ErrTxClosed {
			this.options.Log.Errorf("postgresql transaction rollback failed: %v", rbErr)
		}
		return
	}
	err = tx.Commit(ctx)
	return
}
