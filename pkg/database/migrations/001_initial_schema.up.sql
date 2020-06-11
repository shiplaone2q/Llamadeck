-- -----------------------------------------------------------------------------
-- Instances Table: Central configuration and state for LLM backends
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS instances (
    -- Primary identification
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,

    -- Instance state
    status TEXT NOT NULL CHECK(status IN ('stopped', 'running', 'failed', 'restarting', 'shutting_down')) DEFAULT 'stopped',

    -- Timestamps (created_at stored as Unix timestamp for compatibility with existing JSON format)
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,

    -- All instance options stored as a single JSON blob
