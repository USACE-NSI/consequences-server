echo "Sending request to http://localhost:8080/compute..."

# Run the curl command
# -i is added to show the HTTP response headers (Status 200, etc.)
curl -i -X POST http://localhost:8080/compute \
  -H "Content-Type: application/json" \
  -d @data.json

echo -e "\n\n----------------------------"
echo "Request complete."
read -p "Press [Enter] to close this window..."