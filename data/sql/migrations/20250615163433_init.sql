-- migrate:up
-- Users table for CMS authentication
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    role VARCHAR(50) NOT NULL DEFAULT 'editor', -- admin, editor, viewer
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Categories for organizing programs
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    slug VARCHAR(100) UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- programs table
CREATE TABLE programs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    summary TEXT, -- Short description for listings
    category_id UUID REFERENCES categories (id),
    language VARCHAR(10) DEFAULT 'ar', -- ISO language code
    country VARCHAR(10) DEFAULT 'SA', -- ISO country code
    author VARCHAR(255),
    publisher VARCHAR(255),
    artwork_url TEXT,
    website_url TEXT,
    rss_feed_url TEXT,
    is_explicit BOOLEAN DEFAULT false,
    status VARCHAR(20) DEFAULT 'active', -- active, inactive, draft
    total_episodes INTEGER DEFAULT 0,
    average_duration INTEGER, -- in seconds
    rating DECIMAL(3, 2), -- average rating
    total_ratings INTEGER DEFAULT 0,
    source VARCHAR(50) NOT NULL DEFAULT 'local', -- local, itunes, spotify, etc.
    created_by UUID REFERENCES users (id),
    created_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        published_at TIMESTAMP
    WITH
        TIME ZONE
);

-- Episodes table for individual podcast episodes
CREATE TABLE episodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    program_id UUID NOT NULL REFERENCES programs (id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    summary TEXT,
    episode_number INTEGER,
    season_number INTEGER,
    duration INTEGER, -- in seconds
    audio_url TEXT,
    file_size BIGINT, -- in bytes
    mime_type VARCHAR(50),
    artwork_url TEXT,
    is_explicit BOOLEAN DEFAULT false,
    status VARCHAR(20) DEFAULT 'published', -- published, draft, scheduled
    play_count INTEGER DEFAULT 0,
    rating DECIMAL(3, 2),
    total_ratings INTEGER DEFAULT 0,
    published_at TIMESTAMP
    WITH
        TIME ZONE,
        created_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Enhanced external sources tracking
CREATE TABLE external_sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    program_id UUID REFERENCES programs (id) ON DELETE CASCADE,
    episode_id UUID REFERENCES episodes (id) ON DELETE CASCADE,
    source_name VARCHAR(50) NOT NULL, -- itunes, spotify, google_podcasts
    external_id VARCHAR(255) NOT NULL,
    external_url TEXT,
    raw_data JSONB, -- Store original API response
    last_synced_at TIMESTAMP
    WITH
        TIME ZONE,
        created_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        UNIQUE (source_name, external_id)
);

-- Tags for flexible categorization
CREATE TABLE tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    name VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- migrate:down
DROP TABLE IF EXISTS external_sources;

DROP TABLE IF EXISTS episodes;

DROP TABLE IF EXISTS programs;

DROP TABLE IF EXISTS categories;

DROP TABLE IF EXISTS tags;

DROP TABLE IF EXISTS users;
