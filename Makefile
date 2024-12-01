run:
	air
build:
	templ generate
	npx tailwindcss -i ./static/css/tailwind.css -o ./static/css/tailwind.min.css
	go build -o bin/main cmd/main.go 

migrate-up:
	DB_USER=sagra DB_SERVER=localhost DB_PASSWORD=sagra DB_NAME=sagra_go sql-migrate up -config=config/dbconfig.yml
	DB_USER=sagra DB_SERVER=localhost DB_PASSWORD=sagra sql-migrate up -config=config/dbconfig.yml -env test

migrate-down:
	DB_USER=sagra DB_SERVER=localhost DB_PASSWORD=sagra DB_NAME=sagra_go sql-migrate down -config=config/dbconfig.yml
	DB_USER=sagra DB_SERVER=localhost DB_PASSWORD=sagra  sql-migrate down -config=config/dbconfig.yml -env test
migrate-fresh:
	DB_USER=sagra DB_SERVER=localhost DB_PASSWORD=sagra DB_NAME=sagra_go sql-migrate down -config=config/dbconfig.yml
	DB_USER=sagra DB_SERVER=localhost DB_PASSWORD=sagra sql-migrate down -config=config/dbconfig.yml -env test
	DB_USER=sagra DB_SERVER=localhost DB_PASSWORD=sagra DB_NAME=sagra_go sql-migrate up -config=config/dbconfig.yml
	DB_USER=sagra DB_SERVER=localhost DB_PASSWORD=sagra sql-migrate up -config=config/dbconfig.yml -env test
