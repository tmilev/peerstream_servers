# peerstream_servers

Repository for solving peerstream challenges.

Repository source:

https://github.com/tmilev/peerstream_server

## Ubuntu instructions
1. To build:

```
cd client
go build
cd ../server
go build
```

2. To run:

2.1 Run client
```
cd client
./client
```

2.2 To run 1 copy of the server:

```
cd server
./server
```

Press enter for default. 
The default auto-selects the first port available 
between localhost:9000 and localhost:9010.
