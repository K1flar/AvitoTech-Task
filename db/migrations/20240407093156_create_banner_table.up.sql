CREATE TABLE IF NOT EXISTS banners(
    id SERIAL PRIMARY KEY,
    content JSON,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    is_activ BOOL NOT NULL DEFAULT true
);