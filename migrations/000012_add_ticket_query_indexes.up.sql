-- Indexes for the ticket query paths that currently sequential-scan.

-- reporter_line_id backs the LINE bot's "my tickets" lookup
-- (GetTicketsByLineUserID / FindPendingTicketByEquipmentAndUser), hit by every
-- end user. Only reporter_id (the admin FK) was indexed before.
CREATE INDEX IF NOT EXISTS idx_tickets_reporter_line_id ON tickets (reporter_line_id);

-- status / priority back the admin dashboard list filters (GetAllTickets) and
-- the GROUP BY status stats (GetTicketStats).
CREATE INDEX IF NOT EXISTS idx_tickets_status   ON tickets (status);
CREATE INDEX IF NOT EXISTS idx_tickets_priority ON tickets (priority);
