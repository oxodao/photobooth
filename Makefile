run:
	go run .

create-db:
	rm -rf 0_DATA/photomaton.db
	sqlite3 0_DATA/photomaton.db < sql/init.sql

reset-db:
	rm -rf 0_DATA/photomaton.db
	sqlite3 0_DATA/photomaton.db < sql/init.sql
	sqlite3 0_DATA/photomaton.db < sql/fixtures.sql
