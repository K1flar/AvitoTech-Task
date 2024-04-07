CREATE TABLE IF NOT EXISTS banner_versions(
    id SERIAL PRIMARY KEY,
    banner_id INTEGER REFERENCES banners(id) ON DELETE CASCADE,
    content JSON,
    updated_at TIMESTAMP 
);