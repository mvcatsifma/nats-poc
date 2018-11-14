CREATE FOREIGN TABLE product_stream (
  uuid         varchar(255),
  machine_id   varchar(255),
  product_type varchar(255),
  ok           boolean,
  created_at   timestamptz
)
SERVER pipelinedb;

CREATE VIEW products_total WITH (action = materialize) AS
  SELECT count(*)
  FROM product_stream;

CREATE VIEW products_ok WITH (action = materialize) AS
  SELECT ok, count(*)
  FROM product_stream
  GROUP BY ok;

CREATE VIEW products_machine WITH (action = materialize) AS
  SELECT machine_id, count(*)
  FROM product_stream
  GROUP BY machine_id;

CREATE VIEW products_product_type WITH (action = materialize) AS
  SELECT product_type, count(*)
  FROM product_stream
  GROUP BY product_type;

--
select *
from products_ok;
--
select *
from products_machine;
--
select *
from products_product_type;
--
select *
from products_total;
--

drop view products_total;
drop view products_ok;
drop view products_machine;
drop view products_product_type;

drop foreign table product_stream;

drop table products;
drop table statuses;

delete from products;
delete from statuses;
