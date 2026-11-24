SELECT result
FROM cryptocurrency_runs
WHERE convert_id = ? AND coin_ids = ? AND status = ?
ORDER BY id DESC
LIMIT 1;
