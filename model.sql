create database acourse;

set database = acourse;

create table sessions (
  k text,
  v blob,
  e timestamp,
  primary key (k),
  index (e)
);

create table users (
  id string,
  username string not null,
  name string not null,
  email string,
  about_me string not null,
  image string not null,
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),
  primary key (id),
  unique (username),
  unique (email),
  index (created_at),
  index (updated_at)
);

create table roles (
  user_id string,
  admin bool not null default false,
  instructor bool not null default false,
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),
  primary key (user_id),
  foreign key (user_id) references users (id),
  index (admin),
  index (instructor),
  index (created_at),
  index (updated_at)
);

create table courses (
  id serial,
  user_id string not null,
  title string not null,
  short_desc string not null,
  long_desc string not null,
  image string not null,
  start timestamp default null,
  url string default null,
  type int not null,
  price decimal(9,2) not null default 0,
  discount decimal(9,2) default 0,
  enroll_detail string not null default '',
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),
  primary key (id),
  foreign key (user_id) references users (id),
  unique (url),
  index (created_at),
  index (updated_at),
  index (title),
  index (start),
  index (type),
  index (price)
);

create table course_options (
  id int,
  public bool not null default false,
  enroll bool not null default false,
  attend bool not null default false,
  assignment bool not null default false,
  discount bool not null default false,
  primary key (id),
  foreign key (id) references courses (id),
  index (public),
  index (enroll),
  index (public, enroll),
  index (public, discount),
  index (public, discount, enroll)
);

create table course_contents (
  course_id int,
  i int,
  title string not null,
  long_desc string not null,
  video_id string default null,
  video_type int default null,
  download_url string default null,
  primary key (course_id, i),
  foreign key (course_id) references courses (id)
);

create table assignments (
  id serial,
  course_id int not null,
  title string not null,
  long_desc string not null,
  open bool not null default false,
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),
  primary key (id),
  foreign key (course_id) references courses (id),
  index (created_at)
);

create table user_assignments (
  id serial,
  user_id string not null,
  assignment_id int not null,
  download_url string not null,
  created_at timestamp not null default now(),
  primary key (id),
  foreign key (user_id) references users (id),
  foreign key (assignment_id) references assignments (id),
  index (created_at)
);

create table enrolls (
  user_id string,
  course_id int not null,
  created_at timestamp not null default now(),
  primary key (user_id, course_id),
  foreign key (user_id) references users (id),
  foreign key (course_id) references courses (id),
  index (created_at),
  index (user_id, created_at),
  index (course_id, created_at)
);

create table attends (
  id serial,
  user_id string not null,
  course_id int not null,
  created_at timestamp not null default now(),
  primary key (id),
  foreign key (user_id) references users (id),
  foreign key (course_id) references courses (id),
  index (created_at),
  index (user_id, created_at),
  index (course_id, created_at),
  index (user_id, course_id, created_at)
);

create table payments (
  id serial,
  user_id string not null,
  course_id int not null,
  image string not null,
  price decimal(9, 2) not null,
  original_price decimal(9, 2) not null,
  code string not null,
  status int not null,
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),
  at timestamp default null,
  primary key (id),
  foreign key (user_id) references users (id),
  foreign key (course_id) references courses (id),
  index (created_at),
  index (updated_at),
  index (code),
  index (course_id, code),
  index (status, created_at),
  index (status, updated_at),
  index (at)
);
