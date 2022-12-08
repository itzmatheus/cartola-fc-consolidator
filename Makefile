docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

migrateup-db:
	migrate -path sql/migrations/ -database "postgresql://consolidator:consolidator@localhost:5432/consolidator?sslmode=disable" -verbose up

migratedown-db:
	migrate -path sql/migrations/ -database "postgresql://consolidator:consolidator@localhost:5432/consolidator?sslmode=disable" -verbose down
