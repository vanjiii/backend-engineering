#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/socket.h>
#include <netinet/in.h>

// maximum application buffer
// technically the request size
#define APP_MAX_BUFFER 1024
#define PORT 8082

int main () {
    // define the server and client file descriptors
    // one server instance has one server_fd for accepting conns
    // and multiple clients fd for every established conns
    int server_fd, client_fd;

    // define the socket address
    // struct with params like ip, port, etc.
    struct sockaddr_in address;
    int address_len = sizeof(address);

    // define the application buffer where we receive the requests
    // data will be copied from `os received buffer` to this app buffer
    char buffer[APP_MAX_BUFFER] = {0};

    // create the server socket
    // AF_INET - ipv4
    // SOCK_STREAM - streaming protocol, basically TCP
    // 0 - select any streaming protocol (TCP is the only streaming protocol :D)
    if ((server_fd = socket(AF_INET, SOCK_STREAM, 0)) == 0) {
        perror("socket creation failed");

        exit(EXIT_FAILURE);
    }

    // binding socket
    address.sin_family = AF_INET;           // ipv4
    address.sin_addr.s_addr = INADDR_ANY;   // listen on 0.0.0.0
    address.sin_port = htons(PORT);

    if (bind(server_fd, (struct sockaddr *)&address, address_len) < 0) {
        perror("binding socket fails");

        exit(EXIT_FAILURE);
    }

    // create the queues
    // listen for clients with 10 backlog of 10 conns
    //
    // 10 - number of conns that can sit in the queue
    if (listen(server_fd, 10) < 0) {
        perror("listening fails");

        exit(EXIT_FAILURE);
    }


    // handling the request
    while (1) {
        printf("\nwaiting for connections ... \n");

        if ((client_fd = accept(server_fd, (struct sockaddr*)&address, (socklen_t*)&address_len)) < 0) {
            perror("accepting conn failed");

            exit(EXIT_FAILURE);
        }

        read(client_fd, buffer, APP_MAX_BUFFER);
        printf("%s \n", buffer);

        // here we can check what is the request (method, path, headers ...)

        char *http_responce = "HTTP/1.1 200 OK\n"
            "Content-Type: text/plain\n"
            "Content-Length: 13\n\n"
            "Hello World!\n";

        // write back the response to the client
        write(client_fd, http_responce, strlen(http_responce));

        // u don't actually close the conn immediately..!
        close(client_fd);
    }

    return 0;
}
