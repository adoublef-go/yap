CREATE TABLE `profiles` (
    -- xid format
    `id` TEXT,
    `login` NOT NULL UNIQUE,
    `name` TEXT,
    `photo` TEXT,
    -- constraints
    PRIMARY KEY (id)
);

CREATE TABLE `accounts` (
    `id` TEXT,
    -- xid format
    `profile` TEXT,
    `email` TEXT NOT NULL,
    -- constraints
    CONSTRAINT valid_email CHECK (email <> ''),
    FOREIGN KEY (profile) REFERENCES `profiles` (id),
    PRIMARY KEY (id)
);