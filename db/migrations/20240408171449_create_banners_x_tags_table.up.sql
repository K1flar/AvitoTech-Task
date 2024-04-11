CREATE TABLE IF NOT EXISTS banner_x_tag(
    banner_id INTEGER REFERENCES banners(id) ON DELETE CASCADE NOT NULL,
    tag_id INTEGER REFERENCES tags(id) ON DELETE CASCADE NOT NULL,
    feature_id INTEGER REFERENCES features(id) ON DELETE CASCADE NOT NULL,
    UNIQUE (tag_id, feature_id)
);