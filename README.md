# Mandatory Hand-In 3 (Chitty-Chat Distributed System in Go)

Repository for assignment 3 by the group "Cat Squish Gang".

## Example Output

```
// Server
2022/10/20 15:25:55 Server listening at 127.0.0.1:50005
2022/10/20 15:26:07 Client 0 connected to the server...
2022/10/20 15:26:12 Client 1 connected to the server...
2022/10/20 15:26:17 Client 0 wrote the message: hello world!
2022/10/20 15:26:26 Client 1 wrote the message: hello world! to you too!
2022/10/20 15:26:31 Client 2 connected to the server...
2022/10/20 15:26:38 Client 2 wrote the message: hello world! to both of you!
2022/10/20 15:26:43 Client 0 disconnected from the server...
2022/10/20 15:26:46 Client 1 disconnected from the server...
2022/10/20 15:26:54 Client 2 disconnected from the server...

// Client 1
2022/10/20 15:26:07 Client connecting to server...
2022/10/20 15:26:07 Client connected successfully...
2022/10/20 15:26:07 Participant 0 joined Chitty-Chat at Lamport time 0
2022/10/20 15:26:07 ================================
                        Welcome to Chitty Chat!

                     Commands:
                     --help            Display Help
                     --exit            Leave Server
                     --change-lobby X  Change Lobby
                    ================================
2022/10/20 15:26:12 Participant 1 joined Chitty-Chat at Lamport time 1
hello world!
2022/10/20 15:26:26 Participant 1 wrote at Lamport time 5: hello world! to you too!
2022/10/20 15:26:31 Participant 2 joined Chitty-Chat at Lamport time 6
2022/10/20 15:26:38 Participant 2 wrote at Lamport time 7: hello world! to both of you!
--exit
2022/10/20 15:26:43 Closing server connection...
2022/10/20 15:26:43 Server connection ended...

// Client 2
2022/10/20 15:26:12 Client connecting to server...
2022/10/20 15:26:12 Client connected successfully...
2022/10/20 15:26:12 Participant 1 joined Chitty-Chat at Lamport time 0
2022/10/20 15:26:12 ================================
                        Welcome to Chitty Chat!

                     Commands:
                     --help            Display Help
                     --exit            Leave Server
                     --change-lobby X  Change Lobby
                    ================================
2022/10/20 15:26:17 Participant 0 wrote at Lamport time 3: hello world!
hello world! to you too!
2022/10/20 15:26:31 Participant 2 joined Chitty-Chat at Lamport time 5
2022/10/20 15:26:38 Participant 2 wrote at Lamport time 6: hello world! to both of you!
2022/10/20 15:26:43 Participant 0 left Chitty-Chat at Lamport time 9
--exit
2022/10/20 15:26:46 Closing server connection...
2022/10/20 15:26:46 Server connection ended...

// Client 3
2022/10/20 15:26:31 Client connecting to server...
2022/10/20 15:26:31 Client connected successfully...
2022/10/20 15:26:31 Participant 2 joined Chitty-Chat at Lamport time 0
2022/10/20 15:26:31 ================================
                        Welcome to Chitty Chat!

                     Commands:
                     --help            Display Help
                     --exit            Leave Server
                     --change-lobby X  Change Lobby
                    ================================
hello world! to both of you!
2022/10/20 15:26:43 Participant 0 left Chitty-Chat at Lamport time 9
2022/10/20 15:26:46 Participant 1 left Chitty-Chat at Lamport time 11
Curabitur lobortis sit amet massa sed pretium. Fusce eu metus maximus, malesuada sapien eu, pharetra erat. Proin sit amet nullam.
2022/10/20 15:26:52 Failed to send message: more than 128 characters
--exit
2022/10/20 15:26:54 Closing server connection...
2022/10/20 15:26:54 Server connection ended...

```
