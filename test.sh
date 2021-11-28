sleep 5
if curl -x POST "http://localhost:8081" -d "account_id='1'&{}" | grep -q 'Success'; then
  echo "Tests passed!"
  exit 0
else
  echo "Tests failed!"
  exit 1
