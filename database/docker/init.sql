CREATE TABLE clients (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    api_key TEXT UNIQUE NOT NULL
);

CREATE TABLE logs (
    id SERIAL PRIMARY KEY,
    client_id INT NOT NULL,
    endpoint TEXT NOT NULL,
    method TEXT NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (client_id) REFERENCES clients(id)
);

CREATE TABLE data (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL
);

-- ตารางผู้ใช้ (User)
CREATE TABLE "users" (
    userID SERIAL PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255),
    role VARCHAR(20) CHECK (role IN ('admin', 'affiliator','member')) DEFAULT 'member'
);

-- ตารางบริการ (Service)
CREATE TABLE "Service" (
    serviceID SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(100),
    description TEXT,
    price DECIMAL(10,2)
);

-- ตารางโปรโมชั่น (Promotion)
CREATE TABLE "Promotion" (
    promotionID SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE,
    discount DECIMAL(5,2),
    description TEXT
);

-- ตารางคำสั่งซื้อ (Order)
CREATE TABLE "Order" (
    orderID SERIAL PRIMARY KEY,
    userID INT,
    serviceID INT,
    promotionID INT,
    FOREIGN KEY (userID) REFERENCES "users"(userID) ON DELETE SET NULL,
    FOREIGN KEY (serviceID) REFERENCES "Service"(serviceID) ON DELETE CASCADE,
    FOREIGN KEY (promotionID) REFERENCES "Promotion"(promotionID) ON DELETE SET NULL
);

-- ตารางตารางนัดหมาย (Schedule)
CREATE TABLE "Schedule" (
    scheduleID SERIAL PRIMARY KEY,
    orderID INT,
    scheduleDate TIMESTAMP NOT NULL,
    FOREIGN KEY (orderID) REFERENCES "Order"(orderID) ON DELETE CASCADE
);

-- ตารางบันทึกการใช้งาน (Log)
CREATE TABLE "Log" (
    logID SERIAL PRIMARY KEY,
    userID INT,
    action VARCHAR(50),
    target VARCHAR(255),
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    referrer VARCHAR(255),
    parameter JSON,
    FOREIGN KEY (userID) REFERENCES "users"(userID) ON DELETE SET NULL
);

-- ตารางเว็บ Affiliator
CREATE TABLE "AffiliatorSite" (
    siteID SERIAL PRIMARY KEY,
    userID INT,
    url VARCHAR(255),
    FOREIGN KEY (userID) REFERENCES "users"(userID) ON DELETE CASCADE
);


-- User
INSERT INTO "users" (email, password, role) VALUES
('admin@example.com', 'hashed_password_123', 'admin'),
('john@example.com', 'hashed_password_456', 'affiliator'),
('jane@example.com', 'hashed_password_789', 'member');

-- Service
INSERT INTO "Service" (name, type, description, price) VALUES
('ติดตั้งแอร์', 'ติดตั้ง', 'บริการติดตั้งเครื่องปรับอากาศ', 1500.00),
('ซ่อมปลั๊กไฟ', 'ซ่อมแซม', 'บริการซ่อมปลั๊กไฟที่ชำรุด', 500.00),
('ล้างแอร์', 'ล้าง', 'บริการล้างเครื่องปรับอากาศ', 800.00);

-- Promotion
INSERT INTO "Promotion" (code, discount, description) VALUES
('SUMMER10', 10.00, 'ส่วนลดช่วงหน้าร้อน'),
('NEWUSER20', 20.00, 'ส่วนลดสำหรับผู้ใช้ใหม่');

-- Order
INSERT INTO "Order" (userID, serviceID, promotionID) VALUES
(2, 1, 1),
(3, 2, NULL);

-- Schedule
INSERT INTO "Schedule" (orderID, scheduleDate) VALUES
(1, '2025-04-10 14:00:00'),
(2, '2025-04-11 09:00:00');

-- Log
INSERT INTO "Log" (userID, action, target, referrer, parameter) VALUES
(2, 'search', 'ติดตั้งแอร์', 'affiliator-website.com', '{"keyword": "แอร์"}'),
(3, 'click', 'ล้างแอร์', 'affiliator-website.com', '{"serviceID": 3}');

-- AffiliatorSite
INSERT INTO "AffiliatorSite" (userID, url) VALUES
(2, 'https://johns-repair-affiliate.com'),
(3, 'https://janes-fixit-partner.net');
