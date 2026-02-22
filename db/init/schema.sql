-- PostgreSQL schema for makes, models, trims, and manuals
-- Derived from:
--   Table makes { id integer [pk]; name varchar; year varchar }
--   Table models { id integer [pk]; name varchar; makeId integer [Ref: > makes.id] }
--   Table trims  { id integer [pk]; name varchar; modelId integer [Ref: > models.id]; manualId integer [Ref: > manuals.id] }
--   Table manuals{ id integer [pk]; url varchar }
--
-- Recommendations:
-- - Use BIGSERIAL/BIGINT for IDs (safer for Go int64 usage).
-- - Use NOT NULL where appropriate.
-- - Add foreign keys with sensible ON DELETE behavior.
-- - Add indexes for lookup performance.

BEGIN;

-- Manuals: simple store for manual URLs
CREATE TABLE IF NOT EXISTS manual (
  id       BIGSERIAL PRIMARY KEY,
  url      TEXT NOT NULL
);

-- Makes: manufacturer names, optional year field (kept as text to allow ranges or descriptors)
CREATE TABLE IF NOT EXISTS make (
  id       BIGSERIAL PRIMARY KEY,
  name     TEXT NOT NULL,
  year     VARCHAR(20)
);

-- Models: belong to a make
CREATE TABLE IF NOT EXISTS model (
  id        BIGSERIAL PRIMARY KEY,
  name      TEXT NOT NULL,
  make_id   BIGINT NOT NULL REFERENCES make(id) ON DELETE CASCADE
);

-- Trims: belong to a model and optionally reference a manual
CREATE TABLE IF NOT EXISTS trim (
  id         BIGSERIAL PRIMARY KEY,
  name       TEXT NOT NULL,
  model_id   BIGINT NOT NULL REFERENCES model(id) ON DELETE CASCADE,
  manual_id  BIGINT REFERENCES manual(id) ON DELETE SET NULL
);

-- Useful indexes
CREATE INDEX IF NOT EXISTS idx_model_make_id ON model(make_id);
CREATE INDEX IF NOT EXISTS idx_trim_model_id ON trim(model_id);
CREATE INDEX IF NOT EXISTS idx_trim_manual_id ON trim(manual_id);

-- Optional uniqueness constraints (uncomment if desired)
-- Ensure make names are unique (across case-sensitive values)
-- CREATE UNIQUE INDEX IF NOT EXISTS ux_makes_name ON makes(LOWER(name));
-- Ensure model names are unique per make
-- CREATE UNIQUE INDEX IF NOT EXISTS ux_models_make_name ON models(make_id, LOWER(name));
-- Ensure trim names are unique per model
-- CREATE UNIQUE INDEX IF NOT EXISTS ux_trims_model_name ON trims(model_id, LOWER(name));

COMMIT;
