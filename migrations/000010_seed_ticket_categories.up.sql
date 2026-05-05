-- Seed default ticket categories
INSERT INTO ticket_categories (name, name_en, color, icon, sort_order, is_active)
VALUES
    ('แจ้งซ่อม',          'Repair',           '#EF5350', '🔧', 1, TRUE),
    ('บำรุงรักษา',        'Maintenance',      '#FFA726', '🛠️', 2, TRUE),
    ('สอบถามการใช้งาน',   'Usage Inquiry',    '#42A5F5', '❓', 3, TRUE),
    ('อื่นๆ',             'Other',            '#78909C', '📝', 4, TRUE)
ON CONFLICT DO NOTHING;
