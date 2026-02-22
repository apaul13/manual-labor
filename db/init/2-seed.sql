BEGIN;

-- Makes (UPPERCASE)
INSERT INTO make (id, name, year) VALUES
  (1, 'TOYOTA', '2020'),
  (2, 'HONDA', '2021'),
  (3, 'FORD',  '2019');

-- Models (UPPERCASE)
INSERT INTO model (id, name, make_id) VALUES
  (1, 'CAMRY', 1),
  (2, 'COROLLA', 1),
  (3, 'CIVIC', 2),
  (4, 'ACCORD', 2),
  (5, 'F-150', 3);

-- Manuals (URLs unchanged)
INSERT INTO manual (id, url) VALUES
  (1, 'https://example.com/manuals/camry-le.pdf'),
  (2, 'https://example.com/manuals/camry-xle.pdf'),
  (3, 'https://example.com/manuals/corolla-se.pdf'),
  (4, 'https://example.com/manuals/corolla-xse.pdf'),
  (5, 'https://example.com/manuals/civic-lx.pdf'),
  (6, 'https://example.com/manuals/civic-ex.pdf'),
  (7, 'https://example.com/manuals/accord-sport.pdf'),
  (8, 'https://example.com/manuals/accord-touring.pdf'),
  (9, 'https://example.com/manuals/f150-xl.pdf'),
  (10,'https://example.com/manuals/f150-xlt.pdf');

-- Trims (UPPERCASE)
INSERT INTO trim (id, name, model_id, manual_id) VALUES
  (1, 'LE',        1, 1),
  (2, 'XLE',       1, 2),
  (3, 'SE',        2, 3),
  (4, 'XSE',       2, 4),
  (5, 'LX',        3, 5),
  (6, 'EX',        3, 6),
  (7, 'SPORT',     4, 7),
  (8, 'TOURING',   4, 8),
  (9, 'XL',        5, 9),
  (10,'XLT',       5, 10);

COMMIT;