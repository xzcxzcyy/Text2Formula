cd ./t2f-bot
go mod tidy
go build -o ../t2f-runtime/t2f-bot ./main
cd ..
cd t2f-runtime
./t2f-bot
