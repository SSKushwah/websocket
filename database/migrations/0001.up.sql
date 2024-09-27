CREATE TABLE IF NOT EXIST USER (
    id              INT PRIMARY KEY NOT NULL,
    name            TEXT NOT NULL,
    age             INT NOT NULL,
    address         CHAR(50),
    salary          CHAR
)