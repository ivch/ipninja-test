create table notes
(
  id         INTEGER primary key autoincrement,
  title      TEXT     not null,
  body       TEXT     not null,
  created_at DATETIME not null,
  expires_at DATETIME,
  canceled   tinyint default 0
);

