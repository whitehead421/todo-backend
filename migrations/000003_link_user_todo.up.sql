DELETE FROM
    todos;

ALTER TABLE
    todos
ADD
    COLUMN user_id INT NOT NULL;

ALTER TABLE
    todos
ADD
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id);