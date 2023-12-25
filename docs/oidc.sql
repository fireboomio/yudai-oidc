create table oidc.admin_api_role
(
    id           int auto_increment
        primary key,
    path         int                                not null,
    method       varchar(16)                        not null,
    role_id      int                                not null,
    created_time datetime default CURRENT_TIMESTAMP not null
)
    collate = utf8mb4_general_ci;

create table oidc.admin_log
(
    id        int auto_increment
        primary key,
    createdAt datetime(3) default CURRENT_TIMESTAMP(3) not null,
    ip        varchar(191)                             not null,
    ua        varchar(191)                             not null,
    method    varchar(191)                             not null,
    path      varchar(191)                             not null,
    status    varchar(36)                              null,
    error     text                                     null,
    userId    varchar(36)                              null
)
    collate = utf8mb4_unicode_ci;

create table oidc.admin_menu
(
    id          int auto_increment
        primary key,
    parent_id   int                                        null,
    label       varchar(32)                                not null,
    path        varchar(100)                               null,
    icon        longtext                                   null,
    sort        int                                        not null,
    level       int                                        not null,
    name        varchar(32)                                null,
    is_bottom   tinyint unsigned default '1'               not null,
    create_time datetime         default CURRENT_TIMESTAMP not null,
    perms       varchar(100)                               null,
    menu_type   char             default ''                null
)
    collate = utf8mb4_unicode_ci;

create table oidc.admin_menu_role
(
    id           int auto_increment
        primary key,
    menu_id      int                                not null,
    role_id      int                                not null,
    created_time datetime default CURRENT_TIMESTAMP not null
)
    collate = utf8mb4_general_ci;

create table oidc.admin_role
(
    id           int auto_increment
        primary key,
    code         varchar(36)                        not null,
    name         varchar(64)                        not null,
    description  varchar(255)                       null,
    created_time datetime default CURRENT_TIMESTAMP not null,
    updated_time datetime                           not null,
    constraint code
        unique (code)
)
    collate = utf8mb4_general_ci;

create table oidc.admin_role_user
(
    id           int auto_increment
        primary key,
    role_id      int                                not null,
    user_id      varchar(36)                        not null,
    created_time datetime default CURRENT_TIMESTAMP not null
)
    collate = utf8mb4_general_ci;

create table oidc.provider
(
    owner         varchar(100)  not null,
    name          varchar(100)  not null,
    created_time  varchar(100)  null,
    type          varchar(100)  null,
    client_id     varchar(100)  null,
    client_secret varchar(2000) null,
    sign_name     varchar(100)  null,
    template_code varchar(100)  null,
    primary key (owner, name),
    constraint UQE_provider_name
        unique (name)
)
    collate = utf8mb4_general_ci;

create table oidc.token
(
    id                  int auto_increment
        primary key,
    user_id             varchar(255)           null,
    token               longtext               null,
    expire_time         datetime               null,
    flush_time          datetime               null,
    refresh_token       longtext               null,
    refresh_expire_time datetime               null,
    banned              tinyint(1)             null,
    platform            varchar(36) default '' not null
)
    collate = utf8mb4_general_ci;

create table oidc.user
(
    id            int auto_increment
        primary key,
    created_at    datetime     null,
    name          varchar(32)  not null,
    avatar        varchar(255) null,
    phone         char(13)     null,
    password_salt varchar(100) null,
    password      varchar(100) null,
    country_code  varchar(6)   null,
    password_type varchar(100) null,
    user_id       varchar(255) null,
    wx_unionid    varchar(100) null
)
    collate = utf8mb4_general_ci;

create table oidc.userwx
(
    created_at varchar(100) null,
    unionid    varchar(100) not null,
    openid     varchar(100) not null,
    platform   varchar(36)  null,
    primary key (unionid, openid)
);

create table oidc.verificationrecord
(
    name         varchar(100) not null
        primary key,
    created_time varchar(100) null,
    type         varchar(10)  null,
    user         varchar(100) not null,
    provider     varchar(100) not null,
    receiver     varchar(100) not null,
    code         varchar(10)  not null,
    time         bigint       not null,
    is_used      tinyint(1)   null
)
    collate = utf8mb4_general_ci;

