
CREATE TABLE IF NOT EXISTS pwd_reset_requests (
	id 			text 		PRIMARY KEY,
	created_at 	timestamptz NOT NULL DEFAULT NOW(),
	user_id 	integer 	NOT NULL REFERENCES users(id) ON DELETE CASCADE
);