CREATE EXTENSION citext;

CREATE UNLOGGED TABLE Users
(
    Nickname    citext      COLLATE "ucs_basic"  NOT NULL PRIMARY KEY,
    Fullname    varchar(100)      NOT NULL,
    About       text              NOT NULL,
    Email       citext      NOT NULL UNIQUE
);
CREATE INDEX users_nickname ON Users using hash (Nickname);

CREATE UNLOGGED TABLE Forum
(
    Slug         citext             NOT NULL PRIMARY KEY,
    Title        varchar(100)      NOT NULL,
    Nickname     citext           NOT NULL REFERENCES Users(Nickname),
    Posts        int              NOT NULL DEFAULT 0,
    Threads      int               NOT NULL DEFAULT 0
);
CREATE INDEX forum_slug ON Forum using hash (Slug);

CREATE UNLOGGED TABLE Thread
(
    Id           serial            NOT NULL PRIMARY KEY,
    Title        varchar(100)      NOT NULL,
    Author       citext             NOT NULL REFERENCES Users(Nickname),
    Forum        citext              NOT NULL REFERENCES Forum(Slug),
    Message      text              NOT NULL,
    Votes        int               NOT NULL DEFAULT 0,
    Slug         citext,
    Created      timestamp WITH TIME ZONE NOT NULL
);
CREATE INDEX thread_slug ON Thread using hash (Slug);
CREATE INDEX forum_thread ON Thread (Forum, Created);

CREATE UNLOGGED TABLE Posts
(
    Id           serial            NOT NULL PRIMARY KEY,
    Parent       int               NOT NULL DEFAULT 0,
    Author       citext            NOT NULL REFERENCES Users(Nickname),
    Message      text              NOT NULL,
    IsEdited     bool              NOT NULL DEFAULT false,
    Forum        citext            NOT NULL REFERENCES Forum(Slug),
    Thread       serial            NOT NULL REFERENCES Thread(Id),
    Created      timestamp WITH TIME ZONE NOT NULL,
    TreePath     int[]             DEFAULT ARRAY[] :: INT[]
);
CREATE INDEX posts_select ON Posts (Thread, TreePath);
CREATE INDEX posts_select_parent_tree ON Posts ((TreePath[1]), TreePath);

CREATE UNLOGGED TABLE Vote
(
    IdThread     int               NOT NULL REFERENCES Thread(Id),
    Nickname     citext             NOT NULL REFERENCES Users(Nickname),
    Voice        int               NOT NULL DEFAULT 0,
    PRIMARY KEY(IdThread, Nickname)
);

CREATE UNLOGGED TABLE UsersForum
(
    Forum        citext           COLLATE "ucs_basic" NOT NULL REFERENCES Forum(Slug),
    Nickname     citext            NOT NULL REFERENCES Users(Nickname),
    PRIMARY KEY(Forum, Nickname)
);
CREATE INDEX usersforum_nickname ON UsersForum using hash (Nickname);

-- INSERT INTO Users(Nickname, Fullname, About, Email)
-- VALUES ('Test', 'NikitaGureev', 'About 1st user', 'test@mail.ru');

CREATE OR REPLACE FUNCTION update_users_forum() RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO UsersForum(Forum, Nickname)
    VALUES (new.Forum, new.Author)
    ON CONFLICT ON CONSTRAINT usersforum_pkey
    DO NOTHING;
    RETURN new;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER update_users_forum_posts_trigger AFTER INSERT ON Posts FOR EACH ROW EXECUTE PROCEDURE update_users_forum();
CREATE TRIGGER update_users_forum_thread_trigger AFTER INSERT ON Thread FOR EACH ROW EXECUTE PROCEDURE update_users_forum();

CREATE OR REPLACE FUNCTION update_post_path() RETURNS TRIGGER AS $$
BEGIN
    new.TreePath = (SELECT TreePath FROM Posts WHERE id = new.parent) || new.id;
    RETURN new;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER update_path BEFORE INSERT ON Posts FOR EACH ROW EXECUTE PROCEDURE update_post_path();

CREATE OR REPLACE FUNCTION update_post_count() RETURNS TRIGGER AS $$
BEGIN
    UPDATE forum
    SET Posts = forum.Posts + 1
    WHERE Slug = new.Forum;
    RETURN new;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_posts_count_trigger AFTER INSERT ON Posts FOR EACH ROW EXECUTE PROCEDURE update_post_count();

CREATE OR REPLACE FUNCTION update_thread_count() RETURNS TRIGGER AS $$
BEGIN
    UPDATE forum
    SET Threads = forum.Threads + 1
    WHERE Slug = new.Forum;
    RETURN new;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER update_thread_count_trigger AFTER INSERT ON Thread FOR EACH ROW EXECUTE PROCEDURE update_thread_count();

CREATE OR REPLACE FUNCTION update_vote_count() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'UPDATE' THEN
        UPDATE Thread
        SET Votes = Votes - old.Voice + new.Voice
        WHERE Id = new.IdThread;
        RETURN new;
    ELSE
        UPDATE Thread
        SET Votes = Votes + new.Voice
        WHERE Id = new.IdThread;
        RETURN new;
    END IF;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER update_vote_count_trigger AFTER UPDATE OR INSERT ON Vote FOR EACH ROW EXECUTE PROCEDURE update_vote_count();

VACUUM ANALYSE;