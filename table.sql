CREATE TABLE users (
    id          VARCHAR NOT NULL,
    username    VARCHAR NOT NULL,
    name        VARCHAR NOT NULL,
    email       VARCHAR,
    about_me    VARCHAR NOT NULL DEFAULT '',
    image       VARCHAR NOT NULL DEFAULT '',
    created_at  TIMESTAMP NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY (id)
);
CREATE UNIQUE INDEX ON users (username);
CREATE UNIQUE INDEX ON users (email);
CREATE INDEX ON users (created_at DESC);

CREATE TABLE roles (
  user_id       VARCHAR,
  admin         BOOL NOT NULL DEFAULT FALSE,
  instructor    BOOL NOT NULL DEFAULT FALSE,
  created_at    TIMESTAMP NOT NULL DEFAULT now(),
  updated_at    TIMESTAMP NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id),
  FOREIGN KEY (user_id) REFERENCES users (id)
);
CREATE INDEX ON roles (admin);
CREATE INDEX ON roles (instructor);

CREATE TABLE courses (
  id            UUID DEFAULT gen_random_uuid(),
  user_id       VARCHAR NOT NULL,
  title         VARCHAR NOT NULL,
  short_desc    VARCHAR NOT NULL,
  long_desc     VARCHAR NOT NULL,
  image         VARCHAR NOT NULL,
  start         TIMESTAMP DEFAULT NULL,
  url           VARCHAR DEFAULT NULL,
  type          INT NOT NULL DEFAULT 0,
  price         DECIMAL(9,2) NOT NULL DEFAULT 0,
  discount      DECIMAL(9,2) DEFAULT 0,
  enroll_detail VARCHAR NOT NULL DEFAULT '',
  created_at    TIMESTAMP NOT NULL DEFAULT now(),
  updated_at    TIMESTAMP NOT NULL DEFAULT now(),
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users (id)
);
CREATE UNIQUE INDEX ON courses (url);
CREATE INDEX ON courses (created_at DESC);
CREATE INDEX ON courses (updated_at DESC);

CREATE TABLE course_options (
  course_id     UUID,
  public        BOOL NOT NULL DEFAULT FALSE,
  enroll        BOOL NOT NULL DEFAULT FALSE,
  attend        BOOL NOT NULL DEFAULT FALSE,
  assignment    BOOL NOT NULL DEFAULT FALSE,
  discount      BOOL NOT NULL DEFAULT FALSE,
  PRIMARY KEY (course_id),
  FOREIGN KEY (course_id) REFERENCES courses (id)
);
CREATE INDEX ON course_options (public);
CREATE INDEX ON course_options (enroll);
CREATE INDEX ON course_options (public, enroll);
CREATE INDEX ON course_options (public, discount);
CREATE INDEX ON course_options (public, discount, enroll);

CREATE TABLE course_contents (
  id            UUID DEFAULT gen_random_uuid(),
  course_id     UUID NOT NULL,
  i             INT NOT NULL DEFAULT 0,
  title         VARCHAR NOT NULL DEFAULT '',
  long_desc     VARCHAR NOT NULL DEFAULT '',
  video_id      VARCHAR NOT NULL DEFAULT '',
  video_type    INT NOT NULL DEFAULT 0,
  download_url  VARCHAR NOT NULL DEFAULT '',
  created_at    TIMESTAMP NOT NULL DEFAULT now(),
  updated_at    TIMESTAMP NOT NULL DEFAULT now(),
  PRIMARY KEY (id),
  FOREIGN KEY (course_id) REFERENCES courses (id)
);
CREATE INDEX ON course_contents (course_id, i);

CREATE TABLE assignments (
  id            UUID DEFAULT gen_random_uuid(),
  course_id     UUID NOT NULL,
  i             INT NOT NULL,
  title         VARCHAR NOT NULL,
  long_desc     VARCHAR NOT NULL,
  open          BOOL NOT NULL DEFAULT FALSE,
  created_at    TIMESTAMP NOT NULL DEFAULT now(),
  updated_at    TIMESTAMP NOT NULL DEFAULT now(),
  PRIMARY KEY (id),
  FOREIGN KEY (course_id) REFERENCES courses (id)
);
CREATE INDEX ON assignments (course_id, i);

CREATE TABLE user_assignments (
  id            UUID DEFAULT gen_random_uuid(),
  user_id       VARCHAR NOT NULL,
  assignment_id UUID NOT NULL,
  download_url  VARCHAR NOT NULL,
  created_at    TIMESTAMP NOT NULL DEFAULT now(),
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users (id),
  FOREIGN KEY (assignment_id) REFERENCES assignments (id)
);
CREATE INDEX ON user_assignments (created_at);

CREATE TABLE enrolls (
  user_id       VARCHAR,
  course_id     UUID NOT NULL,
  created_at    TIMESTAMP NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, course_id),
  FOREIGN KEY (user_id) REFERENCES users (id),
  FOREIGN KEY (course_id) REFERENCES courses (id)
);
CREATE INDEX ON enrolls (created_at);
CREATE INDEX ON enrolls (user_id, created_at);
CREATE INDEX ON enrolls (course_id, created_at);

CREATE TABLE attends (
  id            UUID DEFAULT gen_random_uuid(),
  user_id       VARCHAR NOT NULL,
  course_id     UUID NOT NULL,
  created_at    TIMESTAMP NOT NULL DEFAULT now(),
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users (id),
  FOREIGN KEY (course_id) REFERENCES courses (id)
);
CREATE INDEX ON attends (created_at);
CREATE INDEX ON attends (user_id, created_at);
CREATE INDEX ON attends (course_id, created_at);
CREATE INDEX ON attends (user_id, course_id, created_at);

CREATE TABLE payments (
  id                UUID DEFAULT gen_random_uuid(),
  user_id           VARCHAR NOT NULL,
  course_id         UUID NOT NULL,
  image             VARCHAR NOT NULL,
  price             DECIMAL(9, 2) NOT NULL,
  original_price    DECIMAL(9, 2) NOT NULL,
  code              VARCHAR NOT NULL,
  status            INT NOT NULL,
  created_at        TIMESTAMP NOT NULL DEFAULT now(),
  updated_at        TIMESTAMP NOT NULL DEFAULT now(),
  at                TIMESTAMP DEFAULT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users (id),
  FOREIGN KEY (course_id) REFERENCES courses (id)
);
CREATE INDEX ON payments (created_at DESC);
CREATE INDEX ON payments (code);
CREATE INDEX ON payments (course_id, code);
CREATE INDEX ON payments (status, created_at DESC);
