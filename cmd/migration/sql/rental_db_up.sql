CREATE TABLE IF NOT EXISTS orders (
    id TEXT NOT NULL PRIMARY KEY,
    "dateFrom" timestamp, -- datetime
    "dateTo" timestamp, -- datetime
    "dateReservFrom" timestamp NOT NULL, -- datetime
    "dateReservTo" timestamp NOT NULL, -- datetime
    status INTEGER NOT NULL,
    "stationFromId" TEXT NOT NULL,
    "stationToId" TEXT NOT NULL,
    discount REAL,
    tax REAL
);

CREATE TABLE IF NOT EXISTS ocars (
    id TEXT NOT NULL,
    "orderId" TEXT NOT NULL,
    age TEXT NOT NULL,
    plate TEXT NOT NULL,
    document TEXT NOT NULL,
    "carModel" TEXT NOT NULL,
    "initialKM" INTEGER NOT NULL,
    "finalKM" INTEGER NOT NULL,
    status INTEGER NOT NULL,
    "stationId" TEXT NOT NULL,
    FOREIGN KEY ("orderId") REFERENCES orders(id),
    PRIMARY KEY (id, "orderId")
);

CREATE TABLE IF NOT EXISTS opolicies (
    id TEXT NOT NULL,
    "orderId" TEXT NOT NULL,
    name TEXT NOT NULL,
    price REAL NOT NULL,
    unit INTEGER NOT NULL,
    "minUnit" INTEGER NOT NULL,
    "carModel" TEXT NOT NULL,
    "categoryId" TEXT NOT NULL,
    FOREIGN KEY ("orderId") REFERENCES orders(id),
    PRIMARY KEY (id, "orderId")
);