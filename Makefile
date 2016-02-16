build:
	go build -o yaml-go-json main.go yaml-to-json.go

install:
	cp yaml-go-json /usr/local/bin/
