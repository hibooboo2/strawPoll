DROP TABLE IF EXISTS polls
CREATE TABLE polls (
    id BIGSERIAL PRIMARY KEY,
    multi_select BOOLEAN,
    per_browser BOOLEAN,
    per_ip BOOLEAN,
    question TEXT NOT NULL
)


DROP TABLE IF EXISTS poll_answer
CREATE TABLE poll_answer (
    id BIGSERIAL PRIMARY KEY,
    answer TEXT NOT NULL,
    poll_id BIGINT NOT NULL REFERENCES polls (id)
)
CREATE INDEX ON poll_answer (id);

DROP TABLE IF EXISTS identities;
CREATE TABLE identities(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255),
    remote_address VARCHAR(500),

)

DROP TABLE IF EXISTS votes;
CREATE TABLE votes(
    id BIGSERIAL PRIMARY KEY,
    poll_answer_id BIGINT  NOT NULL REFERENCES poll_answer (id),
    poll_id BIGINT  NOT NULL REFERENCES polls (id),
    identity_id BIGINT NOT NULL REFERENCES identities (id)
)
