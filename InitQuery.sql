-- Users table
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(50) PRIMARY KEY,   -- TODO store UUID 
    name VARCHAR(100) NOT NULL
);

-- Assets table
CREATE TABLE IF NOT EXISTS assets (
    asset_id VARCHAR(50) PRIMARY KEY,  
    title TEXT NOT NULL,
    description TEXT,
    asset_type VARCHAR(20) NOT NULL,   -- 'chart', 'insight', 'audience'
    user_id VARCHAR(50) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Charts
CREATE TABLE IF NOT EXISTS charts (
    id VARCHAR(50) PRIMARY KEY,       
    title VARCHAR(200) NOT NULL,
    description TEXT,
    x_axis_title VARCHAR(100),
    y_axis_title VARCHAR(100),
    FOREIGN KEY (id) REFERENCES assets(asset_id) ON DELETE CASCADE
);

-- Chart Data Points
CREATE TABLE IF NOT EXISTS chart_data (
    id SERIAL PRIMARY KEY,
    chart_id VARCHAR(50) NOT NULL,
    datapoint_code VARCHAR(50) NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    FOREIGN KEY (chart_id) REFERENCES charts(id) ON DELETE CASCADE
);

-- Insights
CREATE TABLE IF NOT EXISTS insights (
    id VARCHAR(50) PRIMARY KEY,       
    description TEXT,
    FOREIGN KEY (id) REFERENCES assets(asset_id) ON DELETE CASCADE
);

-- Audiences
CREATE TABLE IF NOT EXISTS audiences (
    id VARCHAR(50) PRIMARY KEY,       
    gender VARCHAR(10),
    country VARCHAR(50),
    age_group VARCHAR(20),
    social_hours INT,
    purchases INT,
    description TEXT,
    FOREIGN KEY (id) REFERENCES assets(asset_id) ON DELETE CASCADE
);

-- Favourites (link table)
CREATE TABLE IF NOT EXISTS favourites (
    user_id VARCHAR(50) NOT NULL,
    asset_id VARCHAR(50) NOT NULL,
    asset_type VARCHAR(20) NOT NULL,
    PRIMARY KEY (user_id, asset_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (asset_id) REFERENCES assets(asset_id) ON DELETE CASCADE
);
