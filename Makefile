export GOPATH=${PWD}

all: demon

exec = hadonis

demon: src/*.go
	go build -o ../dist/usr/local/bin/$(exec) src/*.go

clean:
	rm -f $(exec) *~ *#
	rm -rf pkg
	rm -rf src/gopkg.in

deps:
	go get gopkg.in/yaml.v2

re: clean all
