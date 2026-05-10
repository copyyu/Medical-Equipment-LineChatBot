-- Enable pg_trgm extension for trigram similarity search (used by FindBestMatch / FindSimilarSorted)
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- GIN index on equipments.id_code for fast trigram similarity lookups
CREATE INDEX IF NOT EXISTS idx_equipments_id_code_trgm ON equipments USING GIN (id_code gin_trgm_ops);
