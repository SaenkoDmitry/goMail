DROP DATABASE IF EXISTS tarantool_spaces_store;
CREATE DATABASE tarantool_spaces_store;
USE tarantool_spaces_store;

CREATE TABLE IF NOT EXISTS users
(id SERIAL,
name VARCHAR(20) NOT NULL UNIQUE,
hash_password VARCHAR(200) NOT NULL,
PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS spaces
(id SERIAL,
name VARCHAR(50) NOT NULL UNIQUE,
user_id BIGINT UNSIGNED NOT NULL,
PRIMARY KEY (id),
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS history
(id SERIAL,
user_id BIGINT UNSIGNED NOT NULL,
space_id BIGINT UNSIGNED NOT NULL,
command VARCHAR(200) NOT NULL,
result VARCHAR(200) NOT NULL,
PRIMARY KEY (id),
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
FOREIGN KEY (space_id) REFERENCES spaces(id) ON DELETE NO ACTION
);

CREATE TABLE IF NOT EXISTS permissions
(id SERIAL,
user_id BIGINT UNSIGNED NOT NULL,
space_id BIGINT UNSIGNED NOT NULL,
PRIMARY KEY (id),
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
FOREIGN KEY (space_id) REFERENCES spaces(id) ON DELETE CASCADE
);