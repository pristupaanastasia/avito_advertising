CREATE TABLE advertisement (
        id         SERIAL         NOT NULL,
        price           int         NOT NULL,
        name     CHAR(200)   NOT NULL,
        description     CHAR(1000)   NOT NULL,
        image               char(100)[]         NOT NULL,
        update          date        NOT NULL,
        PRIMARY KEY (id)
);