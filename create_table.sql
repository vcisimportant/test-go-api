DROP TABLE IF EXISTS characters;

CREATE TABLE characters (
  id            SERIAL PRIMARY KEY,
  name          VARCHAR(128) NOT NULL,
  description   VARCHAR(255) NOT NULL,
  book          VARCHAR(255) NOT NULL, 
  rating        INT NOT NULL
);

INSERT INTO characters (name, description, book, rating)
VALUES
  ('Darrow', 'true red god', 'Red Rising', 10),
  ('Sevro', 'OG howler', 'Red Rising', 9),
  ('Lyssander', 'true piece of shit', 'Red Rising', 3),
  ('Kaladin', 'bridgement and herald', 'Stormlight', 10);
