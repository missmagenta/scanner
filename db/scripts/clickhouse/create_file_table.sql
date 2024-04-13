CREATE TABLE IF NOT EXISTS "file" (
                        id UUID DEFAULT generateUUIDv4(),
                        filename String,
                        status String,
                        date_processed Nullable(DATETIME),
                        error Nullable(String)
) ENGINE = MergeTree()
PRIMARY KEY (id)
ORDER BY id;