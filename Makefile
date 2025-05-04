local-e2e-test:
	go build -o ../gococo-examples/multi-module-project/gococo.exe ./cmd/
	cd ../gococo-examples/multi-module-project && ./gococo.exe
	