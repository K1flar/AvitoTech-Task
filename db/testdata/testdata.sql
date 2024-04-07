SELECT setseed(0.5); 

WITH random_content AS (
    SELECT 
    ('{
        "title": "' || md5(random()::text) || '", 
        "url": "http://example.com/' || md5(random()::text) || '", 
        "image": "http://example.com/image/' || md5(random()::text) || '"
    }')::JSONB AS content,
    (CURRENT_TIMESTAMP-'7 year'::interval + (random()*('7 year'::interval))::interval) AS created_dttm,
    random() > 0.5 AS is_active
    FROM generate_series(1, 1000) 
)

INSERT INTO banners(content, created_dttm, is_active)
SELECT content, created_dttm, is_active FROM random_content;

INSERT INTO tags(id)
VALUES (generate_series(1,1000));

INSERT INTO features(id)
VALUES (generate_series(1,1000));

INSERT INTO banners_info(feature_id, tag_id, banner_id)
SELECT
    sub.feature_id AS feature_id,
    unnest(sub.tag_ids) AS tag_id,
    sub.banner_id AS banner_id
FROM (
    SELECT
        feature_id,
        banner_id,
        array_agg(tag_id) AS tag_ids
    FROM
        (SELECT
            generate_series(1, 1000) AS feature_id,
            generate_series(1, 1000) AS banner_id,
            trunc(random() * 10)::int + 1 AS num_tags
         ) AS generated_data
    CROSS JOIN LATERAL (
        SELECT (1 + random()*999)::integer AS tag_id FROM generate_series(1, num_tags)
    ) AS tag_data
    GROUP BY feature_id, banner_id
) AS sub
ON CONFLICT (feature_id, tag_id) DO NOTHING;

INSERT INTO banner_versions(banner_id, content, updated_dttm)
SELECT
    sub.banner_id AS banner_id,
    unnest(sub.contents) AS content,
	unnest(sub.update_at) AS update_at
FROM (
    SELECT
        banner_id,
        array_agg(content) AS contents,
		array_agg(update_at) AS update_at
    FROM
        (SELECT
            (1+random()*999)::integer AS banner_id,
            (1+ random() * 2)::int AS num_versions
		 FROM generate_series(1, 200)
         ) AS generated_data
    CROSS JOIN LATERAL (
        SELECT ('{
			"title": "' || md5(random()::text) || '", 
			"url": "http://example.com/' || md5(random()::text) || '", 
			"image": "http://example.com/image/' || md5(random()::text) || '"
		}')::JSONB AS content, 
		(CURRENT_TIMESTAMP-'1 year'::interval + (random()*('1 year'::interval))::interval) AS update_at
		FROM generate_series(1, num_versions)
    )
    GROUP BY banner_id
) AS sub;