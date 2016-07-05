create table t_user_state (
  user_id char(8) not null,
  date char(8) not null,
  time char(4) not null,
  state char(1) not null,
  primary key(user_id, date, time)
);
