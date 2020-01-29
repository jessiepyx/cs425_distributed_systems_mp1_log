To run server:

	go run server/server.go

To run client:

    go run client/client.go <pattern>

To clear up:

	go clean server/*
	go clean client/*

To test:

First generate the test files:

	go run clientTest/sendFile.go
	go run serverTest/receiveFile.go

On the server side, replace logfile.log with serverTest/fakeLog.out

Start servers with

    go run server/server.go

On the client side, run:

    go run client/client.go "true pattern"

Check grep.out if the result contains 1 line of "true pattern"
	

