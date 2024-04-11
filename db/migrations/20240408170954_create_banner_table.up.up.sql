CREATE TABLE IF NOT EXISTS banners(
    id SERIAL PRIMARY KEY,
    content JSONB,
    created_dttm TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_dttm TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_active BOOL NOT NULL DEFAULT true,
    feature_id INTEGER REFERENCES features(id) ON DELETE CASCADE NOT NULL
);