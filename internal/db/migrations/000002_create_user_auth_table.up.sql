CREATE TYPE auth_type AS ENUM ('password');

CREATE TABLE IF NOT EXISTS user_auth(
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  auth_type auth_type NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT fk_user_auth_user 
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_auth_user_type_identifier ON user_auth(user_id,auth_type);