-- create app tables
CREATE TABLE IF NOT EXISTS accounts(
	id UUID PRIMARY KEY,
	name VARCHAR(512),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS users(
	id UUID PRIMARY KEY,
	email VARCHAR(320) NOT NULL,
	name VARCHAR(512) NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS account_users(
	account_id UUID NOT NULL REFERENCES accounts(id),
	user_id UUID NOT NULL REFERENCES users(id),
	PRIMARY KEY (account_id, user_id)
);


CREATE TABLE IF NOT EXISTS toggles(
	id UUID PRIMARY KEY,
	account_id UUID NOT NULL REFERENCES accounts(id),
	key VARCHAR(512) NOT NULL,
	active BOOLEAN NOT NULL DEFAULT TRUE,
	rules JSONB,
	description VARCHAR(2048),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (account_id, key)
);

-- create app user
DO
$do$
BEGIN
	IF NOT EXISTS (
		SELECT FROM pg_catalog.pg_user
		WHERE usename = 'toggle') THEN

		CREATE USER toggle WITH ENCRYPTED PASSWORD 'toggle';
		REVOKE CONNECT ON DATABASE toggle FROM PUBLIC;

		GRANT CONNECT
		ON DATABASE toggle
		TO toggle;

		GRANT SELECT, INSERT, UPDATE, DELETE
		ON ALL TABLES IN SCHEMA public
		TO toggle;

	END IF;
END
$do$;


-- add some test data
INSERT INTO accounts (id, name) VALUES
('8dc8c3cd-7c2a-4a4c-bc1e-7ba042096029', 'Toggle Test')
ON CONFLICT DO NOTHING;
