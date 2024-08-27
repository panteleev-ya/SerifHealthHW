# Golang parsing and uploading to ClickHouse

## Research
You can find the research that I made in `RESEARCH.md` file

## Why this tech stack?
I used Go for these reasons:
- Serif Health is actively migrating its APIs and backend to the Go.
- Go is high performance programming language.

I used ClickHouse for the storage for these reasons:
- It is a columnar database, which allows highly efficient data access when querying specific columns.
- It is high scalable (horizontally), which makes it a good solution for storing huge amount of data.

## How to run
Make sure you have enough space on disk, it took ~70 GB when I was testing it

1. `docker-compose up -d`
2. Execute `create_table_urls.sql`
3. `go run main.go`
4. Execute `search_for_ny_ppo.sql`

Make sure to delete `storage` directory after the demo is finished

## Pros
- Queries to the ClickHouse work faster than parsing whole JSON
- Script can be executed multiple times with different JSON files to combine the data from all of them
- Different filtering queries can be executed easily using SQL

## Cons
- Data retrieving runs slower than with Python script, so if it needs to be run once Python is better
- Requires disk space

## How long it took me to make?
~2h researching
~2h building the solution

## How long does it run?
~320s (extracting) using local resources
~15-20s (filtering) using local resources
