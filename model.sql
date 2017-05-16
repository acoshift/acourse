CREATE DATABASE acourse;

SET DATABASE = acourse;

CREATE TABLE users (
  id STRING PRIMARY KEY NOT NULL,
  username STRING UNIQUE NOT NULL,
  name STRING NOT NULL,
  email STRING UNIQUE,
  about_me STRING NOT NULL,
  image STRING NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),
  INDEX (username, email, created_at, updated_at)
);

CREATE TABLE roles (
  id STRING PRIMARY KEY NOT NULL REFERENCES users (id),
  admin BOOL NOT NULL DEFAULT false,
  instructor BOOL NOT NULL DEFAULT false,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),
  INDEX (admin, instructor, created_at, updated_at)
);

CREATE TABLE courses (
  id SERIAL PRIMARY KEY NOT NULL,
  user_id STRING NOT NULL REFERENCES users (id),
  title STRING NOT NULL,
  short_desc STRING NOT NULL,
  long_desc STRING NOT NULL,
  image STRING NOT NULL,
  start TIMESTAMP DEFAULT NULL,
  url STRING UNIQUE DEFAULT NULL,
  type int NOT NULL,
  price DECIMAL(9,2) NOT NULL DEFAULT 0,
  discount DECIMAL(9,2) DEFAULT 0,
  enroll_detail STRING NOT NULL DEFAULT '',
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),
  INDEX (title, start, url, type, user_id, created_at, updated_at)
);

CREATE TABLE course_options (
  id INT PRIMARY KEY NOT NULL REFERENCES courses (id),
  public BOOL NOT NULL DEFAULT false,
  enroll BOOL NOT NULL DEFAULT false,
  attend BOOL NOT NULL DEFAULT false,
  assignment BOOL NOT NULL DEFAULT false,
  discount BOOL NOT NULL DEFAULT false,
  INDEX (public, enroll, attend, assignment, discount)
);

CREATE TABLE course_contents (
  id SERIAL PRIMARY KEY NOT NULL,
  course_id INT NOT NULL REFERENCES courses (id),
  title STRING NOT NULL,
  long_desc STRING NOT NULL,
  video_id STRING DEFAULT NULL,
  video_type INT DEFAULT NULL,
  download_url STRING DEFAULT NULL,
  INDEX (course_id)
);

CREATE TABLE assignments (
  id SERIAL PRIMARY KEY NOT NULL,
  course_id INT NOT NULL REFERENCES courses (id),
  title STRING NOT NULL,
  long_desc STRING NOT NULL,
  open BOOL NOT NULL DEFAULT false,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),
  INDEX (course_id, open, created_at)
);

CREATE TABLE user_assignments (
  id SERIAL PRIMARY KEY NOT NULL,
  user_id STRING NOT NULL REFERENCES users (id),
  assignment_id INT NOT NULL REFERENCES assignments (id),
  download_url STRING NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  INDEX (user_id, assignment_id, created_at)
);

CREATE TABLE enrolls (
  id SERIAL PRIMARY KEY NOT NULL,
  user_id STRING NOT NULL REFERENCES users (id),
  course_id INT NOT NULL REFERENCES courses (id),
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  INDEX (user_id, course_id, created_at),
  UNIQUE (user_id, course_id)
);

CREATE TABLE attends (
  id SERIAL PRIMARY KEY NOT NULL,
  user_id STRING NOT NULL REFERENCES users (id),
  course_id INT NOT NULL REFERENCES courses (id),
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  INDEX (user_id, course_id, created_at)
);

CREATE TABLE payments (
  id SERIAL PRIMARY KEY NOT NULL,
  user_id STRING NOT NULL REFERENCES users (id),
  course_id INT NOT NULL REFERENCES courses (id),
  image STRING NOT NULL,
  price DECIMAL(9, 2) NOT NULL,
  original_price DECIMAL(9, 2) NOT NULL,
  code STRING NOT NULL,
  status INT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),
  at TIMESTAMP DEFAULT NULL,
  INDEX (user_id, course_id, code, status, created_at, updated_at, at)
);
