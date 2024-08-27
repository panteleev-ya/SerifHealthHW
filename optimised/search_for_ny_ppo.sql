SELECT description, location
FROM urls
WHERE (lower(description) LIKE '%new york%' or lower(description) LIKE '%ny%')
  AND lower(description) LIKE '%ppo%';