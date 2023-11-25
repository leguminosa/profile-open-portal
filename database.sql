/**
    This is the SQL script that will be used to initialize the database schema.
    We will evaluate you based on how well you design your database.
    1. How you design the tables.
    2. How you choose the data types and keys.
    3. How you name the fields.
    In this assignment we will use PostgreSQL as the database.
*/

CREATE TABLE users (
    id              SERIAL                                                  not null
        primary key,
    fullname        VARCHAR                                                 not null,
    phone_number    VARCHAR                                                 not null    unique,
    password        TEXT                                                    not null,
    login_count     INTEGER                     default 0                   not null,
    created_at      TIMESTAMP WITH TIME ZONE    default CURRENT_TIMESTAMP   not null,
    updated_at      TIMESTAMP WITH TIME ZONE
);
