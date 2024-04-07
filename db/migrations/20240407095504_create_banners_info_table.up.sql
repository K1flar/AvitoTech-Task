CREATE TABLE IF NOT EXISTS banners_info(
    feature_id INTEGER REFERENCES features(id),
    tag_id INTEGER REFERENCES tags(id),
    banner_id INTEGER REFERENCES banners(id) ON DELETE CASCADE NOT NULL,
    PRIMARY KEY (feature_id, tag_id)
);