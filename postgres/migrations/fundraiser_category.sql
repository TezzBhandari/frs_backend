CREATE TABLE IF NOT EXISTS fundraiser_category (
    id BIGINT  PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(150) NOT NULL
);