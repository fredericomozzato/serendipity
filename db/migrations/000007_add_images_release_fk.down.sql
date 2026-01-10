BEGIN;

ALTER TABLE images
	DROP CONSTRAINT images_release_id_fkey;

ALTER TABLE images
	DROP COLUMN release_id;

COMMIT;
