#include <unistd.h>
#include <stdio.h>
#include <sys/socket.h>
#include <stdlib.h>
#include <netinet/in.h>
#include <string.h>
#include <sys/types.h>
#include <arpa/inet.h>
#include <sys/time.h>
#include <errno.h>

#define PORT 1337
#define MAX_SUBSCRIBERS 5
#define TRUE 1

int main(int argc, char const *argv[])
{
  int opt = TRUE, clientNumber = 0;
  int subscriber[MAX_SUBSCRIBERS];
  struct sockaddr_in subscriberAddresses[MAX_SUBSCRIBERS];

  int serverSocket = socket(PF_INET, SOCK_STREAM, 0);
  fd_set readfds;
  if (serverSocket < 0)
  {
    perror("Problems during creating socket");
    exit(EXIT_FAILURE);
  }

  for (int i = 0; i < MAX_SUBSCRIBERS; i++)
    subscriber[i] = 0;

  //set master socket to allow multiple connections ,
  //this is just a good habit, it will work without this
  if (setsockopt(serverSocket, SOL_SOCKET, SO_REUSEADDR, (char *)&opt,
                 sizeof(opt)) < 0)
  {
    perror("setsockopt");
    exit(EXIT_FAILURE);
  }

  struct sockaddr_in addr;
  addr.sin_family = PF_INET;
  addr.sin_addr.s_addr = INADDR_ANY;
  addr.sin_port = htons(PORT);

  int bindStatus = bind(serverSocket, (const struct sockaddr *)&addr, sizeof(addr));
  if (bindStatus < 0)
  {
    perror("Problems during binding");
    exit(EXIT_FAILURE);
  }

  int listenStatus = listen(serverSocket, MAX_SUBSCRIBERS);
  if (listenStatus < 0)
  {
    perror("Can not listen");
    exit(EXIT_FAILURE);
  }

  int clientServer = -1, newSocket = -1;
  int addrLen = sizeof(addr);
  struct sockaddr_in *currentAddress = NULL;
  int maxSocketDescriptor = 0, socketDescriptor = 0, activity = 0;
  while (1)
  {
    //clear the socket set
    FD_ZERO(&readfds);

    //add master socket to set
    FD_SET(serverSocket, &readfds);
    maxSocketDescriptor = serverSocket;

    //add child sockets to set
    for (int i = 0; i < MAX_SUBSCRIBERS; i++)
    {
      //socket descriptor
      socketDescriptor = subscriber[i];

      //if valid socket descriptor then add to read list
      if (socketDescriptor > 0)
        FD_SET(socketDescriptor, &readfds);

      //highest file descriptor number, need it for the select function
      if (socketDescriptor > maxSocketDescriptor)
        maxSocketDescriptor = socketDescriptor;
    }

    //wait for an activity on one of the sockets , timeout is NULL ,
    //so wait indefinitely
    activity = select(maxSocketDescriptor + 1, &readfds, NULL, NULL, NULL);

    if ((activity < 0) && (errno != EINTR))
    {
      printf("select error");
    }

    //If something happened on the master socket ,
    //then its an incoming connection
    if (FD_ISSET(serverSocket, &readfds))
    {
      if ((newSocket = accept(serverSocket,
                              (struct sockaddr *)&addr, (socklen_t *)&addrLen)) < 0)
      {
        perror("accept");
        exit(EXIT_FAILURE);
      }

      //inform user of socket number - used in send and receive commands
      printf("New connection , socket fd is %d , ip is : %s , port : %d\n " ,
        newSocket , inet_ntoa(addr.sin_addr) , ntohs(addr.sin_port));
      subscriberAddresses[clientNumber++] = addr;

      recv()
    }
  }

  close(serverSocket);

  return 0;
}
