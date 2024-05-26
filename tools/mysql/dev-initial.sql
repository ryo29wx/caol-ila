-- create database
CREATE DATABASE IF NOT EXISTS caolila;
USE caolila;

-- create user info table
CREATE TABLE IF NOT EXISTS user_info (
    user_id VARCHAR(36) PRIMARY KEY,
    first_name VARCHAR(36) NOT NULL,
    last_name VARCHAR(36) NOT NULL,
    mail_address1 VARCHAR(128) NOT NULL UNIQUE,
    mail_address2 VARCHAR(128),
    mail_address3 VARCHAR(128),
    phone_num1  VARCHAR(16) NOT NULL,
    phone_num2  VARCHAR(16),
    phone_num3  VARCHAR(16),
    address1 VARCHAR(128),
    address2 VARCHAR(128),
    address3 VARCHAR(128),
    post_code INT,
    sex INT NOT NULL,
    pay_rank INT NOT NULL,
    regist_day TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    birthday TIMESTAMP NOT NULL
);


INSERT IGNORE INTO user_info (user_id, first_name, last_name, mail_address1, mail_address2, mail_address3, phone_num1, phone_num2, phone_num3, address1, address2, address3, post_code, sex, pay_rank, regist_day, birthday)
VALUES
('c7897f35-4864-d1e5-cc15-3d0df635cbfb', 'Taro', 'Admin', 'ryo29wx@gmail.com', '', '', '07012345678', '', '', 'Tokyo Minato-ku', '', '', 1234567, 1, 1, NOW(), 763916400),
('c7897f35-4864-d1e5-cc15-3d0df635cbfa', 'Ryo', 'Kiuchi', 'ryo29wx@gmail.com', '', '', '07012345678', '', '', 'Tokyo Minato-ku', '', '', 1234567, 1, 1, NOW(), 763916400),
('463bd76b-a6ab-ae4e-18c3-75995d965b3a', 'Yoshiaki', 'Toyama', 'sample1@sample.com', '', '', '012123415093', '', '', 'Tokyo Minato-ku', '', '', 1234567, 1, 1, NOW(), 763916400),
('2b4cdc4a-5b36-85e6-d2e8-4228ca44bae9', 'test01', 'test11', 'sample2@sample.com', '', '', '07012345678', '', '', 'Tokyo Minato-ku', '', '', 1234567, 1, 1, NOW(), 763916400),
('a85a238d-de1d-3e5a-3ea8-5619a6012dc3', '亮', '木内', 'sample3@sample.com', '', '', '07012345678', '', '', 'Tokyo Minato-ku', '', '', 1234567, 1, 1, NOW(), 763916400),
('8c5c1259-8546-34f8-508b-660e6f53bc57', 'テスト', 'テスト03', 'sample4@sample.com', '', '', '07012345678', '', '', 'Tokyo Minato-ku', '', '', 1234567, 1, 1, NOW(), 763916400),
('99bfac31-7784-5675-0215-886408dc37a9', 'test03', 'test13', 'sample5@sample.com', '', '', '07012345678', '', '', 'Tokyo Minato-ku', '', '', 1234567, 1, 1, NOW(), 763916400),
('a1b1d4f5-6dd6-d993-2257-1b9ebe5315f4', 'test04', 'test14', 'sample6@sample.com', '', '', '07012345678', '', '', 'Tokyo Minato-ku', '', '', 1234567, 1, 1, NOW(), 763916400),
('5e3f0c9f-9f26-ef67-8b09-c2dc18e07cd9', 'test05', 'test15', 'sample7@sample.com', '', '', '07012345678', '', '', 'Tokyo Minato-ku', '', '', 1234567, 1, 1, NOW(), 763916400),
('270b38a2-119f-19f6-fa91-500cfeef902b', 'test06', 'test16', 'sample8@sample.com', '', '', '07012345678', '', '', 'Tokyo Minato-ku', '', '', 1234567, 1, 1, NOW(), 763916400),
('d33f77fe-1949-4f9d-e9de-0e4009367748', 'test07', 'test17', 'sample9@sample.com', '', '', '07012345678', '', '', 'Tokyo Minato-ku', '', '', 1234567, 1, 1, NOW(), 763916400);

-- create user_payment_info table
CREATE TABLE IF NOT EXISTS user_payment_info (
    user_id VARCHAR(36) PRIMARY KEY,
    pay_type INT NOT NULL,
    p_status BOOLEAN,
    cc_number INT,
    cc_expiration_day INT,
    cc_sec INT,
    cc_signature  VARCHAR(16),
    FOREIGN KEY (user_id) REFERENCES user_info(user_id)
);