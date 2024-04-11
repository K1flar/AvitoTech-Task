SELECT setseed(0.5); 

INSERT INTO tags(id)
VALUES (generate_series(1,1000));

INSERT INTO features(id)
VALUES (generate_series(1,1000));

WITH random_content AS (
    SELECT 
    ('{
        "title": "' || md5(random()::text) || '", 
        "url": "http://example.com/' || md5(random()::text) || '", 
        "image": "http://example.com/image/' || md5(random()::text) || '"
    }')::JSONB AS content,
    random_date AS created_dttm,
    random_date AS updated_dttm,
    random() > 0.5 AS is_active,
    (1 + random()*999)::integer AS feature_id
    FROM ( 
        SELECT (CURRENT_TIMESTAMP-'7 year'::interval + (random()*('6 year'::interval))::interval) as random_date
        FROM generate_series(1,1000)
    )
)

-- INSERT banners
INSERT INTO banners(content, created_dttm, updated_dttm, is_active, feature_id)
SELECT content, created_dttm, updated_dttm, is_active, feature_id FROM random_content;

WITH generated_data AS (
    SELECT * FROM (SELECT
        id as banner_id,
        feature_id,
        (3 + random()*7)::integer AS num_tags
        FROM banners
    ) 
    CROSS JOIN LATERAL (
        SELECT (1 + random()*999)::integer AS tag_id FROM generate_series(1, num_tags)
    )
)

INSERT INTO banner_x_tag(banner_id, tag_id, feature_id)
SELECT
    sub.banner_id AS banner_id,
    unnest(sub.tag_ids) AS tag_id,
    feature_id
FROM (
    SELECT
        banner_id,
        array_agg(tag_id) AS tag_ids,
        feature_id
    FROM generated_data AS gd
    WHERE 1=(SELECT COUNT(feature_id) FROM generated_data WHERE tag_id=gd.tag_id AND feature_id=gd.feature_id)
    GROUP BY banner_id, feature_id
	ORDER BY banner_id
) AS sub
ON CONFLICT (feature_id, tag_id) DO NOTHING;