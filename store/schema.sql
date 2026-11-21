-- Forums from the forum index, plus sync bookkeeping. Forum ids are unique
-- across sports; topics and posts belong to a sport through their forum.
CREATE TABLE forums (
    id              INTEGER PRIMARY KEY,
    sport           TEXT,               -- MZ's name: "soccer" (football) or "hockey"
    name            TEXT NOT NULL,      -- "ManagerZone talk"
    path            TEXT,               -- "English » ManagerZone talk" (from the topic list)
    lang            TEXT,               -- index section the forum was listed in, e.g. "en"
    position        INTEGER,            -- order on the index page
    description     TEXT,
    total_posts     INTEGER,
    total_topics    INTEGER,
    latest_topic_id INTEGER,
    latest_topic    TEXT,
    latest_post_id  INTEGER,
    latest_post_by  TEXT,
    latest_post_at  TEXT,
    synced_post_id  INTEGER,            -- latest_post_id as of the last completed sync
    synced_at       TEXT
);

-- Forum authors, with their latest known country, federation and roles.
CREATE TABLE users (
    id            INTEGER PRIMARY KEY,
    name          TEXT NOT NULL,
    country       TEXT,
    federation_id INTEGER,
    roles         TEXT                  -- JSON array, e.g. ["MA","Crew"]
);

CREATE TABLE topics (
    id               INTEGER PRIMARY KEY,
    forum_id         INTEGER NOT NULL,
    title            TEXT NOT NULL,     -- full title from the topic page
    pinned           INTEGER NOT NULL DEFAULT 0,
    views            INTEGER,
    replies          INTEGER,           -- posts after the opening one, as MZ counts them
    author_id        INTEGER REFERENCES users(id),
    last_post_id     INTEGER,
    last_post_offset INTEGER,
    last_post_at     TEXT,
    synced_offset    INTEGER NOT NULL DEFAULT 0, -- offset of the last page fetched
    synced_at        TEXT
);
CREATE INDEX topics_forum ON topics(forum_id, last_post_at);

CREATE TABLE posts (
    id        INTEGER PRIMARY KEY,      -- MZ post ids grow over time: ORDER BY id = oldest first
    topic_id  INTEGER NOT NULL REFERENCES topics(id),
    forum_id  INTEGER NOT NULL,
    title     TEXT,                     -- "Re: <topic title>" for replies
    author_id INTEGER REFERENCES users(id),
    team_id   INTEGER,
    posted_at TEXT NOT NULL,
    text      TEXT NOT NULL,            -- the author's own words, quotes removed
    html      TEXT NOT NULL,            -- rendered body including quotes
    quotes    TEXT                      -- JSON [{"Author":..,"Text":..}]
);
CREATE INDEX posts_topic ON posts(topic_id, id);
CREATE INDEX posts_author ON posts(author_id);
CREATE INDEX posts_posted ON posts(posted_at);

-- Full-text index over posts' title and own text (porter stemming: "tactics" matches "tactic").
CREATE VIRTUAL TABLE posts_fts USING fts5(
    title, text, content='posts', content_rowid='id', tokenize='porter unicode61');

CREATE TRIGGER posts_ai AFTER INSERT ON posts BEGIN
    INSERT INTO posts_fts(rowid, title, text) VALUES (new.id, new.title, new.text);
END;
CREATE TRIGGER posts_ad AFTER DELETE ON posts BEGIN
    INSERT INTO posts_fts(posts_fts, rowid, title, text) VALUES ('delete', old.id, old.title, old.text);
END;
CREATE TRIGGER posts_au AFTER UPDATE ON posts BEGIN
    INSERT INTO posts_fts(posts_fts, rowid, title, text) VALUES ('delete', old.id, old.title, old.text);
    INSERT INTO posts_fts(rowid, title, text) VALUES (new.id, new.title, new.text);
END;
