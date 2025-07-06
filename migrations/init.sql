CREATE TABLE roles (
    id TEXT NOT NULL,
    role TEXT CHECK (role IN ('student', 'parent')),
    PRIMARY KEY (id, role)
);

CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    role TEXT CHECK (role IN ('student', 'parent')) NOT NULL,
    user_id TEXT NOT NULL,
    date TIMESTAMP NOT NULL DEFAULT now(),
    amount NUMERIC(10, 2) NOT NULL,
    quantity INTEGER NOT NULL,
    FOREIGN KEY (user_id, role) REFERENCES roles(id, role) ON DELETE CASCADE
);