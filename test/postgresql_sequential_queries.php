<?php
$host = "127.0.0.1";
$port = "45432";
$user="user_7n2xtj";
$password="password_QKPYJb";
// 创建 PostgreSQL 连接
$conn = pg_connect("host=".$host." port=".$port." dbname=postgres user=".$user." password=".$password);
if (!$conn) {
    echo "Error: Unable to connect to PostgreSQL.\n";
    exit;
}

// 创建数据库（需要检查数据库是否已存在）
$sql = "SELECT 1 FROM pg_database WHERE datname = 'test'";
$result = pg_query($conn, $sql);
if (pg_num_rows($result) == 0) {
    $sql = "CREATE DATABASE test";
    if (pg_query($conn, $sql)) {
        echo "Database created successfully.\n";
    } else {
        echo "Error creating database: " . pg_last_error($conn) . "\n";
    }
} else {
    echo "Database 'test' already exists.\n";
}

// 连接到目标数据库 'test'
pg_close($conn); // 关闭当前连接
$conn = pg_connect("host=".$host." port=".$port." dbname=test user=".$user." password=".$password);

if (!$conn) {
    echo "Error: Unable to connect to database 'test'.\n";
    exit;
}

// 创建表
$sql = "CREATE TABLE IF NOT EXISTS users (name VARCHAR(100))";
if (pg_query($conn, $sql)) {
    echo "Table users created successfully.\n";
} else {
    echo "Error creating table: " . pg_last_error($conn) . "\n";
}

// 插入数据
$sql = "INSERT INTO users (name) VALUES ('test3'), ('test4'), ('test5')";
if (pg_query($conn, $sql)) {
    echo "Records inserted successfully.\n";
} else {
    echo "Error inserting records: " . pg_last_error($conn) . "\n";
}

// 更新数据
$sql = "UPDATE users SET name = 'test2' WHERE name = 'test3'";
if (pg_query($conn, $sql)) {
    echo "Record updated successfully.\n";
} else {
    echo "Error updating record: " . pg_last_error($conn) . "\n";
}

// 查询数据
$sql = "SELECT * FROM users";
$result = pg_query($conn, $sql);
if (pg_num_rows($result) > 0) {
    while ($row = pg_fetch_assoc($result)) {
        print_r($row);
    }
    pg_free_result($result);
} else {
    echo "No records found.\n";
}

// 删除数据
$sql = "DELETE FROM users WHERE name = 'test2'";
if (pg_query($conn, $sql)) {
    echo "Record deleted successfully.\n";
} else {
    echo "Error deleting record: " . pg_last_error($conn) . "\n";
}

// 删除表
$sql = "DROP TABLE IF EXISTS users";
if (pg_query($conn, $sql)) {
    echo "Table users dropped successfully.\n";
} else {
    echo "Error dropping table: " . pg_last_error($conn) . "\n";
}

// 删除数据库（请确保没有连接到要删除的数据库）
pg_close($conn); // 关闭当前连接
$conn = pg_connect("host=".$host." port=".$port." dbname=postgres user=".$user." password=".$password);
$sql = "DROP DATABASE IF EXISTS test";
if (pg_query($conn, $sql)) {
    echo "Database dropped successfully.\n";
} else {
    echo "Error dropping database: " . pg_last_error($conn) . "\n";
}

// 模拟报错（会报错，因为 'aaa' 不是有效的列名）
$sql = "SELECT 123 * aaa;";
if (pg_query($conn, $sql)) {
    echo "Records selected successfully.\n";
} else {
    echo "Error selecting records: " . pg_last_error($conn) . "\n";
}

// 关闭连接
pg_close($conn);
?>