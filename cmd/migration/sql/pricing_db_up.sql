-- PRAGMA foreign_keys = on;

CREATE TABLE IF NOT EXISTS categories (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS cmodels (
    id SERIAL PRIMARY KEY, -- INTEGER AUTOINCREMENT
    name TEXT NOT NULL,
    "categoryId" TEXT NOT NULL,
    FOREIGN KEY ("categoryId") REFERENCES categories(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS cpolicies (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    price FLOAT NOT NULL,
    unit INT NOT NULL,
    "minUnit" INT NOT NULL,
    "categoryId" TEXT NOT NULL,
    FOREIGN KEY ("categoryId") REFERENCES categories(id) ON DELETE CASCADE
);