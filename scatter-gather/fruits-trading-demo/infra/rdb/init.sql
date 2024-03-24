CREATE TABLE IF NOT EXISTS results (
    id serial PRIMARY KEY,
    feature_type text NOT NULL,
    value float NOT NULL
);
