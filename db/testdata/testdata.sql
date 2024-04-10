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

-- INSERT versions
INSERT INTO banners(id, content, created_dttm, updated_dttm, is_active, feature_id)
SELECT id, content, created_dttm, updated_dttm, is_active, feature_id FROM (
    SELECT 
        id AS id,
        ('{"version":"new version"}')::JSONB AS content,
        created_dttm,
        updated_dttm,
        random() > 0.5 AS is_active,
        (1 + random()*999)::integer AS feature_id
    FROM (
        SELECT 
            (1 + random()*999)::integer AS id, 
            (CURRENT_TIMESTAMP-'7 year'::interval + (random()*('7 year'::interval))::interval) AS created_dttm,
            (CURRENT_TIMESTAMP-'1 year'::interval + (random()*('1 year'::interval))::interval) AS updated_dttm
        FROM generate_series(1,1000) 
    )
);

INSERT INTO banners_x_tags(tag_id, banner_id, banner_updated_dttm)
SELECT
    unnest(sub.tag_ids) AS tag_id,
    sub.banner_id AS banner_id,
	sub.banner_updated_dttm AS banner_updated_dttm 
FROM (
    SELECT
        array_agg(tag_id) AS tag_ids,
        banner_id,
		banner_updated_dttm
    FROM
        (SELECT
            id as banner_id,
            updated_dttm as banner_updated_dttm,
            trunc(random() * 10)::int + 1 AS num_tags
		 FROM banners
         ) AS generated_data
    CROSS JOIN LATERAL (
        SELECT (1 + random()*999)::integer AS tag_id FROM generate_series(1, num_tags)
    ) AS tag_data
    GROUP BY banner_id, banner_updated_dttm
	ORDER BY banner_id
) AS sub
ON CONFLICT (tag_id, banner_id, banner_updated_dttm) DO NOTHING;
