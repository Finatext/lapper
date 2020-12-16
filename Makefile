build: lapper/build examples/go/build examples/shell/build


test:
	go test

lapper/build:
	docker build -t lapper:latest .

examples/%/build:
	docker build -t lapper-example-$*:latest examples/$*/

examples/%/run:
	docker run --env-file .env --rm -p 9000:8080 lapper-example-$*:latest
