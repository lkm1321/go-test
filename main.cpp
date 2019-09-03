#include <iostream>
#include <string> 
#include "trial.h" 

GoString toGoString(std::string cppstring)
{
	return {cppstring.c_str(), cppstring.size()}; 
}

int main()
{
	std::string message = "hello world"; 
	std::cout << Log(toGoString(message));
	std::cout << Log(toGoString(message)); 
	return 0; 
}
