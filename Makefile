start_dev: stop_dev
	@echo "Starting Development Container"
	@docker-compose -f docker-compose-dev.yml up --build
	@echo "Development Container Running"

stop_dev:
	@echo "Stopping Development Container"
	@docker-compose -f docker-compose-dev.yml down
	@echo "Development Container Stopped"