<?php
// 创建 MySQL 连接
$conn = new mysqli("127.0.0.1", "root", "123456", "", 43306);

// 检查连接
if ($conn->connect_error) {
    die("Connection failed: " . $conn->connect_error);
}

// 定义 SQL 语句
$sql = "
CREATE DATABASE IF NOT EXISTS test;
USE test;
CREATE TABLE IF NOT EXISTS users (name VARCHAR(100));
INSERT INTO users (name) VALUES ('张三'),('李四'),('王五');
UPDATE users SET name = '熊二' WHERE name = '张三';
SELECT * FROM users;
DELETE FROM users WHERE name = '熊二';
DROP TABLE IF EXISTS users;
DROP DATABASE IF EXISTS test;

SELECT 123*aaa;
";

// 执行多条 SQL 语句
if ($conn->multi_query($sql)) {
    do {
        if ($result = $conn->store_result()) {
            while ($row = $result->fetch_assoc()) {
                print_r($row);
            }
            $result->free();
        }
    } while ($conn->next_result());
} else {
    echo "Error: " . $conn->error;
}

// 关闭连接
$conn->close();
?>
