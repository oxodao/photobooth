run:
	mkdir -p gui/dist gui_admin/dist
	touch gui/dist/random_data gui_admin/dist/random_data # Make go fckng happy
	docker compose up -d
	go run .

create-db:
	rm -rf 0_DATA/photobooth.db
	sqlite3 0_DATA/photobooth.db < sql/init.sql

reset-db:
	rm -rf 0_DATA/photobooth.db
	sqlite3 0_DATA/photobooth.db < sql/init.sql
	sqlite3 0_DATA/photobooth.db < sql/fixtures.sql

take-picture:
	docker compose exec mosquitto mosquitto_pub -h 127.0.0.1 -t photobooth/button_press -m "TAKE_PICTURE"

show-debug:
	docker compose exec mosquitto mosquitto_pub -h 127.0.0.1 -t photobooth/button_press -m "DISPLAY_DEBUG"

set-mode-photobooth:
	docker compose exec mosquitto mosquitto_pub -h 127.0.0.1 -t photobooth/admin/set_mode -m "PHOTOBOOTH"

set-mode-quiz:
	docker compose exec mosquitto mosquitto_pub -h 127.0.0.1 -t photobooth/admin/set_mode -m "QUIZ"

set-mode-disabled:
	docker compose exec mosquitto mosquitto_pub -h 127.0.0.1 -t photobooth/admin/set_mode -m "DISABLED"
