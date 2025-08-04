-- Create application_metrics table
CREATE TABLE IF NOT EXISTS application_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    cpu_usage DECIMAL(5,2) NOT NULL CHECK (cpu_usage >= 0),
    memory_usage DECIMAL(10,2) NOT NULL CHECK (memory_usage >= 0),
    request_count BIGINT NOT NULL DEFAULT 0 CHECK (request_count >= 0),
    error_count BIGINT NOT NULL DEFAULT 0 CHECK (error_count >= 0),
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_application_metrics_app_id ON application_metrics(application_id);
CREATE INDEX IF NOT EXISTS idx_application_metrics_timestamp ON application_metrics(timestamp);

-- Create composite index for efficient metrics queries
CREATE INDEX IF NOT EXISTS idx_application_metrics_app_time ON application_metrics(application_id, timestamp DESC);

-- Create index for time-based queries
CREATE INDEX IF NOT EXISTS idx_application_metrics_time_range ON application_metrics(timestamp DESC); 