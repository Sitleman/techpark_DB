CREATE TABLE Users
(
    Nickname    varchar(100)    NOT NULL PRIMARY KEY,
    Fullname    varchar(100)      NOT NULL,
    About       text      NOT NULL,
    Email       varchar(100)      NOT NULL
);

INSERT INTO Users(Nickname, Fullname, About, Email)
VALUES ('Test', 'NikitaGureev', 'About 1st user', 'test@mail.ru');

CREATE TABLE Forum
(
    Slug         varchar(100)      NOT NULL PRIMARY KEY,
    Title        varchar(100)      NOT NULL,
    Nickname     varchar(100)      NOT NULL REFERENCES Users(Nickname) ON DELETE CASCADE,
    Posts        int               NOT NULL DEFAULT 0,
    Threads      int               NOT NULL DEFAULT 0
);

CREATE TABLE Thread
(
    Id           serial            NOT NULL PRIMARY KEY,
    Title        varchar(100)      NOT NULL,
    Author       varchar(100)      NOT NULL REFERENCES Users(Nickname) ON DELETE CASCADE,
    Forum        varchar(100)      NOT NULL REFERENCES Forum(Slug) ON DELETE CASCADE,
    Message      text      NOT NULL,
    Votes        int               NOT NULL DEFAULT 0,
    Slug         varchar(100)      ,
    Created      timestamp WITH TIME ZONE NOT NULL
);