create table oidc.admin_log
(
    id        int auto_increment
        primary key,
    createdAt datetime default CURRENT_TIMESTAMP not null,
    ip        varchar(36)                        not null,
    ua        varchar(256)                       not null,
    method    varchar(16)                        not null,
    path      varchar(512)                       not null,
    status    varchar(36)                        null,
    error     text                               null,
    userId    varchar(36)                        null,
    cost      double                             null,
    requestId varchar(64)                        not null
)
    collate = utf8mb4_unicode_ci;

create table oidc.admin_menu
(
    id            int auto_increment
        primary key,
    label         varchar(32)                          not null,
    path          varchar(128)                         null,
    icon          longtext                             null,
    sort          int                                  not null,
    level         int                                  not null,
    createdAt     datetime   default CURRENT_TIMESTAMP not null,
    parentId      int                                  null,
    permission    varchar(100)                         null,
    type          char                                 null,
    updatedAt     datetime                             not null,
    code          varchar(64)                          null,
    `schema`      text                                 null,
    visibleInMenu tinyint(1) default 1                 not null
)
    collate = utf8mb4_unicode_ci;

create table oidc.admin_role
(
    id          int auto_increment
        primary key,
    code        varchar(36)                        not null,
    name        varchar(64)                        not null,
    description varchar(255)                       null,
    createdAt   datetime default CURRENT_TIMESTAMP not null,
    updatedAt   datetime                           not null,
    constraint code
        unique (code)
)
    collate = utf8mb4_general_ci;

create table oidc.admin_api_role
(
    id        int auto_increment
        primary key,
    path      varchar(256)                       not null,
    method    varchar(16)                        not null,
    createdAt datetime default CURRENT_TIMESTAMP not null,
    roleId    int                                not null,
    constraint admin_api_role_roleId_fkey
        foreign key (roleId) references oidc.admin_role (id)
            on update cascade on delete cascade
)
    collate = utf8mb4_general_ci;

create table oidc.admin_menu_role
(
    id        int auto_increment
        primary key,
    createdAt datetime default CURRENT_TIMESTAMP not null,
    menuId    int                                not null,
    roleId    int                                not null,
    constraint admin_menu_role_menuId_roleId_key
        unique (menuId, roleId),
    constraint admin_menu_role_menuId_fkey
        foreign key (menuId) references oidc.admin_menu (id)
            on update cascade on delete cascade,
    constraint admin_menu_role_roleId_fkey
        foreign key (roleId) references oidc.admin_role (id)
            on update cascade on delete cascade
)
    collate = utf8mb4_general_ci;

create table oidc.admin_role_user
(
    id        int auto_increment
        primary key,
    createdAt datetime default CURRENT_TIMESTAMP not null,
    roleId    int                                not null,
    userId    varchar(36)                        not null,
    constraint admin_role_user_roleId_fkey
        foreign key (roleId) references oidc.admin_role (id)
            on update cascade on delete cascade
)
    collate = utf8mb4_general_ci;

create table oidc.demo_post_category
(
    id          int auto_increment
        primary key,
    createdAt   datetime default CURRENT_TIMESTAMP not null,
    updatedAt   datetime                           not null,
    name        varchar(64)                        not null,
    description varchar(255)                       null
)
    collate = utf8mb4_unicode_ci;

create table oidc.demo_post
(
    id         int auto_increment
        primary key,
    createdAt  datetime default CURRENT_TIMESTAMP not null,
    updatedAt  datetime                           not null,
    title      varchar(64)                        not null,
    poster     varchar(128)                       null,
    content    text                               not null,
    userId     varchar(36)                        not null,
    categoryId int                                not null,
    constraint demo_post_categoryId_fkey
        foreign key (categoryId) references oidc.demo_post_category (id)
            on update cascade on delete cascade
)
    collate = utf8mb4_unicode_ci;

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
    user_id       varchar(255) null,
    created_at    datetime     null,
    updated_at    datetime     null,
    name          varchar(32)  not null,
    avatar        varchar(255) null,
    password_type varchar(100) null,
    password_salt varchar(100) null,
    password      varchar(100) null,
    phone         char(13)     null,
    country_code  varchar(6)   null
)
    collate = utf8mb4_general_ci;

create table oidc.user_social
(
    id                int auto_increment
        primary key,
    user_id           varchar(36) null,
    provider          varchar(64) not null,
    provider_user_id  varchar(64) not null,
    provider_platform varchar(64) null,
    provider_unionid  varchar(64) null,
    created_at        datetime    not null,
    constraint user_social_pk
        unique (provider_user_id)
);

create table oidc.verification_record
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

