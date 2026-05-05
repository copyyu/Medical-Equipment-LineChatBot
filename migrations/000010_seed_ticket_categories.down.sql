-- Remove seeded ticket categories (only the default ones)
DELETE FROM ticket_categories
WHERE name IN ('แจ้งซ่อม', 'บำรุงรักษา', 'สอบถามการใช้งาน', 'อื่นๆ');
