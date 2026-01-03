#!/bin/bash
echo "Starting Web UI on http://localhost:5000 ..."
uv sync --no-install-project
uv run python app.py
