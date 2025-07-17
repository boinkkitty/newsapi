CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS news (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    author TEXT NOT NULL,
    title TEXT NOT NULL,
    summary TEXT NOT NULL,
    content TEXT NOT NULL,
    source TEXT NOT NULL,
    tags TEXT[] NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

INSERT INTO news(id, author, title, summary, content, source, tags, created_at, updated_at)
VALUES(
    uuid_generate_v4(),
    'Alice Smith',
    'PostgreSQL 16 Released',
    'What''s new in PostgreSQL 16?',
    'PostgreSQL 16 brings performance and monitoring improvements...',
    'postgresweekly.com',
    ARRAY['postgresql', 'release', 'database'],
    NOW(),
    NOW()
);
INSERT INTO news(id, author, title, summary, content, source, tags, created_at, updated_at)
VALUES(
    uuid_generate_v4(),
    'Bob Lee',
    'Understanding UUIDs in Databases',
    'Why UUIDs are useful for distributed systems.',
    'UUIDs help prevent key collisions across systems, especially in microservices...',
    'dbtips.io',
    ARRAY['uuid', 'database', 'design'],
    NOW(),
    NOW()
);

INSERT INTO news (
    id, author, title, summary, content, source, tags, created_at, updated_at, deleted_at
) VALUES (
     uuid_generate_v4(),
     'Jane Doe',
     'Go Makes Testing Easy',
     'A short summary about Go testing',
     'This article discusses how Go makes writing tests easier by providing built-in support for testing.',
     'example.com',
     ARRAY['go', 'testing', 'backend'],
     NOW(),
     NOW(),
     NOW()
 );

