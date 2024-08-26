INSERT INTO users (username, password_hash) VALUES
('testuser1', '$2a$10$ZmpzL3j.fEOgsfno.MHzNuYdkfQr5PRoUTWUbkJHhVvF6HMcwcwSW'), --password is "password"
('testuser2', '$2a$10$ZmpzL3j.fEOgsfno.MHzNuYdkfQr5PRoUTWUbkJHhVvF6HMcwcwSW'),
('testuser3', '$2a$10$ZmpzL3j.fEOgsfno.MHzNuYdkfQr5PRoUTWUbkJHhVvF6HMcwcwSW');

INSERT INTO user_notes (user_id, note_text) VALUES
(1, 'Test note 1 for user 1'),
(1, 'Test note 2 for user 1'),
(2, 'Test note 1 for user 2'),
(3, 'Test note 1 for user 3'),
(3, 'Test note 2 for user 3'),
(3, 'Test note 3 for user 3');
