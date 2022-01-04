-- inserts some example data into the database

-- create 2 categories
INSERT INTO categories (id, title, user_id)
VALUES
	(0, 'Skool', 0),
	(1, 'CCA', 0);

-- create 5 tasks
INSERT INTO tasks (id, category_id, title, description, deadline, completed, user_id)
VALUES
	(0, 0, 'Do Lab 3', 'prolly would need 3 hours (ah who am I kidding make that 9).', NULL, FALSE, 0),
	(1, 0, 'Revise for Midterms', 'gg bellcurve-god save me.', NULL, FALSE, 0),
	(2, 1, 'Prepare for CCA meeting on Friday', 'best not to show up empty-handed.', CURRENT_TIMESTAMP + INTERVAL '5 days', FALSE, 0);



