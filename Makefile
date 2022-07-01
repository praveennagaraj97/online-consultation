start_dev: stop_dev
	@echo "Starting Development Container"
	@docker-compose -f docker-compose-dev.yml up --build
	@echo "Development Container Running"

stop_dev:
	@echo "Stopping Development Container"
	@docker-compose -f docker-compose-dev.yml down
	@echo "Development Container Stopped"

serve_ngrok:
	@echo "starting ngrok reverse proxy"
	@./ngrok http 4000
	@echo "ngrok reverse proxy server started"