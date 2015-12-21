CREATE TABLE hosts
(
  id serial NOT NULL,
  host character varying NOT NULL,
  thumb_url character varying NULL,
  created_at timestamp without time zone NOT NULL,
  updated_at timestamp without time zone NOT NULL
);

ALTER TABLE "hosts"
  ADD CONSTRAINT "hosts_id" PRIMARY KEY ("id"),
  ADD CONSTRAINT "hosts_host" UNIQUE ("host");

CREATE TABLE links
(
  id serial NOT NULL,
  host_id integer,
  link character varying,
  created_at timestamp without time zone NOT NULL,
  updated_at timestamp without time zone NOT NULL
);

ALTER TABLE "links"
  ADD CONSTRAINT "links_id" PRIMARY KEY ("id");

CREATE INDEX index_links_on_host_id ON links USING btree (host_id);

ALTER TABLE links
    ADD CONSTRAINT links_hosts FOREIGN KEY (host_id) REFERENCES hosts(id);

