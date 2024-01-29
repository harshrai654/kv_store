CREATE TABLE IF NOT EXISTS kv(
  k varchar(30) NOT NULL,
  v blob NOT NULL,
  expired_at timestamp NULL DEFAULT NULL,
  PRIMARY KEY (k),
  KEY idx_expired_at (expired_at)
);

CREATE EVENT KeyDeleteJob
ON SCHEDULE EVERY 5 MINUTE
DO
  DELETE FROM kv WHERE expired_at <= NOW();