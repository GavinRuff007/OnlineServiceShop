CREATE TABLE products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    productCode VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    price INT NOT NULL,
    status VARCHAR(50) NOT NULL,
    inventory VARCHAR(50) NOT NULL
);
