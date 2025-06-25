INSERT INTO notifications
    (
        device_token,
        city,
        key_words,
        rent_type,
        room_type,
        districts,
        enabled,
        updated_date
    )
VALUES
    (
        @device_token,
        @city,
        @key_words,
        @rent_type,
        @room_type,
        @districts,
        @enabled,
        EXTRACT(EPOCH FROM CURRENT_TIMESTAMP)::int
    )
ON CONFLICT (device_token) DO UPDATE
    SET city=@city,
        key_words=@key_words,
        rent_type=@rent_type,
        room_type=@room_type,
        districts=@districts,
        enabled=@enabled,
        updated_date=EXTRACT(EPOCH FROM CURRENT_TIMESTAMP)::int