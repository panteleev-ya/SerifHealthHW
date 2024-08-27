CREATE TABLE IF NOT EXISTS urls
(
    id          UUID,
    description String,
    location    String
) ENGINE = MergeTree()
      ORDER BY id;
