REQUESTS=${1:-500}  # Default to 500 if no argument
echo "Starting $REQUESTS user simulation..."
start_time=$(date +%s)

for i in $(seq 1 $REQUESTS)
do
  curl -X POST -H "Content-Type: application/json" \
    -d "{\"user_id\": \"user$i\", \"data\": \"data$i\"}" \
    http://localhost:8080/upload &
done
wait

end_time=$(date +%s)
duration=$((end_time - start_time))
echo "Simulation completed in $duration seconds"