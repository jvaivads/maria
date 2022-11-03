create table user
(
    id           int                                  not null,
    user_name    varchar(100)                         not null,
    email        varchar(100)                         not null,
    active       tinyint(1)                           not null,
    date_created datetime default current_timestamp() not null,

    constraint user_pk
        primary key (id)
);

create table role
(
    id           int                                  not null,
    role_name    varchar(100)                         not null,
    type        varchar(100)                          not null,
    active       tinyint(1)                           not null,
    date_created datetime default current_timestamp() not null,

    constraint role_pk
        primary key (id)
);

create table task
(
    id           int                                  not null,
    task_name    varchar(100)                         not null,
    type        varchar(100)                          not null,
    active       tinyint(1)                           not null,
    date_created datetime default current_timestamp() not null,

    constraint task_pk
        primary key (id)
);

create table client
(
    id           int                                  not null,
    client_name    varchar(100)                       not null,
    active       tinyint(1)                           not null,
    date_created datetime default current_timestamp() not null,

    constraint client_pk
        primary key (id)
);

create table user_role
(
    id           int                                  not null,
    user_id    int                       not null,
    role_id       int                           not null,
    date_expired datetime  null,
    date_created datetime default current_timestamp() not null,

    constraint user_role_pk
        primary key (id),
    constraint user_role_user_id_fk
        foreign key (user_id) references user (id),
    constraint user_role_role_id_fk
        foreign key (role_id) references role (id)
);

create table user_client
(
    id           int                                  not null,
    user_id    int                       not null,
    client_id       int                           not null,
    date_expired datetime  null,
    date_created datetime default current_timestamp() not null,

    constraint user_client_pk
        primary key (id),
    constraint user_client_user_id_fk
        foreign key (user_id) references user (id),
    constraint user_client_client_id_fk
        foreign key (client_id) references client (id)
);

create table user_task
(
    id           int                                  not null,
    user_id    int                       not null,
    task_id    int                       not null,
    client_id       int                           not null,
    status       varchar(100)                           not null,
    date_created datetime default current_timestamp() not null,

    constraint user_task_pk
        primary key (id),
    constraint user_task_user_id_fk
        foreign key (user_id) references user (id),
    constraint user_task_task_id_fk
        foreign key (task_id) references task (id),
    constraint user_task_client_id_fk
        foreign key (client_id) references client (id)
);