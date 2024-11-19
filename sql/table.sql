-- user テーブル
CREATE TABLE user (
    user_id VARCHAR(255) PRIMARY KEY,
    firebase_uid VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    bio TEXT,
    profile_img_url VARCHAR(2083) -- URLの最大長に合わせる
);

-- post テーブル
CREATE TABLE post (
    post_id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255),
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    edited_at DATETIME,
    deleted_at DATETIME,
    parent_post_id VARCHAR(255),
    FOREIGN KEY (user_id) REFERENCES user(user_id) ON DELETE CASCADE,
    FOREIGN KEY (parent_post_id) REFERENCES post(post_id) ON DELETE SET NULL
);

-- like テーブル
CREATE TABLE `like` (
    user_id VARCHAR(255),
    post_id VARCHAR(255),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, post_id),
    FOREIGN KEY (user_id) REFERENCES user(user_id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES post(post_id) ON DELETE CASCADE
);

-- follower テーブル
CREATE TABLE follower (
    user_id VARCHAR(255),
    followed_user_id VARCHAR(255),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, followed_user_id),
    FOREIGN KEY (user_id) REFERENCES user(user_id) ON DELETE CASCADE,
    FOREIGN KEY (followed_user_id) REFERENCES user(user_id) ON DELETE CASCADE
);
