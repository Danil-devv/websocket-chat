CREATE TABLE IF NOT EXISTS messages (
  id BIGSERIAL,
  username CHARACTER VARYING(128) NOT NULL,
  data TEXT NOT NULL
)