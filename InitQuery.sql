

-- Users table (optional, just in case you want to reference)
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

-- Charts
CREATE TABLE IF NOT EXISTS charts (
    id VARCHAR(50) PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    x_axis_title VARCHAR(100),
    y_axis_title VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS chart_data (
    id SERIAL PRIMARY KEY,
    chart_id VARCHAR(50) REFERENCES charts(id) ON DELETE CASCADE,
    datapoint_code VARCHAR(50) NOT NULL,
    value DOUBLE PRECISION NOT NULL
);

-- Insights
CREATE TABLE IF NOT EXISTS insights (
    id VARCHAR(50) PRIMARY KEY,
    description TEXT
);

-- Audiences
CREATE TABLE IF NOT EXISTS audiences (
    id VARCHAR(50) PRIMARY KEY,
    gender VARCHAR(10),
    country VARCHAR(50),
    age_group VARCHAR(20),
    social_hours INT,
    purchases INT,
    description TEXT
);

-- Favourites linking user and assets (asset_type to identify which table)
CREATE TABLE IF NOT EXISTS favourites (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    asset_id VARCHAR(50) NOT NULL,
    asset_type VARCHAR(20) NOT NULL, -- 'chart', 'insight', 'audience'
    UNIQUE(user_id, asset_id, asset_type)
);








