# 

PGPASSWORD=password psql -U postgres -h localhost -p 5432 -d homework

# The query logic 
host_000008,2017-01-01 08:59:22,2017-01-01 09:59:22

# Finished query
SELECT
    time_bucket('1 minute', ts) AS minute,
    host,
    max(usage) AS max_usage,
    min(usage) AS min_usage
  FROM cpu_usage
  WHERE ts >= '2017-01-01 08:59:22' AND ts < '2017-01-01 09:59:22'
  GROUP BY minute, host
  ORDER BY minute DESC;

SELECT
    time_bucket('1 minute', ts) AS minute,
    host,
    max(usage) AS max_usage,
    min(usage) AS min_usage
  FROM cpu_usage
  WHERE ts >= '2017-01-01 08:59:22' AND ts < '2017-01-01 09:59:22'
  AND host = 'host_000017'
  GROUP BY minute, host
  ORDER BY minute DESC;