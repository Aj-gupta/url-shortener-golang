DO
$do$
BEGIN
   IF EXISTS (SELECT FROM pg_database WHERE datname = 'postgres') THEN
      RAISE NOTICE 'Database already exists';  -- optional
   ELSE
      PERFORM dblink_exec('dbname=' || current_database()  -- current db
                        , 'CREATE DATABASE postgres');
   END IF;
END
$do$;
CREATE TABLE public.short_urls (
  id bigserial NOT NULL, 
  created_at timestamp without time zone NOT NULL DEFAULT now(), 
  updated_at character varying(255) NOT NULL DEFAULT now(), 
  short_url character varying(255) NOT NULL, 
  long_url character varying(255) NOT NULL
);

ALTER TABLE 
  public.short_urls 
ADD 
  CONSTRAINT short_urls_pkey PRIMARY KEY (id);
CREATE UNIQUE INDEX short_url_idx ON short_urls (short_url);
CREATE UNIQUE INDEX long_url_idx ON short_urls (long_url);

CREATE 
OR REPLACE FUNCTION base62_encode(IN digits bigint) RETURNS varchar AS $$DECLARE chars char [];
reslt varchar;
val bigint;
BEGIN chars := ARRAY [ '0', 
'1', 
'2', 
'3', 
'4', 
'5', 
'6', 
'7', 
'8', 
'9', 
'a', 
'b', 
'c', 
'd', 
'e', 
'f', 
'g', 
'h', 
'i', 
'j', 
'k', 
'l', 
'm', 
'n', 
'o', 
'p', 
'q', 
'r', 
's', 
't', 
'u', 
'v', 
'w', 
'x', 
'y', 
'z', 
'A', 
'B', 
'C', 
'D', 
'E', 
'F', 
'G', 
'H', 
'I', 
'J', 
'K', 
'L', 
'M', 
'N', 
'O', 
'P', 
'Q', 
'R', 
'S', 
'T', 
'U', 
'V', 
'W', 
'X', 
'Y', 
'Z' ];
val := digits;
reslt := '';
IF val < 0 THEN val := val * -1;
END IF;
WHILE val != 0 LOOP reslt := chars [(val % 62) + 1] || reslt;
val := val / 62;
END LOOP;
RETURN reslt;
END;
$$LANGUAGE plpgsql IMMUTABLE;

CREATE 
OR REPLACE FUNCTION trigger_set_timestamp_after_update() RETURNS TRIGGER AS $$BEGIN New.updated_at = now();
RETURN NEW;
END;
$$LANGUAGE plpgsql;

CREATE 
OR REPLACE FUNCTION trigger_short_urls_create() RETURNS TRIGGER AS $$ BEGIN NEW.short_url = base62_encode(NEW.id);
NEW.created_at = now();
NEW.updated_at = NEW.created_at;
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER short_urls_create BEFORE INSERT ON short_urls FOR EACH ROW EXECUTE PROCEDURE trigger_short_urls_create();

CREATE TRIGGER short_urls_update BEFORE 
UPDATE 
  ON short_urls FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp_after_update();