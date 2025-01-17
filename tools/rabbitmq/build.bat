set GOOS=linux
set CGO_ENABLED=0

go build -o ./rabbitmq.a11 ./main.go

@REM go build -o ./bin/go_manwu/a11/st005.a11 ./services/st005/main.go
@REM cd ./bin
@REM xcopy /E /I /Y .\json .\docker_manwugame\json
REM pause