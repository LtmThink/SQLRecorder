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

// 创建数据库
$query = "CREATE DATABASE test;";
$result = pg_query($conn, $query);

// 检查创建数据库是否成功
if (!$result) {
    echo "Error: " . pg_last_error();
    pg_close($conn);
    exit;
}

// 关闭当前连接，连接到新创建的数据库
pg_close($conn);
$conn = pg_connect("host=".$host." port=".$port." dbname=test user=".$user." password=".$password);

if (!$conn) {
    die("Connection failed: " . pg_last_error());
}

// 开始事务
pg_query($conn, "BEGIN");

$query = "
    CREATE TABLE IF NOT EXISTS users (name VARCHAR(100));
    INSERT INTO users (name) VALUES ('test3'), ('test4'), ('test5');
    UPDATE users SET name = 'test2' WHERE name = 'test3';
    SELECT * FROM users;
    DELETE FROM users WHERE name = 'test2';
    DROP TABLE IF EXISTS users;
    select a;
";

// 执行查询
$result = pg_query($conn, $query);

// 检查查询是否成功
if (!$result) {
    echo "Error: " . pg_last_error();
    // 如果出现错误，回滚事务
    pg_query($conn, "ROLLBACK");
} else {
    // 提交事务
    pg_query($conn, "COMMIT");
}

// 关闭连接到新数据库
pg_close($conn);

// 重新连接到默认数据库（postgres）
$conn = pg_connect("host=".$host." port=".$port." dbname=postgres user=".$user." password=".$password);

if (!$conn) {
    die("Connection failed: " . pg_last_error());
}

// 删除数据库
$query = "DROP DATABASE IF EXISTS test;";
$result = pg_query($conn, $query);

// 检查删除数据库是否成功
if (!$result) {
    echo "Error: " . pg_last_error();
} else {
    echo "Database 'test' deleted successfully.";
}

// 关闭连接
pg_close($conn);
?>