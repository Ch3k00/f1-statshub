-- Tabla de pilotos
CREATE TABLE IF NOT EXISTS drivers (
    driver_number INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    name_acronym TEXT NOT NULL,
    team_name TEXT NOT NULL,
    country_code TEXT NOT NULL
);

-- Tabla de sesiones (carreras)
CREATE TABLE IF NOT EXISTS sessions (
    session_key INTEGER PRIMARY KEY,
    session_name TEXT NOT NULL,
    session_type TEXT NOT NULL,
    location TEXT NOT NULL,
    country_name TEXT NOT NULL,
    year INTEGER NOT NULL,
    circuit_short_name TEXT NOT NULL,
    date_start TEXT NOT NULL
);

-- Tabla de posiciones
CREATE TABLE IF NOT EXISTS positions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    driver_number INTEGER NOT NULL,
    session_key INTEGER NOT NULL,
    position INTEGER NOT NULL,
    date TEXT NOT NULL,
    FOREIGN KEY (driver_number) REFERENCES drivers(driver_number),
    FOREIGN KEY (session_key) REFERENCES sessions(session_key)
);

-- Tabla de vueltas
CREATE TABLE IF NOT EXISTS laps (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    driver_number INTEGER NOT NULL,
    session_key INTEGER NOT NULL,
    lap_number INTEGER NOT NULL,
    lap_duration REAL,
    duration_sector_1 REAL,
    duration_sector_2 REAL,
    duration_sector_3 REAL,
    st_speed REAL,
    date_start TEXT NOT NULL,
    FOREIGN KEY (driver_number) REFERENCES drivers(driver_number),
    FOREIGN KEY (session_key) REFERENCES sessions(session_key)
);