#include "timer.h" 
#include <signal.h> 
#include <iostream> 

bool should_die  = false; 
bool should_toggle = false; 

void sigint_handler(int signum)
{
	should_die = true; 
}

void sigquit_handler(int signum)
{
	should_toggle = true; 
}

main()
{
	signal(SIGINT, sigint_handler); 
	signal(SIGQUIT, sigquit_handler); 

	Timer timer = {
		.millis = 100, 
		.message = "hello",
		.enabled = false
	}; 

	auto goTimer = InitGoTimer(&timer); 
	timer.enabled = true; 

	StartGoTimer(goTimer); 
	while(!should_die){
		if (should_toggle)
		{
			should_toggle = false;
			timer.enabled = !timer.enabled;
			if(timer.enabled) StartGoTimer(goTimer); 
		}
	}
	std::cout << "finished!" << std::endl; 
}




