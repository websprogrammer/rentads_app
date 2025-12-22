SELECT
    id,
    post_id,
    city,
    feedback_type,
    message
FROM feedbacks
WHERE
    (@id = 0 OR id < @id)
    AND city=@city
ORDER BY id DESC
LIMIT @page_size;