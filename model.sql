CREATE DATABASE IF NOT EXISTS acourse;

SET DATABASE = acourse;

CREATE TABLE IF NOT EXISTS users (
  id STRING PRIMARY KEY,
  username STRING UNIQUE NOT NULL,
  name STRING,
  email STRING UNIQUE,
  about_me STRING,
  image STRING,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),
  INDEX (username, email, created_at, updated_at)
);

CREATE TABLE IF NOT EXISTS roles (
  id STRING PRIMARY KEY REFERENCES users (id),
  admin BOOL NOT NULL DEFAULT false,
  instructor BOOL NOT NULL DEFAULT false,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),
  INDEX (admin, instructor, created_at, updated_at)
);

CREATE TABLE IF NOT EXISTS courses (
  id SERIAL PRIMARY KEY,
  user_id STRING NOT NULL REFERENCES users (id),
  title STRING NOT NULL,
  short_desc STRING,
  long_desc STRING,
  image STRING,
  start TIMESTAMP,
  url STRING UNIQUE,
  type int NOT NULL,
  price DECIMAL(9,2) NOT NULL DEFAULT 0,
  discount DECIMAL(9,2),
  enroll_detail STRING,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),
  INDEX (title, start, url, type, user_id, created_at, updated_at)
);

CREATE TABLE IF NOT EXISTS course_options (
  id INT PRIMARY KEY REFERENCES courses (id),
  public BOOL NOT NULL DEFAULT false,
  enroll BOOL NOT NULL DEFAULT false,
  attend BOOL NOT NULL DEFAULT false,
  assignment BOOL NOT NULL DEFAULT false,
  discount BOOL NOT NULL DEFAULT false,
  INDEX (public, enroll, attend, assignment, discount)
);

CREATE TABLE IF NOT EXISTS course_contents (
  id SERIAL PRIMARY KEY,
  course_id INT NOT NULL REFERENCES courses (id),
  title STRING NOT NULL,
  long_desc STRING,
  video_id STRING,
  video_type INT,
  download_url STRING,
  INDEX (course_id)
);

CREATE TABLE IF NOT EXISTS course_assignments (
  id SERIAL PRIMARY KEY,
  course_id INT NOT NULL REFERENCES courses (id),
  title STRING NOT NULL,
  long_desc STRING,
  open BOOL NOT NULL DEFAULT false,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),
  INDEX (course_id, open, created_at)
);

CREATE TABLE IF NOT EXISTS user_assignments (
  id SERIAL PRIMARY KEY,
  user_id STRING NOT NULL REFERENCES users (id),
  course_assignment_id INT NOT NULL REFERENCES course_assignments (id),
  download_url STRING NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  INDEX (user_id, course_assignment_id, created_at)
);

CREATE TABLE IF NOT EXISTS payments (
  id SERIAL PRIMARY KEY,
  user_id STRING NOT NULL REFERENCES users (id),
  course_id INT NOT NULL REFERENCES courses (id),
  image STRING NOT NULL,
  price DECIMAL(9, 2) NOT NULL,
  original_price DECIMAL(9, 2) NOT NULL,
  code STRING,
  status INT,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),
  at TIMESTAMP,
  INDEX (user_id, course_id, code, status, created_at, updated_at, at)
);
