-- +goose Up

create table users (
    created_at    timestamp       not null default now(),
    id            bigint(20)      not null primary key,
    lang          varchar(2)      not null default '',
    blocked       boolean         not null default 0
);

create table links (
    created_at    timestamp       not null default now(),
    id            bigint(64)      not null primary key auto_increment,
    user_id       bigint(20)      not null default 0 references users (id),
    url           varchar(1024)   not null default '' collate utf8mb4_unicode_ci,
    deleted       boolean         not null default 0,

    index (user_id)
);

create table views (
    created_at    timestamp       not null default now(),
    id            bigint(64)      not null primary key auto_increment,
    link_id       bigint(64)      not null default 0 references links (id),
    ip            varchar(16)     not null default '',
    user_agent    varchar(256)    not null default '',

    index (link_id)
);
