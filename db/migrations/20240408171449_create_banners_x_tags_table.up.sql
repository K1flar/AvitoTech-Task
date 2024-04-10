CREATE TABLE IF NOT EXISTS banners_x_tags(
    tag_id INTEGER REFERENCES tags(id) ON DELETE CASCADE,
    banner_id INTEGER,
    banner_updated_dttm TIMESTAMP NOT NULL,
    PRIMARY KEY (tag_id, banner_id, banner_updated_dttm),
    FOREIGN KEY (banner_id, banner_updated_dttm) REFERENCES banners(id, updated_dttm) ON DELETE CASCADE  
);