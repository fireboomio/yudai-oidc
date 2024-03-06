create table provider
(
    owner         varchar(100)  not null,
    name          varchar(100)  not null,
    created_at    datetime      not null,
    type          varchar(100)  null,
    client_id     varchar(100)  not null,
    client_secret varchar(2000) not null,
    sign_name     varchar(100)  null,
    template_code varchar(100)  null,
    primary key (owner, name),
    constraint UQE_provider_name
        unique (name)
)
    charset = utf8mb4;

create table token
(
    id                  int auto_increment
        primary key,
    platform            varchar(36) null,
    user_id             varchar(36) not null,
    token               text        not null,
    created_at          datetime    not null,
    expire_time         datetime    null,
    refresh_token       text        null,
    refresh_expire_time datetime    null,
    banned              tinyint(1)  null
)
    charset = utf8mb4;

create table user
(
    id            int auto_increment
        primary key,
    user_id       varchar(36)  not null,
    avatar        varchar(255) null,
    created_at    datetime     not null,
    updated_at    datetime     null,
    name          varchar(64)  null,
    password      varchar(100) null,
    password_salt varchar(100) null,
    password_type varchar(100) null,
    phone         varchar(20)  null,
    country_code  varchar(6)   null
)
    charset = utf8mb4;

create index IDX_user_name
    on user (name);

create index IDX_user_phone
    on user (phone);

create table user_social
(
    id                int auto_increment
        primary key,
    user_id           varchar(36) null,
    provider          varchar(64) not null,
    provider_user_id  varchar(64) not null,
    provider_platform varchar(64) null,
    provider_unionid  varchar(64) null,
    created_at        datetime    null
)
    charset = utf8mb4;

create index IDX_user_social_created_at
    on user_social (created_at);

create table verification_record
(
    id         int auto_increment
        primary key,
    created_at datetime     not null,
    type       varchar(10)  null,
    user_id    varchar(36)  not null,
    provider   varchar(100) not null,
    receiver   varchar(100) not null,
    code       varchar(10)  not null,
    is_used    tinyint(1)   null
)
    charset = utf8mb4;