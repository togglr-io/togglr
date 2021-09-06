-- create update trigger
CREATE OR REPLACE FUNCTION updated_at_trigger()
RETURNS TRIGGER AS $$
BEGIN
	NEW.updated_at = now();
	RETURN NEW;
END;
$$ language 'plpgsql';

-- create app tables
CREATE TABLE IF NOT EXISTS accounts(
	id UUID PRIMARY KEY,
	name VARCHAR(512),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

DROP TRIGGER IF EXISTS accounts_updated_at ON accounts;
CREATE TRIGGER accounts_updated_at BEFORE UPDATE
ON accounts FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();


CREATE TABLE IF NOT EXISTS identity_types(
	name VARCHAR(64) PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users(
	id UUID PRIMARY KEY,
	email VARCHAR(320) NOT NULL,
	name VARCHAR(512) NOT NULL,
	identity VARCHAR(512) NOT NULL,
	identity_type VARCHAR(64) NOT NULL REFERENCES identity_types(name),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (identity, identity_type)
);

DROP TRIGGER IF EXISTS users_updated_at ON users;
CREATE TRIGGER users_updated_at BEFORE UPDATE
ON users FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();



CREATE TABLE IF NOT EXISTS account_users(
	account_id UUID NOT NULL REFERENCES accounts(id),
	user_id UUID NOT NULL REFERENCES users(id),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
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

DROP TRIGGER IF EXISTS toggles_updated_at ON toggles;
CREATE TRIGGER toggles_updated_at BEFORE UPDATE
ON toggles FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();



CREATE TABLE IF NOT EXISTS metadata_keys(
	id UUID PRIMARY KEY,
	account_id UUID NOT NULL REFERENCES accounts(id),
	key VARCHAR(512) NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (account_id, key)
);

DROP TRIGGER IF EXISTS metadata_keys_updated_at ON metadata_keys;
CREATE TRIGGER metadata_keys_updated_at BEFORE UPDATE
ON metadata_keys FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();



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


-- initial data migration
INSERT INTO identity_types (name) VALUES
('github'),
('google'),
('basic')
ON CONFLICT DO NOTHING;

-- add some test data
INSERT INTO accounts (id, name) VALUES
('8dc8c3cd-7c2a-4a4c-bc1e-7ba042096029', 'Toggle Test')
ON CONFLICT DO NOTHING;


