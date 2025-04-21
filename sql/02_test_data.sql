-- Insert test clients
INSERT INTO clients (name, email, address) VALUES
    ('Test Client 1', 'contact@client1.example', 
     E'Smart Building, Floor 3\nTeststreet 1\n12345 Test City\nRomania'),
    ('Test Client 2', 'info@client2.example', 
     E'Business Center Tower\nSample Road 42\n54321 Example Town\nRomania'),
    ('Tech Solutions Ltd', 'projects@techsolutions.example', 
     E'Innovation Park Complex\nBuilding B, Floor 5\n98765 Tech City\nRomania'),
    ('Creative Agency Inc', 'offers@creative.example', 
     E'Design Avenue Plaza\nCreative Quarter\n34567 Creative Town\nRomania'),
    ('Startup Hub Co', 'contracts@startuphub.example', 
     E'Startup Campus\nInnovation District\n89012 Innovation City\nRomania');

-- Insert test employees with different roles (salaries in EUR/year)
INSERT INTO employees (name, email, role, salary) VALUES
    ('John Principal', 'john@company.example', 'Principal', 85000),
    ('Sarah Senior', 'sarah@company.example', 'Senior', 60000),
    ('Mike Professional', 'mike@company.example', 'Professional', 40000),
    ('Lisa Junior', 'lisa@company.example', 'Junior', 28000),
    ('David Senior', 'david@company.example', 'Senior', 62000),
    ('Emma Professional', 'emma@company.example', 'Professional', 42000),
    ('Tom Junior', 'tom@company.example', 'Junior', 27000);

-- Insert employee costs (annual costs in EUR)
INSERT INTO employee_costs (name, cost_per_year) VALUES
    ('Equipment (Laptop, Monitors, Peripherals)', 2000),
    ('Software Licenses', 1200),
    ('Training Budget', 1500),
    ('Office Cost', 3600),
    ('Insurance (Health, Liability)', 1200),
    ('Management Overhead', 7000);

-- Insert offers with various team sizes and statuses
-- Draft offers (3)
INSERT INTO offers (client_id, timeframe, requirements, multiplier, status)
VALUES 
    (1, 2, E'# User Authentication System\n\nImplement a secure authentication system using JWT tokens and role-based access control.\n\n| Feature | Priority |\n|---------|----------|\n| Login | High |\n| Password Reset | Medium |', 1.8, 'draft'),
    (2, 6, E'# E-commerce Integration\n\nDevelop API endpoints for product management and order processing.', 1.7, 'draft'),
    (3, 2, E'# Mobile App Backend\n\nBuild RESTful APIs for the mobile application with proper documentation.', 1.8, 'draft');

-- Sent offers (3)
INSERT INTO offers (client_id, timeframe, requirements, multiplier, status)
VALUES 
    (2, 2, E'# Payment Gateway Integration\n\nImplement secure payment processing with Stripe.', 1.8, 'sent'),
    (3, 6, E'# Data Analytics Dashboard\n\nCreate a real-time analytics dashboard with filtering capabilities.', 1.7, 'sent'),
    (4, 2, E'# Content Management System\n\nDevelop a headless CMS with markdown support.', 1.8, 'sent');

-- Accepted offers (2)
INSERT INTO offers (client_id, timeframe, requirements, multiplier, discount_amount, discount_explanation, status)
VALUES 
    (3, 6, E'# Enterprise Resource Planning System\n\nDevelop core modules for:\n- Human Resources\n- Inventory Management\n- Financial Reporting\n\n## Technical Requirements\n\n| Module | Integration |\n|--------|-------------|\n| HR | Active Directory |\n| Inventory | Existing WMS |\n| Finance | SAP |', 1.7, 3.00, 'Long-term client discount', 'accepted'),
    (5, 2, E'# API Gateway Development\n\nImplement a scalable API gateway with rate limiting.', 1.8, 2.00, 'Strategic partnership discount', 'accepted');

-- Rejected offers (2)
INSERT INTO offers (client_id, timeframe, requirements, multiplier, status)
VALUES 
    (4, 6, E'# Legacy System Migration\n\nMigrate from monolithic architecture to microservices.', 1.9, 'rejected'),
    (5, 2, E'# Reporting Engine\n\nBuild a customizable report generation system.', 1.8, 'rejected');

-- Add employees to offers with different team sizes
-- Draft offers team assignments
INSERT INTO offer_employees (offer_id, employee_id, role) VALUES
    (1, 1, 'Principal'),                              -- 1 person team
    (2, 2, 'Senior'),                                 -- 1 person team
    (3, 1, 'Principal'),                              -- 3 person team
    (3, 3, 'Professional'),
    (3, 4, 'Junior');

-- Sent offers team assignments
INSERT INTO offer_employees (offer_id, employee_id, role) VALUES
    (4, 1, 'Principal'),                              -- 2 person team
    (4, 2, 'Senior'),
    (5, 2, 'Senior'),                                 -- 3 person team
    (5, 3, 'Professional'),
    (5, 4, 'Junior'),
    (6, 5, 'Senior');                                 -- 1 person team

-- Accepted offers team assignments
INSERT INTO offer_employees (offer_id, employee_id, role) VALUES
    (7, 1, 'Principal'),                              -- 3 person team
    (7, 6, 'Professional'),
    (7, 7, 'Junior'),
    (8, 2, 'Senior'),                                 -- 2 person team
    (8, 3, 'Professional');

-- Rejected offers team assignments
INSERT INTO offer_employees (offer_id, employee_id, role) VALUES
    (9, 1, 'Principal'),                              -- 2 person team
    (9, 3, 'Professional'),
    (10, 5, 'Senior');                                -- 1 person team