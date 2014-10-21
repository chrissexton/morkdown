build: server.go
	go build
	bower install

run: build morkdown
	-pkill morkdown
	./morkdown&
