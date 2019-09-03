all: 
	go build -buildmode=c-shared -o trial.so trial.go
	g++ -o main main.cpp trial.so
