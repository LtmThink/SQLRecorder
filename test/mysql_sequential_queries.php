<?php
// 创建 MySQL 连接
$conn = new mysqli("127.0.0.1", "root", "123456", "", 43306);

// 检查连接
if ($conn->connect_error) {
    die("Connection failed: " . $conn->connect_error);
}

// 创建数据库
$sql = "CREATE DATABASE IF NOT EXISTS test";
if ($conn->query($sql) === TRUE) {
    echo "Database created successfully.\n";
} else {
    echo "Error creating database: " . $conn->error . "\n";
}

// 切换到数据库
$sql = "USE test";
if ($conn->query($sql) === TRUE) {
    echo "Using database test.\n";
} else {
    echo "Error using database: " . $conn->error . "\n";
}

// 创建表
$sql = "CREATE TABLE IF NOT EXISTS users (name VARCHAR(100))";
if ($conn->query($sql) === TRUE) {
    echo "Table users created successfully.\n";
} else {
    echo "Error creating table: " . $conn->error . "\n";
}

// 插入数据
$sql = "INSERT INTO users (name) VALUES ('test3'),('test4'),('test5')";
if ($conn->query($sql) === TRUE) {
    echo "Records inserted successfully.\n";
} else {
    echo "Error inserting records: " . $conn->error . "\n";
}

// 更新数据
$sql = "UPDATE users SET name = 'test2' WHERE name = 'test3'";
if ($conn->query($sql) === TRUE) {
    echo "Record updated successfully.\n";
} else {
    echo "Error updating record: " . $conn->error . "\n";
}

// 查询数据
$sql = "SELECT * FROM users";
$result = $conn->query($sql);
if ($result->num_rows > 0) {
    while ($row = $result->fetch_assoc()) {
        print_r($row);
    }
    $result->free();
} else {
    echo "No records found.\n";
}

// 删除数据
$sql = "DELETE FROM users WHERE name = 'test2'";
if ($conn->query($sql) === TRUE) {
    echo "Record deleted successfully.\n";
} else {
    echo "Error deleting record: " . $conn->error . "\n";
}

// 删除表
$sql = "DROP TABLE IF EXISTS users";
if ($conn->query($sql) === TRUE) {
    echo "Table users dropped successfully.\n";
} else {
    echo "Error dropping table: " . $conn->error . "\n";
}

// 删除数据库
$sql = "DROP DATABASE IF EXISTS test";
if ($conn->query($sql) === TRUE) {
    echo "Database dropped successfully.\n";
} else {
    echo "Error dropping database: " . $conn->error . "\n";
}
// 模拟报错
$sql = "SELECT 123*aaa;";
if ($conn->query($sql) === TRUE) {
    echo "Records selected successfully.\n";
} else {
    echo "Error selecting records: " . $conn->error . "\n";
}

// 关闭连接
$conn->close();
?>
