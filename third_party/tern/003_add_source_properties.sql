ALTER TABLE nrtm_source
	ADD COLUMN properties jsonb NOT NULL default '{}'::jsonb
	;

-----------------------------------
---- create above / drop below ----
-----------------------------------

	ALTER TABLE nrtm_source
	DROP COLUMN properties
	;
