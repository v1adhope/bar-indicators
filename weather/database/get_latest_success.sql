SELECT result
FROM weather_runs
WHERE city = ? AND status = ?
ORDER BY id DESC
LIMIT 1;
