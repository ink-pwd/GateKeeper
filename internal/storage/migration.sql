-- public."User" определение

-- Drop table

-- DROP TABLE public."User";

CREATE TABLE public."User" (
	id int4 GENERATED ALWAYS AS IDENTITY NOT NULL,
	email varchar NOT NULL,
	password_hash varchar NOT NULL,
	CONSTRAINT user_pk PRIMARY KEY (id),
	CONSTRAINT user_unique UNIQUE (email)
);