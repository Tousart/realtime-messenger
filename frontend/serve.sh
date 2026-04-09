#!/usr/bin/env bash
# Serves the frontend on http://localhost:3000
# Requires Python 3 (preinstalled on most systems)
cd "$(dirname "$0")"
echo "Frontend available at http://localhost:3000"
echo "  Register: http://localhost:3000/auth/register.html"
echo "  Login:    http://localhost:3000/auth/login.html"
python3 -m http.server 3000
