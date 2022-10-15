CREATE TABLE IF NOT EXISTS cars (
    id TEXT NOT NULL PRIMARY KEY,
    age INT NOT NULL,
    plate TEXT NOT NULL,
    document TEXT NOT NULL,
    model TEXT NOT NULL,
    make TEXT NOT NULL,
    "stationId" TEXT NOT NULL,
    km INT NOT NULL,
    status INT NOT NULL
);

CREATE TABLE IF NOT EXISTS stations (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    complement TEXT NOT NULL,
    state TEXT NOT NULL,
    city TEXT NOT NULL,
    cep TEXT NOT NULL,
    capacity INT NOT NULL,
    idle INT NOT NULL
);