BEGIN;

ALTER TABLE images
	ADD COLUMN release_id INTEGER NOT NULL;

ALTER TABLE images
	ADD CONSTRAINT images_release_id_fkey
	FOREIGN KEY(release_id) REFERENCES releases(id)
	ON DELETE CASCADE;

COMMIT;
