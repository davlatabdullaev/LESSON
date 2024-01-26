CREATE TABLE branch (
  id UUID PRIMARY KEY,
  name VARCHAR(50),
  address TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sale (
  id UUID PRIMARY KEY,
  branch_id UUID REFERENCES branch(id),
  shop_assistant_id VARCHAR(10),
  cashier_id UUID,
  payment_type payment_type_enum,
  price NUMERIC(15,4),
  status status_enum,
  client_name VARCHAR(50),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transaction (
  id UUID PRIMARY KEY,
  sale_id UUID REFERENCES sale(id),
  staff_id UUID REFERENCES staff(id),
  transaction_type transaction_type_enum,
  source_type source_type_enum,
  amount NUMERIC(15,4),
  description TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE staff (
  id UUID PRIMARY KEY,
  branch_id UUID REFERENCES branch(id),
  tarif_id UUID REFERENCES staff_tarif(id),
  staff_type staff_type_enum,
  name VARCHAR(50),
  balance NUMERIC(15,4),
  birth_date DATE,
  gender VARCHAR(10),
  login VARCHAR(25),
  password VARCHAR(128),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE staff_tarif (
  id UUID PRIMARY KEY,
  name VARCHAR(50),
  tarif_type tarif_type_enum,
  amount_for_cash NUMERIC(15,4),
  amount_for_card NUMERIC(15,4),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE payment_type_enum AS ENUM ('card', 'cash');
CREATE TYPE status_enum AS ENUM ('in_procces', 'success', 'cancel');
CREATE TYPE transaction_type_enum AS ENUM ('withdraw', 'topup');
CREATE TYPE source_type_enum AS ENUM ('bonus', 'sales');
CREATE TYPE staff_type_enum AS ENUM ('shop_assistant', 'cashier');
CREATE TYPE tarif_type_enum AS ENUM ('percent', 'field');
