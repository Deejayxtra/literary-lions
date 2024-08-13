-- Create the 'users' table to store user account information.
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,       -- Unique identifier for each user, auto-incremented.
    email TEXT NOT NULL UNIQUE,                 -- User's email, must be unique and not null.
    username TEXT NOT NULL UNIQUE,              -- User's username, must be unique and not null.
    password TEXT NOT NULL,                     -- Hashed password for user authentication, not null.
    role TEXT NOT NULL CHECK (role IN ('user', 'admin'))  -- User role, constrained to either 'user' or 'admin'.
);

-- Create the 'sessions' table to store user sessions for authentication.
CREATE TABLE IF NOT EXISTS sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,       -- Unique identifier for each session, auto-incremented.
    user_id INTEGER NOT NULL,                   -- Foreign key referencing the 'users' table, links session to a user.
    token TEXT NOT NULL,                        -- Session token used for authentication.
    expires_at DATETIME NOT NULL,               -- Expiration time of the session token.
    FOREIGN KEY (user_id) REFERENCES users(id)  -- Ensure user_id corresponds to a valid user in the 'users' table.
);

-- Create the 'posts' table to store blog posts made by users.
CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,       -- Unique identifier for each post, auto-incremented.
    user_id INTEGER NOT NULL,                   -- Foreign key referencing the 'users' table, links post to a user.
    category TEXT,                              -- Category of the post, can be null.
    title TEXT NOT NULL,                        -- Title of the post, not null.
    content TEXT NOT NULL,                      -- Content of the post, not null.
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,  -- Timestamp of post creation, defaults to current time.
    FOREIGN KEY (user_id) REFERENCES users(id)  -- Ensure user_id corresponds to a valid user in the 'users' table.
);

-- Create the 'comments' table to store comments made on posts.
CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,       -- Unique identifier for each comment, auto-incremented.
    post_id INTEGER NOT NULL,                   -- Foreign key referencing the 'posts' table, links comment to a post.
    user_id INTEGER NOT NULL,                   -- Foreign key referencing the 'users' table, links comment to a user.
    content TEXT NOT NULL,                      -- Content of the comment, not null.
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,  -- Timestamp of comment creation, defaults to current time.
    FOREIGN KEY (post_id) REFERENCES posts(id), -- Ensure post_id corresponds to a valid post in the 'posts' table.
    FOREIGN KEY (user_id) REFERENCES users(id)  -- Ensure user_id corresponds to a valid user in the 'users' table.
);

-- Create the 'post_likes' table to track likes on posts by users.
CREATE TABLE IF NOT EXISTS post_likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,       -- Unique identifier for each post like, auto-incremented.
    user_id INTEGER NOT NULL,                   -- Foreign key referencing the 'users' table, links like to a user.
    post_id INTEGER,                            -- Foreign key referencing the 'posts' table, links like to a post.
    is_like BOOLEAN NOT NULL,                   -- Boolean indicating whether it is a like (true) or not (false).
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,  -- Timestamp of like creation, defaults to current time.
    FOREIGN KEY (user_id) REFERENCES users(id), -- Ensure user_id corresponds to a valid user in the 'users' table.
    FOREIGN KEY (post_id) REFERENCES posts(id)  -- Ensure post_id corresponds to a valid post in the 'posts' table.
);

-- Create the 'post_dislikes' table to track dislikes on posts by users.
CREATE TABLE IF NOT EXISTS post_dislikes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,       -- Unique identifier for each post dislike, auto-incremented.
    user_id INTEGER NOT NULL,                   -- Foreign key referencing the 'users' table, links dislike to a user.
    post_id INTEGER,                            -- Foreign key referencing the 'posts' table, links dislike to a post.
    is_dislike BOOLEAN NOT NULL,                -- Boolean indicating whether it is a dislike (true) or not (false).
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,  -- Timestamp of dislike creation, defaults to current time.
    FOREIGN KEY (user_id) REFERENCES users(id), -- Ensure user_id corresponds to a valid user in the 'users' table.
    FOREIGN KEY (post_id) REFERENCES posts(id)  -- Ensure post_id corresponds to a valid post in the 'posts' table.
);

-- Create the 'comment_likes' table to track likes on comments by users.
CREATE TABLE IF NOT EXISTS comment_likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,       -- Unique identifier for each comment like, auto-incremented.
    user_id INTEGER NOT NULL,                   -- Foreign key referencing the 'users' table, links like to a user.
    comment_id INTEGER,                         -- Foreign key referencing the 'comments' table, links like to a comment.
    is_like BOOLEAN NOT NULL,                   -- Boolean indicating whether it is a like (true) or not (false).
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,  -- Timestamp of like creation, defaults to current time.
    FOREIGN KEY (user_id) REFERENCES users(id), -- Ensure user_id corresponds to a valid user in the 'users' table.
    FOREIGN KEY (comment_id) REFERENCES comments(id) -- Ensure comment_id corresponds to a valid comment in the 'comments' table.
);

-- Create the 'comment_dislikes' table to track dislikes on comments by users.
CREATE TABLE IF NOT EXISTS comment_dislikes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,       -- Unique identifier for each comment dislike, auto-incremented.
    user_id INTEGER NOT NULL,                   -- Foreign key referencing the 'users' table, links dislike to a user.
    comment_id INTEGER,                         -- Foreign key referencing the 'comments' table, links dislike to a comment.
    is_dislike BOOLEAN NOT NULL,                -- Boolean indicating whether it is a dislike (true) or not (false).
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,  -- Timestamp of dislike creation, defaults to current time.
    FOREIGN KEY (user_id) REFERENCES users(id), -- Ensure user_id corresponds to a valid user in the 'users' table.
    FOREIGN KEY (comment_id) REFERENCES comments(id) -- Ensure comment_id corresponds to a valid comment in the 'comments' table.
);
