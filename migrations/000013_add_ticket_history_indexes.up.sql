-- Index the ticket_histories columns used by the admin activity-log path.
-- GetStatusChangeLogs filters on action, orders by created_at, and the stats
-- query groups by new_value — all sequential-scanned before this.
CREATE INDEX IF NOT EXISTS idx_ticket_histories_action     ON ticket_histories (action);
CREATE INDEX IF NOT EXISTS idx_ticket_histories_created_at ON ticket_histories (created_at);
