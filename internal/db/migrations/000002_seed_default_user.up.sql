INSERT INTO users (id, username, nickname, email)
VALUES (1, 'admin', 'Administrator', 'admin@example.com')
ON CONFLICT (username) DO NOTHING;

SELECT setval(
    pg_get_serial_sequence('users', 'id'),
    COALESCE((SELECT MAX(id) FROM users), 1),
    true
);
