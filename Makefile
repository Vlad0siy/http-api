.PHONY: run
run:
	go run -v ./main.go

.PHONY: kill
kill:
	taskkill /F /IM main.exe
