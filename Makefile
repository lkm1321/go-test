ctimer: main_ctimer.cpp timer.h timer.so
	g++ -o ctimer main_ctimer.cpp ./timer.so

timer: main_timer.cpp timer.h timer.so
	g++ -o timer main_timer.cpp ./timer.so

timer.so: timer.go
	go build -buildmode=c-shared -o timer.so timer.go

timer.h: timer.go
	go build -buildmode=c-shared -o timer.so timer.go

trial: trial.so main.cpp
	g++ -o main main.cpp trial.so

trial.so: trial.go
	go build -buildmode=c-shared -o trial.so trial.go
