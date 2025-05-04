local-e2e-test:
	go build -o ../gococo-examples/multi-module-project/gococo.exe main.go
	cd ../gococo-examples/multi-module-project && ./gococo.exe
