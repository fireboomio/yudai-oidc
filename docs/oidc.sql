create table admin_log
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
    requestId varchar(64)                        not null,
    body      text                               null,
    constraint admin_log_userId_fkey
        foreign key (userId) references user (user_id)
            on update cascade on delete cascade
)
    collate = utf8mb4_general_ci;

create table admin_menu
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
    visibleInMenu tinyint(1) default 1                 not null,
    apis          varchar(191)                         null
)
    collate = utf8mb4_general_ci;

create table admin_role
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

create table admin_menu_role
(
    id        int auto_increment
        primary key,
    createdAt datetime default CURRENT_TIMESTAMP not null,
    menuId    int                                not null,
    roleId    int                                not null,
    constraint admin_menu_role_menuId_roleId_key
        unique (menuId, roleId),
    constraint admin_menu_role_menuId_fkey
        foreign key (menuId) references admin_menu (id)
            on update cascade on delete cascade,
    constraint admin_menu_role_roleId_fkey
        foreign key (roleId) references admin_role (id)
            on update cascade on delete cascade
)
    collate = utf8mb4_general_ci;

create table admin_role_user
(
    id        int auto_increment
        primary key,
    createdAt datetime default CURRENT_TIMESTAMP not null,
    roleId    int                                not null,
    userId    varchar(36)                        not null,
    constraint admin_role_user_roleId_userId_key
        unique (roleId, userId),
    constraint admin_role_user_roleId_fkey
        foreign key (roleId) references admin_role (id)
            on update cascade on delete cascade
)
    collate = utf8mb4_general_ci;

create table demo_area
(
    id         int auto_increment
        primary key,
    createdAt  datetime    default CURRENT_TIMESTAMP not null,
    name       varchar(64)                           not null,
    address    text                                  null,
    code       varchar(36)                           not null,
    parentCode varchar(36) default ''                not null,
    constraint demo_area_pk2
        unique (code)
);

create table demo_area_role
(
    id        int auto_increment
        primary key,
    createdAt datetime default CURRENT_TIMESTAMP not null,
    areaId    int                                not null,
    roleId    int                                not null,
    constraint demo_area_role_areaId_roleId_key
        unique (areaId, roleId),
    constraint demo_area_role_areaId_fkey
        foreign key (areaId) references demo_area (id)
            on update cascade on delete cascade,
    constraint demo_area_role_roleId_fkey
        foreign key (roleId) references admin_role (id)
            on update cascade on delete cascade
)
    collate = utf8mb4_general_ci;

create table demo_post_category
(
    id          int auto_increment
        primary key,
    createdAt   datetime default CURRENT_TIMESTAMP not null,
    updatedAt   datetime                           not null,
    name        varchar(64)                        not null,
    description varchar(255)                       null
)
    collate = utf8mb4_general_ci;

create table demo_post
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
        foreign key (categoryId) references demo_post_category (id)
            on update cascade on delete cascade
)
    collate = utf8mb4_general_ci;
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