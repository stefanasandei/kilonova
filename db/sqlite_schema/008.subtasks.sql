CREATE TABLE IF NOT EXISTS subtasks (
	id 			INTEGER 	PRIMARY KEY,
	created_at 	TIMESTAMP 	NOT NULL DEFAULT CURRENT_TIMESTAMP,
	problem_id 	INTEGER 	NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
	visible_id 	INTEGER 	NOT NULL,
	score 		INTEGER 	NOT NULL,
	tests 		TEXT 		NOT NULL
);
