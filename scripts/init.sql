-- Create software_versions table
CREATE TABLE IF NOT EXISTS software_versions (
    id VARCHAR(255) PRIMARY KEY,
    app_id VARCHAR(255) NOT NULL,
    version VARCHAR(255) NOT NULL,
    filename VARCHAR(255),
    checksum VARCHAR(255),
    size BIGINT,
    created_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT,
    UNIQUE(app_id, version)
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_software_versions_app_id ON software_versions(app_id);
CREATE INDEX IF NOT EXISTS idx_software_versions_version ON software_versions(version);
CREATE INDEX IF NOT EXISTS idx_software_versions_created_at ON software_versions(created_at DESC);

-- Grant permissions
GRANT ALL PRIVILEGES ON software_versions TO postgres;
