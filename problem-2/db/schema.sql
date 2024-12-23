CREATE TABLE IF NOT EXISTS sim_records (
    id SERIAL PRIMARY KEY,
    imsi BIGINT NOT NULL UNIQUE,
    pin1 VARCHAR(50),
    puk1 VARCHAR(50),
    pin2 VARCHAR(50),
    puk2 VARCHAR(50),
    aam1 VARCHAR(50),
    ki_umts_enc VARCHAR(50),
    acc VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
