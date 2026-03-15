build:
	go build -o pv ./cmd/pv

install:
	cp pv /usr/local/bin/

clean:
	rm -f pv

run:
	go run ./cmd/pv
