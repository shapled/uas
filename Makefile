build:
	go build -o dist/uas.exe

run:
	cd dist && .\uas.exe -c test.yml
