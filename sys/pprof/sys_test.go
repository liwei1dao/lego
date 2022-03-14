package pprof_test

/*
查看堆栈调用信息
go tool pprof http://localhost:6060/debug/pprof/heap
查看 30 秒内的 CPU 信息
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
查看 goroutine 阻塞
go tool pprof http://localhost:6060/debug/pprof/block
收集 5 秒内的执行路径
go tool pprof http://localhost:6060/debug/pprof/trace?seconds=5
争用互斥持有者的堆栈跟踪
go tool pprof http://localhost:6060/debug/pprof/mutex
*/
/*
UI web 界面
curl -sK -v http://localhost:6060/debug/pprof/heap > heap.out
go tool pprof -http=:8080 heap.out
*/
