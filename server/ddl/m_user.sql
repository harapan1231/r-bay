create table m_user (
  user_id char(8) not null,
  user_name varchar(40) not null,
  job_type char(1) not null,
  team_id char(5),
  primary key(user_id)
);
