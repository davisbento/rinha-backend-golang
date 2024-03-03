create table extratos (
  id SERIAL PRIMARY KEY,
  cliente_id INT REFERENCES clientes(id),
  valor INT,
  tipo VARCHAR(1),
  descricao VARCHAR(100),
  data TIMESTAMP NOT NULL default NOW()
);