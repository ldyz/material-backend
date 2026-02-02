-- Create stock_logs table
CREATE TABLE IF NOT EXISTS stock_logs (
    id BIGSERIAL PRIMARY KEY,
    stock_id BIGINT NOT NULL REFERENCES stocks(id) ON DELETE CASCADE,
    type VARCHAR(10) NOT NULL CHECK (type IN ('in', 'out')),
    quantity DECIMAL(15,3) NOT NULL,
    time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    remark VARCHAR(100),
    project_id BIGINT,
    user_id BIGINT,
    recipient VARCHAR(255),
    purpose VARCHAR(255),
    requisition_id BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_stock_logs_stock_id ON stock_logs(stock_id);
CREATE INDEX IF NOT EXISTS idx_stock_logs_project_id ON stock_logs(project_id);
CREATE INDEX IF NOT EXISTS idx_stock_logs_requisition_id ON stock_logs(requisition_id);

-- Create stock_op_logs table
CREATE TABLE IF NOT EXISTS stock_op_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    op_type VARCHAR(20),
    stock_id BIGINT REFERENCES stocks(id) ON DELETE CASCADE,
    log_id BIGINT,
    detail TEXT,
    time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_stock_op_logs_stock_id ON stock_op_logs(stock_id);
