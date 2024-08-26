# Makefile

# Define the command for running Air
run-air:
	@echo "Starting Air for Go server..."
	air

# Define the command for running Tailwind in watch mode
watch-tailwind:
	@echo "Starting Tailwind CSS in watch mode..."
	npx tailwindcss -i ./static/css/styles.css -o ./static/css/tailwind.css --watch

# Define a command to run both concurrently
dev: 
	@echo "Running Air and Tailwind CSS together..."
	$(MAKE) -j2 run-air watch-tailwind

.PHONY: run-air watch-tailwind dev
