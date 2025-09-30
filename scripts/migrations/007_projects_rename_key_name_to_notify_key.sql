-- Rename projects.key_name to notify_key and related indexes/constraints
BEGIN;

-- 1) Rename column
ALTER TABLE projects RENAME COLUMN key_name TO notify_key;

-- 2) Recreate unique constraint if named (depends on prior name). Drop if exists then add new.
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'projects_company_id_key_name_key'
    ) THEN
        ALTER TABLE projects DROP CONSTRAINT projects_company_id_key_name_key;
    END IF;
END$$;

ALTER TABLE projects ADD CONSTRAINT projects_company_id_notify_key_key UNIQUE (company_id, notify_key);

-- 3) Update any dependent views or materialized views if necessary (none here)

COMMIT;


