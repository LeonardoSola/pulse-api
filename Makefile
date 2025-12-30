.PHONY: help migrate migrate-up migrate-down migrate-status migrate-create

# Run the application in development mode
run:
	bash -c 'set -a && source .env && set +a && air .'

# Show available commands
help:
	@echo "Available commands:"
	@echo "  make run          - Run the application in development mode"
	@echo "  make help         - Show available commands"