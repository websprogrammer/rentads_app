INSERT INTO feedbacks
    (
        city,
        post_id,
        feedback_type,
        message
    )
VALUES
    (
        @city,
        @post_id,
        @feedback_type,
        @message
    )
ON CONFLICT (city, post_id, feedback_type) DO NOTHING;
