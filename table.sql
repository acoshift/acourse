-- create extension if not exists pgcrypto;

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
create unique index on users (username);
create unique index on users (email);
create index on users (created_at desc);

create table roles (
	user_id varchar,
	admin bool not null default false,
	instructor bool not null default false,
	created_at timestamp not null default now(),
	updated_at timestamp not null default now(),
	primary key (user_id),
	foreign key (user_id) references users (id)
);
create index on roles (admin);
create index on roles (instructor);

create table courses (
	id uuid default gen_random_uuid(),
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
create unique index on courses (url);
create index on courses (created_at desc);
create index on courses (updated_at desc);

create table course_options (
	course_id uuid,
	public bool not null default false,
	enroll bool not null default false,
	attend bool not null default false,
	assignment bool not null default false,
	discount bool not null default false,
	primary key (course_id),
	foreign key (course_id) references courses (id)
);
create index on course_options (public);
create index on course_options (enroll);
create index on course_options (public, enroll);
create index on course_options (public, discount);
create index on course_options (public, discount, enroll);

create table course_contents (
	id uuid default gen_random_uuid(),
	course_id uuid not null,
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
create index on course_contents (course_id, i);

create table assignments (
	id uuid default gen_random_uuid(),
	course_id uuid not null,
	i int not null,
	title varchar not null,
	long_desc varchar not null,
	open bool not null default false,
	created_at timestamp not null default now(),
	updated_at timestamp not null default now(),
	primary key (id),
	foreign key (course_id) references courses (id)
);
create index on assignments (course_id, i);

create table user_assignments (
	id uuid default gen_random_uuid(),
	user_id varchar not null,
	assignment_id uuid not null,
	download_url varchar not null,
	created_at timestamp not null default now(),
	primary key (id),
	foreign key (user_id) references users (id),
	foreign key (assignment_id) references assignments (id)
);
create index on user_assignments (created_at);

create table enrolls (
	user_id varchar,
	course_id uuid not null,
	created_at timestamp not null default now(),
	primary key (user_id, course_id),
	foreign key (user_id) references users (id),
	foreign key (course_id) references courses (id)
);
create index on enrolls (created_at);
create index on enrolls (user_id, created_at);
create index on enrolls (course_id, created_at);

create table attends (
	id uuid default gen_random_uuid(),
	user_id varchar not null,
	course_id uuid not null,
	created_at timestamp not null default now(),
	primary key (id),
	foreign key (user_id) references users (id),
	foreign key (course_id) references courses (id)
);
create index on attends (created_at);
create index on attends (user_id, created_at);
create index on attends (course_id, created_at);
create index on attends (user_id, course_id, created_at);

create table payments (
	id uuid default gen_random_uuid(),
	user_id varchar not null,
	course_id uuid not null,
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
create index on payments (created_at desc);
create index on payments (code);
create index on payments (course_id, code);
create index on payments (status, created_at desc);
