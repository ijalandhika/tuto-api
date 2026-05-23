-- migrate:down

DROP TABLE IF EXISTS parent_gate_attempts;
DROP TABLE IF EXISTS auth_sessions;
DROP TABLE IF EXISTS parents;
DROP TYPE  IF EXISTS actor_type;