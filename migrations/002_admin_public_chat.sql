BEGIN;

ALTER TABLE admins
	ADD COLUMN IF NOT EXISTS locationid INT REFERENCES locations(id);

UPDATE admins AS a
SET locationid = l.id
FROM locations AS l
WHERE a.locationid IS NULL
	AND replace(lower(a.username), '-', '') = replace(lower(l.slug), '-', '');

UPDATE admins
SET locationid = CASE id
	WHEN 'adm-kpl' THEN (SELECT id FROM locations WHERE slug = 'kopiloka')
	WHEN 'adm-kkt' THEN (SELECT id FROM locations WHERE slug = 'koktong')
	WHEN 'adm-kbg' THEN (SELECT id FROM locations WHERE slug = 'kopi-braga')
	ELSE locationid
END
WHERE locationid IS NULL;

DO $$
BEGIN
	IF EXISTS (SELECT 1 FROM admins WHERE locationid IS NULL) THEN
		RAISE EXCEPTION 'Every admin must be assigned to a location before this migration can finish';
	END IF;
END
$$;

ALTER TABLE admins
	ALTER COLUMN locationid SET NOT NULL;

ALTER TABLE publicmessages
	ADD COLUMN IF NOT EXISTS adminid TEXT REFERENCES admins(id);

ALTER TABLE publicmessages
	DROP CONSTRAINT IF EXISTS publicmessages_sender_check;

ALTER TABLE publicmessages
	ADD CONSTRAINT publicmessages_sender_check
	CHECK ((senderid IS NOT NULL)::int + (adminid IS NOT NULL)::int = 1);

CREATE INDEX IF NOT EXISTS idx_publicmessages_location_timestamp
	ON publicmessages(locationid, "timestamp" DESC);

COMMIT;
