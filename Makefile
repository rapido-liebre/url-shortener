# Build and start all containers
up:
	docker-compose up --build

# Stop and remove all containers and volumes
down:
	docker-compose down -v

# Show backend logs
logs:
	docker-compose logs -f backend

# Open PostgreSQL shell inside container
psql:
	docker exec -it $$(docker ps -qf "name=db") psql -U postgres -d shortener

# Test POST /links/shorten via curl
curl:
	curl -X POST http://localhost:8080/links/shorten \
	  -H "Content-Type: application/json" \
	  -d '{"long_url":"https://example.org"}' -v

# Start React frontend in local dev mode (with Vite)
frontend-dev:
	cd frontend && npm run dev

# Reset everything: clean volumes and rebuild
reset-db:
	docker-compose down -v && docker-compose up --build
