test:
	go test ./...
lint:
	staticcheck ./...
coverage:
	go test -cover ./...
coverage_file:
	go test -cover -coverprofile=coverage.out ./...
rm_coverage_file:
	rm -f coverage.out
final_coverage:
	cat coverage.out | grep -v "mock" | grep -v "stub" | grep -v "test" > coverage.final.out
coverage_report:
	make final_coverage
	go tool cover -html=coverage.final.out
coverage_project:
	make coverage_file
	make final_coverage
	go tool cover -func coverage.final.out