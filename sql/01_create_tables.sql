-- Create enum for employee roles
CREATE TYPE employee_role AS ENUM ('Principal', 'Senior', 'Professional', 'Junior');
CREATE TYPE offer_status AS ENUM('draft', 'sent', 'accepted', 'rejected');

-- Create clients table
CREATE TABLE IF NOT EXISTS clients (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address TEXT,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create employees table
CREATE TABLE IF NOT EXISTS employees (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    role employee_role NOT NULL,
    salary DECIMAL(10, 2) NOT NULL CHECK (salary > 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create table extra costs for employees for calculation
CREATE TABLE IF NOT EXISTS employee_costs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    cost_per_year DECIMAL(10, 2) NOT NULL CHECK (cost_per_year > 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create offers table
CREATE TABLE IF NOT EXISTS offers (
    id SERIAL PRIMARY KEY,
    client_id INTEGER NOT NULL REFERENCES clients(id),
    timeframe INTEGER NOT NULL CHECK (timeframe IN (2, 6)),
    requirements TEXT,
    multiplier DECIMAL(5, 2),
    discount_amount DECIMAL(5, 2),
    discount_explanation TEXT,
    status offer_status DEFAULT 'draft',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CHECK ((discount_amount IS NULL) OR 
           (discount_amount IS NOT NULL AND discount_explanation IS NOT NULL))
);

-- Create offer_employees junction table
CREATE TABLE IF NOT EXISTS offer_employees (
    id SERIAL PRIMARY KEY,
    role employee_role NOT NULL,
    offer_id INTEGER REFERENCES offers(id),
    employee_id INTEGER REFERENCES employees(id)
);

-- Add indexes
CREATE INDEX idx_clients_email ON clients(email);
CREATE INDEX idx_employees_email ON employees(email);
CREATE INDEX idx_offers_client_id ON offers(client_id);
CREATE INDEX idx_offers_status ON offers(status);
