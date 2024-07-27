CREATE TABLE IF NOT EXISTS fundraisers (
    id BIGINT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    story TEXT NOT NULL,
    target_amount DECIMAL(10, 2) NOT NULL,
    cover_img VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);