dev:
	DB_URL="postgres://postgres@localhost:5432/wongnok?sslmode=disable" go run .

build:
	docker build -t acoshift/wongnok .
	docker push acoshift/wongnok
