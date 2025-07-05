CREATE TABLE project_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    priority VARCHAR(20) NOT NULL DEFAULT 'medium',
    estimated_hours DECIMAL(8,2),
    actual_hours DECIMAL(8,2),
    due_date DATE,
    assigned_to UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_project_items_project_id ON project_items(project_id);
CREATE INDEX idx_project_items_status ON project_items(status);
CREATE INDEX idx_project_items_priority ON project_items(priority);
CREATE INDEX idx_project_items_assigned_to ON project_items(assigned_to);
CREATE INDEX idx_project_items_deleted_at ON project_items(deleted_at);
CREATE INDEX idx_project_items_created_at ON project_items(created_at); 