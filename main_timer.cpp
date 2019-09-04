#include "timer.h" 
#include <signal.h> 
#include <iostream> 

bool should_die  = false; 

void sigint_handler(int signum)
{
	should_die = true; 
}

main()
{
	signal(SIGINT, sigint_handler); 
	StartTimer(100); 
	while(!should_die) {} 
	StopTimer(); 
	std::cout << "finished!" << std::endl; 
}




