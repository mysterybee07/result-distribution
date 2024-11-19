run:
	@docker run -p 8080:8080 --rm -v $(pwd):/app -v /app/tmp --name result-distribution-air result-distribution
