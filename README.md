Summary

This is the first project for the course, which we have called "NU Chitter". It should serve as both an introduction to Go as well as a self-evaluation of proficiency in networking and systems programming. All submissions should be 100% your own work and will be graded on a Pass/Fail scale -- either the submitted code works or it does not.

NU Chitter is a very simple chat client and server with support for private messaging. By default, incoming messages should be broadcast to all connected clients. Incoming messages that contain a colon are considered "commands". There are a total of three commands "whoami:", "all:", and private messages. More details on these commands are provided below.

Note: Your code *must* be thread safe! You should not be sharing objects across goroutines. Code that is not thread safe will be very heavily penalized (50% off). Furthermore, you *must* use channels to communicate between goroutines! Submissions that use locking mechanisms will also be very heavily penalized.

Turning in

For this project, you should be able to fit all your code within a single file -- the reference implementation is relatively short (less than 150 lines of code with comments). Projects 1-3 will require more than a single file. For those, we will provide instructions on how to pack them into a submittable tar.gz file. However, for this project, just submit your single source code file named "chitter.go".

Project description

The specification of this project is relatively straight forward. Remember, it's basically a chatroom that has support for private messages (PMs). The first client acts as a server for the room and other clients connect to it.

Your NU Chitter implementation should require a port number when starting. NU Chitter will listen for connections on this port. Below is an example of starting NU Chitter in the command line.
go run chitter.go 12345
By default, messages should be broadcasted to all clients. In other words, any messages that does not contain a colon should be broadcast. Furthermore, incoming messages that begin with "all:" should also be broadcast (more on commands below). For example, if a client sends any of the following messages (one per line), the body of the message should be forwarded to all clients. In the examples below, the body of each message is in bold text ("Hello, world!").
Hello, world!
all: Hello, world!
all : Hello, world!
When messages are forwarded to other clients, they should be prefixed with the ID of the user that sent the message. In our example above, the "Hello, world!" messages from a user with an ID of 42 would appear as follows:
42: Hello, world!
Commands: Your NU Chitter program will need to support (at least) three types of "commands" that are parsed out by the server. This includes the "all:" command, the "whoami:" command, and sending PMs. If a user sends an invalid command, you can ignore the message.
all: - The "all:" command allows a user to explicitly broadcast a message. This also allows a user to use a colon ":" in their message. Remember that by default, if NU Chitter receives a message with one or more colons in it, everything before the first colon is interpreted as a command. Using the "all:" command allows users to have colons in their messages.
whoami: - Users should be able to get their own ID number by using the whoami: command. Below is an example of user 13 asking what their ID is with the NU Chitter responding only to that user.
whoami:
chitter: 13
Note that anything after a "whoami:" command can simply be ignored by the server.
Private messages - Users should also be able to send a message to a specific user. To do so, their message should be prefixed with the ID of the user they wish to message. Private messages to nonexistent users can just be ignored (i.e., the server does nothing). If user 1 wants to message user 2, their submitted message should be in the following format:
2: Hey, number two!
And should appear as the following on user 2's (and only user 2's) connection:
1: Hey, number two!
Debugging & testing

Testing your own NU Chitter implementation should be straightforward. Additionally, a reference solution will be running on aqualab.cs.northwestern.edu on port 10345. Note that this server is only reachable from Northwestern's network. If you are trying to access it from another network, you will need to connect to Northwestern's VPN service. In order to help simplify testing, we have also implemented a "who" command on our test server. Similar to it's Unix counterpart, "who" returns a list of all the users connected to the server. You do not need to implement this function (though it should be trivial).

Your program also MUST be thread safe. We will be stress testing it to make sure that there are not any incorrect behaviors. You can check for race conditions by using the following command when you test your server:

go run -race chitter.go 12345
Suggested Plan of attack

Below is an outline of our suggest plan of attack. If you are already comfortable with Go, the following can still serve as a useful guide for the suggested order in which to implement each feature.

Write a simple echo server in Go. Make sure the server is able to handle concurrent connections (hint: use Goroutines).
Modify the server so that incoming messages are forwarded to all connected clients. The server will need to maintain a record of clients that are currently connected to the server. You will need to begin using Go's channels and the select statement to do so.
Create a client interface to connect to the server.  The client will serve the same functionality as connecting with telnet or nc.  It will be started by the same program with a different flag passed in, i.e. "-c", to indicate it's a client as well as the IP:port number of the server.
Assign unique IDs to the connected clients. Our reference implementation assigns integers (0, 1, 2, etc.), but you can assign IDs however you prefer -- just do not assign the IDs "all" and "whoami", those are reserved for commands. Once users have IDs, prefix outgoing messages with the ID of the user that sent the message.
Add support for the "whoami" and "all" commands. You can assume any line of text with a colon contains a command, and that all text before the colon is the command. If a user sends an invalid command, you can ignore the message. If there is no colon, it is a broadcast message. Go's "strings" (Links to an external site.) package should make parsing messages from users straightforward.
Add support for PMs. Clients should be able to send a PM by prefacing a message with the desired user's ID and a colon (e.g., 1: <message>).
