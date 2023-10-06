# Simple-Real-Time-Chat-App

## Hub :
➜ The manager of the rooms <br>
➜ Running on a separate go routine rather than the main routine <br>


The websocket connection of gorillamux actually allows one concurrent writer at a time 

browser actually responds to the ping messages by default, all we need to do is to handle the server-side code 
-> this what happend when a client connected and suddenly disconnected, the server will now by catching the "writemsg: websocket: close sent"
```cmd
2023/10/06 10:11:32 [Heartbeat] | ping
2023/10/06 10:11:32 pong
2023/10/06 10:11:41 [Heartbeat] | ping
2023/10/06 10:11:41 pong
2023/10/06 10:11:50 [Heartbeat] | ping
2023/10/06 10:11:50 pong
2023/10/06 10:11:59 [Heartbeat] | ping
2023/10/06 10:11:59 pong
2023/10/06 10:12:08 [Heartbeat] | ping
2023/10/06 10:12:08 pong
2023/10/06 10:12:15 error ==> websocket: close 1000 (normal)
2023/10/06 10:12:17 [Heartbeat] | ping
2023/10/06 10:12:17 writemsg:  websocket: close sent
```