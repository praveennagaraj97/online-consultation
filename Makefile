start_dev: stop_dev
	@echo "Starting Development Container"
	@docker-compose -f docker-compose-dev.yml up --build --remove-orphans
	@echo "Development Container Running"

stop_dev:
	@echo "Stopping Development Container"
	@docker-compose -f docker-compose-dev.yml down
	@echo "Development Container Stopped"


start_prod: stop_dev stop_prod
	@echo "Starting Production Container"
	@docker-compose -f docker-compose.yml up --build --remove-orphans
	@echo "Production Container Running"

stop_prod:
	@echo "Stopping Production Container"
	@docker-compose -f docker-compose.yml down
	@echo "Production Container Stopped"

serve_ngrok:
	@echo "starting ngrok reverse proxy"
	@./ngrok http 4000
	@echo "ngrok reverse proxy server started"