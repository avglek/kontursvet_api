-- Leads table
CREATE TABLE IF NOT EXISTS leads (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone_digital VARCHAR(20) NOT NULL,
    phone_format VARCHAR(30) NOT NULL,
    home_type VARCHAR(100) NOT NULL,
    location VARCHAR(255),
    message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Portfolio cases
CREATE TABLE IF NOT EXISTS portfolio_cases (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    case_item VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    task TEXT,
    location VARCHAR(255),
    term VARCHAR(100),
    team VARCHAR(255),
    period VARCHAR(100),
    features TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Portfolio case works (many-to-one)
CREATE TABLE IF NOT EXISTS portfolio_case_works (
    id BIGSERIAL PRIMARY KEY,
    case_id BIGINT NOT NULL REFERENCES portfolio_cases(id) ON DELETE CASCADE,
    work TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Portfolio case meta (many-to-one)
CREATE TABLE IF NOT EXISTS portfolio_case_meta (
    id BIGSERIAL PRIMARY KEY,
    case_id BIGINT NOT NULL REFERENCES portfolio_cases(id) ON DELETE CASCADE,
    meta_item VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Portfolio photos
CREATE TABLE IF NOT EXISTS portfolio_photos (
    id BIGSERIAL PRIMARY KEY,
    case_id BIGINT NOT NULL REFERENCES portfolio_cases(id) ON DELETE CASCADE,
    gallery_key INT NOT NULL,
    photo_path VARCHAR(500) NOT NULL,
    alt VARCHAR(255),
    caption VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_leads_created_at ON leads(created_at);
CREATE INDEX IF NOT EXISTS idx_portfolio_cases_name ON portfolio_cases(name);
CREATE INDEX IF NOT EXISTS idx_portfolio_photos_case_id ON portfolio_photos(case_id);
CREATE INDEX IF NOT EXISTS idx_portfolio_case_works_case_id ON portfolio_case_works(case_id);
CREATE INDEX IF NOT EXISTS idx_portfolio_case_meta_case_id ON portfolio_case_meta(case_id);
