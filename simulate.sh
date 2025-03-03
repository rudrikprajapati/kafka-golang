#!/bin/bash
echo "Starting 500 user simulation..."
start_time=$(date +%s)

for i in {1..500}
do
  curl -X POST -H "Content-Type: application/json" \
    -d "{\"user_id\": \"user$i\", \"data\": \"data$i\"}" \
    http://localhost:8080/upload &
done
wait

end_time=$(date +%s)
duration=$((end_time - start_time))
echo "Simulation completed in $duration seconds"