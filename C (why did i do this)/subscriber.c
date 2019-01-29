#include <unistd.h>
#include <stdio.h>
#include <sys/socket.h>
#include <stdlib.h>
#include <netinet/in.h>
#include <string.h>
#include <sys/types.h>

int main(int argc, char const *argv[])
{
  if (argc < 1) {
    perror("Please specify the port where the socket should be activated\nHow to run the program: ./StartChat [port number]");
    exit(EXIT_FAILURE);
  }

  int portNumber = atoi(argv[1]);
  if (portNumber < 1 || portNumber > 65535) {
    perror("Port number is between 1 and 65535 inclusive");
    exit(EXIT_FAILURE);
  }

  int clientSocket = socket(PF_INET, SOCK_STREAM, 0);
  if (clientSocket < 0) {
    perror("Socket can not be created");
    exit(EXIT_FAILURE);
  }

  
  return 0;
}
