CREATE TABLE IF NOT EXISTS users (
	id 			INTEGER 	PRIMARY KEY,
	created_at	TIMESTAMP 	NOT NULL DEFAULT CURRENT_TIMESTAMP,
	name 		TEXT 		NOT NULL UNIQUE,
	admin 		INTEGER 	NOT NULL DEFAULT FALSE,
	proposer 	INTEGER 	NOT NULL DEFAULT FALSE,
	email 		TEXT 		NOT NULL UNIQUE,
	password 	TEXT 		NOT NULL,
	bio 		TEXT 		NOT NULL DEFAULT '',

	default_visible INTEGER NOT NULL DEFAULT FALSE,

	verified_email 	INTEGER NOT NULL DEFAULT FALSE,
	email_verif_sent_at TIMESTAMP,

	banned 		INTEGER 	NOT NULL DEFAULT FALSE,
	disabled 	INTEGER 	NOT NULL DEFAULT FALSE,

	preferred_language TEXT NOT NULL DEFAULT 'ro'
);
