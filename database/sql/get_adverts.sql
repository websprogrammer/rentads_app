SELECT
    post_id,
    description,
    publication_date AS date,
    profile_link,
    profile_name,
    photos,
    city,
    district,
    sub_district,
    metro,
    rent_type,
    room_type,
    price
FROM adverts
WHERE
    (@older_than = 0 OR publication_date < @older_than)
    AND city=@city
    AND (@rent_type = 0 OR rent_type = @rent_type)
    AND (@room_type = 0 OR room_type = @room_type)
    AND (cardinality(COALESCE(@districts::text[], '{}')) = 0 OR district = ANY(@districts::text[]))
    AND (cardinality(COALESCE(@sub_districts::text[], '{}')) = 0 OR sub_district = ANY(@sub_districts::text[]))
    AND (cardinality(COALESCE(@metro_stations::text[], '{}')) = 0 OR metro = ANY(@metro_stations::text[]))
    AND (@key_words = '' OR description_tokens @@ to_tsquery('russian', @key_words))
ORDER BY publication_date DESC
LIMIT @page_size;