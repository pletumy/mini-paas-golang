-- Create applications table
CREATE TABLE IF NOT EXISTS applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    git_url VARCHAR(500) NOT NULL,
    branch VARCHAR(100) NOT NULL DEFAULT 'main',
    port INTEGER NOT NULL CHECK (port >= 1 AND port <= 65535),
    environment VARCHAR(20) NOT NULL DEFAULT 'development' CHECK (environment IN ('development', 'staging', 'production')),
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'building', 'deploying', 'running', 'failed', 'stopped')),
    image_url VARCHAR(500),
    deployment_url VARCHAR(500),
    replicas INTEGER NOT NULL DEFAULT 1 CHECK (replicas >= 1 AND replicas <= 10),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_applications_user_id ON applications(user_id);
CREATE INDEX IF NOT EXISTS idx_applications_status ON applications(status);
CREATE INDEX IF NOT EXISTS idx_applications_environment ON applications(environment);
CREATE INDEX IF NOT EXISTS idx_applications_created_at ON applications(created_at);

-- Create unique constraint on name per user
CREATE UNIQUE INDEX IF NOT EXISTS idx_applications_user_name ON applications(user_id, name);

-- Create trigger to update updated_at timestamp
CREATE TRIGGER update_applications_updated_at 
    BEFORE UPDATE ON applications 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column(); 