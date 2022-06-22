CREATE EXTENSION citext;

CREATE UNLOGGED TABLE Users
(
    Nickname    citext      COLLATE "ucs_basic"  NOT NULL PRIMARY KEY,
    Fullname    varchar(100)      NOT NULL,
    About       text              NOT NULL,
    Email       citext      NOT NULL
);

CREATE UNLOGGED TABLE Forum
(
    Slug         citext             NOT NULL PRIMARY KEY,
    Title        varchar(100)      NOT NULL,
    Nickname     citext           NOT NULL REFERENCES Users(Nickname) ON DELETE CASCADE,
    Posts        int              NOT NULL DEFAULT 0,
    Threads      int               NOT NULL DEFAULT 0
);

CREATE UNLOGGED TABLE Thread
(
    Id           serial            NOT NULL PRIMARY KEY,
    Title        varchar(100)      NOT NULL,
    Author       citext             NOT NULL REFERENCES Users(Nickname) ON DELETE CASCADE,
    Forum        citext              NOT NULL REFERENCES Forum(Slug) ON DELETE CASCADE,
    Message      text              NOT NULL,
    Votes        int               NOT NULL DEFAULT 0,
    Slug         citext,
    Created      timestamp WITH TIME ZONE NOT NULL
);
CREATE INDEX thread_select_users ON Thread (Forum, Author);

CREATE UNLOGGED TABLE Posts
(
    Id           serial            NOT NULL PRIMARY KEY,
    Parent       int               NOT NULL DEFAULT 0,
    Author       citext            NOT NULL REFERENCES Users(Nickname) ON DELETE CASCADE,
    Message      text              NOT NULL,
    IsEdited     bool              NOT NULL DEFAULT false,
    Forum        citext            NOT NULL REFERENCES Forum(Slug) ON DELETE CASCADE,
    Thread       serial            NOT NULL REFERENCES Thread(Id) ON DELETE CASCADE,
    Created      timestamp WITH TIME ZONE NOT NULL,
    TreePath     int[]             DEFAULT ARRAY[] :: INT[]
);
CREATE INDEX posts_select_users ON Posts (Forum, Author);
CREATE INDEX posts_select ON Posts (Thread, TreePath);

CREATE UNLOGGED TABLE Vote
(
    IdThread     int               NOT NULL REFERENCES Thread(Id) ON DELETE CASCADE,
    Nickname     citext             NOT NULL REFERENCES Users(Nickname) ON DELETE CASCADE,
    Voice        int               NOT NULL DEFAULT 0,
    PRIMARY KEY(IdThread, Nickname)
);

-- INSERT INTO Users(Nickname, Fullname, About, Email)
-- VALUES ('Test', 'NikitaGureev', 'About 1st user', 'test@mail.ru');

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
