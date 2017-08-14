create database acourse;

set database = acourse;

create table users (
  id varchar not null,
  username varchar not null,
  name varchar not null,
  email varchar,
  about_me varchar not null default '',
  image varchar not null default '',
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),
  primary key (id)
);
create unique index users_username_key on users (username);
create unique index users_email_key on users (email);
create index users_created_at_idx on users (created_at desc);

create table roles (
  user_id varchar,
  admin bool not null default false,
  instructor bool not null default false,
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),
  primary key (user_id),
  foreign key (user_id) references users (id)
);
create index roles_admin_idx on roles (admin);
create index roles_instructor_idx on roles (instructor);

create table courses (
  id serial,
  user_id varchar not null,
  title varchar not null,
  short_desc varchar not null,
  long_desc varchar not null,
  image varchar not null,
  start timestamp default null,
  url varchar default null,
  type int not null default 0,
  price decimal(9,2) not null default 0,
  discount decimal(9,2) default 0,
  enroll_detail varchar not null default '',
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),
  primary key (id),
  foreign key (user_id) references users (id)
);
create unique index courses_url_key on courses (url);
create index courses_created_at_idx on courses (created_at desc);
create index courses_updated_at_idx on courses (updated_at desc);

create table course_options (
  course_id int,
  public bool not null default false,
  enroll bool not null default false,
  attend bool not null default false,
  assignment bool not null default false,
  discount bool not null default false,
  primary key (course_id),
  foreign key (course_id) references courses (id)
);
create index course_options_public_idx on course_options (public);
create index course_options_enroll_idx on course_options (enroll);
create index course_options_public_enroll_idx on course_options (public, enroll);
create index course_options_public_discount_idx on course_options (public, discount);
create index course_options_public_discount_enroll_idx on course_options (public, discount, enroll);

create table course_contents (
  id serial,
  course_id int not null,
  i int not null default 0,
  title varchar not null default '',
  long_desc varchar not null default '',
  video_id varchar not null default '',
  video_type int not null default 0,
  download_url varchar not null default '',
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),
  primary key (id),
  foreign key (course_id) references courses (id)
);
create index course_contents_course_id_i_idx on course_contents (course_id, i);

create table assignments (
  id serial,
  course_id int not null,
  i int not null,
  title varchar not null,
  long_desc varchar not null,
  open bool not null default false,
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),
  primary key (id),
  foreign key (course_id) references courses (id)
);
create index assignments_course_id_idx on assignments (course_id, i);

create table user_assignments (
  id serial,
  user_id varchar not null,
  assignment_id int not null,
  download_url varchar not null,
  created_at timestamp not null default now(),
  primary key (id),
  foreign key (user_id) references users (id),
  foreign key (assignment_id) references assignments (id)
);
create index user_assignments_created_at_idx on user_assignments (created_at);

create table enrolls (
  user_id varchar,
  course_id int not null,
  created_at timestamp not null default now(),
  primary key (user_id, course_id),
  foreign key (user_id) references users (id),
  foreign key (course_id) references courses (id)
);
create index enrolls_created_at_idx on enrolls (created_at);
create index enrolls_user_id_created_at_idx on enrolls (user_id, created_at);
create index enrolls_course_id_created_at_idx on enrolls (course_id, created_at);

create table attends (
  id serial,
  user_id varchar not null,
  course_id int not null,
  created_at timestamp not null default now(),
  primary key (id),
  foreign key (user_id) references users (id),
  foreign key (course_id) references courses (id)
);
create index attends_created_at_idx on attends (created_at);
create index attends_user_id_created_at_idx on attends (user_id, created_at);
create index attends_course_id_created_at_idx on attends (course_id, created_at);
create index attends_user_id_course_id_created_at_idx on attends (user_id, course_id, created_at);

create table payments (
  id serial,
  user_id varchar not null,
  course_id int not null,
  image varchar not null,
  price decimal(9, 2) not null,
  original_price decimal(9, 2) not null,
  code varchar not null,
  status int not null,
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),
  at timestamp default null,
  primary key (id),
  foreign key (user_id) references users (id),
  foreign key (course_id) references courses (id)
);
create index payments_created_at_idx on payments (created_at desc);
create index payments_code_idx on payments (code);
create index payments_course_id_code_idx on payments (course_id, code);
create index payments_status_created_at_idx on payments (status, created_at desc);
