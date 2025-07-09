CREATE TABLE comments(
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    blog_id INT NOT NULL,
    content VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (blog_id) REFERENCES blogs(id)
)