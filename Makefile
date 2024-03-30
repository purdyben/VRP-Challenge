build: 
	go build -o bin/main cmd/main.go 
	chmod +x bin/main

run: build
	./bin/main $(args)
# 	go run cmd/main.go $(args)

util:
	go mod tidy -v
	go mod verify
	go fmt vorto/...
	go vet cmd/main.go


eval: build 
	python3 evaluateShared.py --cmd ./bin/main --problemDir trainingProblems

eval1: build 
	python3 evaluateShared.py --cmd ./bin/main --problemDir trainingProblems1
