CREATE TYPE auth_type AS ENUM ('email','google','phone');

CREATE TABLE IF NOT EXISTS user_auth(
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  auth_type auth_type NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  identifier TEXT NOT NULL,
  credential TEXT,
  deleted_at TIMESTAMPTZ,

  CONSTRAINT fk_user_auth_user 
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_auth_user_type ON user_auth(user_id,auth_type) where deleted_at IS NULL;

CREATE UNIQUE INDEX idx_user_auth_identifier ON user_auth(auth_type,identifier) WHERE deleted_at IS NULL;