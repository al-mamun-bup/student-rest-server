PORT?=8080
serve:
	@go build -o student && ./student serve --port=${PORT}