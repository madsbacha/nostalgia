CREATE TABLE IF NOT EXISTS main.media
(
    id            TEXT              not null primary key,
    title         TEXT              not null,
    file_id       TEXT              not null,
    thumbnail_id  INTEGER           not null,
    uploaded_at   INTEGER           not null,
    uploaded_by   INTEGER           not null,
    description   TEXT              not null,
    tags          TEXT              not null
);

CREATE TABLE IF NOT EXISTS main.user
(
    id            INTEGER           not null primary key,
    discord_id    TEXT              not null unique,
    username      TEXT              not null,
    avatar        TEXT,
    roles         TEXT              not null
);

CREATE TABLE IF NOT EXISTS main.setting
(
    id            INTEGER           not null primary key,
    key           TEXT              not null unique,
    value         TEXT              not null
);

CREATE TABLE IF NOT EXISTS main.file
(
    id            INTEGER           not null primary key,
    path          TEXT              not null,
    mime_type     TEXT              not null
);

CREATE TABLE IF NOT EXISTS main.thumbnail
(
    id            INTEGER           not null primary key,
    file_id       TEXT              not null,
    blurhash      TEXT              not null
);
